package brokerdaemon

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/keychain"
	opsdk "github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/onepassword"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/runtimeauth"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/version"
	"golang.org/x/sys/unix"
)

type Server struct {
	Config config.Config

	mu             sync.RWMutex
	cachedToken    string
	tokenSource    string
	materialCache  map[string]string
	materialSource map[string]string
}

func (s *Server) Run() error {
	socketPath := s.Config.Broker.SocketPath
	if socketPath == "" {
		return fmt.Errorf("broker socket path is empty")
	}

	if err := os.MkdirAll(filepath.Dir(socketPath), 0o700); err != nil {
		return err
	}
	_ = os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}
	defer listener.Close()
	if err := os.Chmod(socketPath, 0o600); err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	if err := ensureSameUID(conn); err != nil {
		_ = writeResponse(conn, map[string]any{"ok": false, "error": err.Error()})
		return
	}

	request, err := readRequest(conn)
	if err != nil {
		_ = writeResponse(conn, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	_ = conn.SetReadDeadline(time.Time{})
	_ = conn.SetWriteDeadline(time.Now().Add(responseTimeout(request.Action)))

	switch request.Action {
	case "status":
		response := s.statusResponse()
		_ = writeResponse(conn, response)
	case "get_token":
		response := s.tokenResponse()
		_ = writeResponse(conn, response)
	case "get_material":
		response := s.materialResponse(request.Reference)
		_ = writeResponse(conn, response)
	case "cache_material":
		response := s.cacheMaterialResponse(request.Reference, request.KeychainService, request.KeychainAccount)
		_ = writeResponse(conn, response)
	default:
		_ = writeResponse(conn, map[string]any{"ok": false, "error": "unknown_action"})
	}
}

func responseTimeout(action string) time.Duration {
	switch action {
	case "get_material":
		return 90 * time.Second
	default:
		return 5 * time.Second
	}
}

func (s *Server) statusResponse() map[string]any {
	token, source, err := s.loadToken()
	return map[string]any{
		"ok":              err == nil,
		"error":           errorString(err),
		"socket":          s.Config.Broker.SocketPath,
		"pid":             os.Getpid(),
		"token_available": runtimeauth.TokenLooksValid(token),
		"token_length":    len(token),
		"token_source":    source,
	}
}

func (s *Server) tokenResponse() map[string]any {
	token, source, err := s.loadToken()
	if err != nil {
		return map[string]any{
			"ok":              false,
			"error":           err.Error(),
			"token_available": false,
			"token_length":    0,
			"token_source":    source,
		}
	}
	return map[string]any{
		"ok":              true,
		"token":           token,
		"token_available": true,
		"token_length":    len(token),
		"token_source":    source,
	}
}

func (s *Server) loadToken() (string, string, error) {
	s.mu.RLock()
	cachedToken := s.cachedToken
	cachedSource := s.tokenSource
	s.mu.RUnlock()
	if runtimeauth.TokenLooksValid(cachedToken) {
		return cachedToken, cachedSource, nil
	}

	token, source, err := runtimeauth.ResolveServiceAccountTokenForBroker(s.Config)
	if err != nil {
		return "", "", err
	}

	s.mu.Lock()
	s.cachedToken = token
	s.tokenSource = source
	if s.materialCache == nil {
		s.materialCache = map[string]string{}
	}
	if s.materialSource == nil {
		s.materialSource = map[string]string{}
	}
	s.mu.Unlock()
	return token, source, nil
}

func (s *Server) materialResponse(reference string) map[string]any {
	material, source, err := s.loadMaterial(reference)
	if err != nil {
		fmt.Fprintf(os.Stderr, "op-sa-broker: get_material failed for %s: %v\n", reference, err)
	}
	return map[string]any{
		"ok":              err == nil,
		"error":           errorString(err),
		"reference":       reference,
		"material":        material,
		"material_length": len(material),
		"material_source": source,
	}
}

func (s *Server) cacheMaterialResponse(reference string, keychainService string, keychainAccount string) map[string]any {
	if keychainService == "" {
		return map[string]any{
			"ok":    false,
			"error": "keychain_service is required",
		}
	}
	material, source, err := s.loadMaterial(reference)
	if err != nil {
		return map[string]any{
			"ok":        false,
			"error":     err.Error(),
			"reference": reference,
		}
	}
	account := keychainAccount
	if account == "" {
		account = keychain.CurrentUser()
	}
	if err := keychain.UpsertGenericPassword(context.Background(), keychainService, account, material); err != nil {
		return map[string]any{
			"ok":               false,
			"error":            err.Error(),
			"reference":        reference,
			"keychain_service": keychainService,
			"keychain_account": account,
			"material_source":  source,
		}
	}
	return map[string]any{
		"ok":               true,
		"reference":        reference,
		"keychain_service": keychainService,
		"keychain_account": account,
		"value_length":     len(material),
		"material_source":  source,
	}
}

func (s *Server) loadMaterial(reference string) (string, string, error) {
	if reference == "" {
		return "", "", fmt.Errorf("reference is required")
	}

	s.mu.RLock()
	if material, ok := s.materialCache[reference]; ok && material != "" {
		source := s.materialSource[reference]
		s.mu.RUnlock()
		return material, source, nil
	}
	s.mu.RUnlock()

	if strings.HasPrefix(reference, "keychain://") {
		service, account, err := keychain.ParseReference(reference)
		if err != nil {
			return "", "", err
		}
		value, err := keychain.ReadGenericPassword(context.Background(), service, account)
		if err != nil {
			return "", "", err
		}
		source := "macos-keychain:security:broker"
		s.mu.Lock()
		if s.materialCache == nil {
			s.materialCache = map[string]string{}
		}
		if s.materialSource == nil {
			s.materialSource = map[string]string{}
		}
		s.materialCache[reference] = value
		s.materialSource[reference] = source
		s.mu.Unlock()
		return value, source, nil
	}

	token, tokenSource, err := s.loadToken()
	if err != nil {
		return "", "", err
	}

	client := opsdk.ServiceAccountClient{
		IntegrationName:    "aih-go-broker",
		IntegrationVersion: version.Version,
		TokenEnv:           s.Config.SecretService.ServiceAccountEnv,
		Token:              token,
	}
	value, err := client.ResolveReference(context.Background(), reference)
	if err != nil {
		return "", "", err
	}

	source := "onepassword:service_account:" + tokenSource + ":broker"
	s.mu.Lock()
	if s.materialCache == nil {
		s.materialCache = map[string]string{}
	}
	if s.materialSource == nil {
		s.materialSource = map[string]string{}
	}
	s.materialCache[reference] = value
	s.materialSource[reference] = source
	s.mu.Unlock()
	return value, source, nil
}

type requestPayload struct {
	Action          string `json:"action"`
	Reference       string `json:"reference,omitempty"`
	KeychainService string `json:"keychain_service,omitempty"`
	KeychainAccount string `json:"keychain_account,omitempty"`
}

func readRequest(conn net.Conn) (requestPayload, error) {
	line, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return requestPayload{}, err
	}
	var payload requestPayload
	if err := json.Unmarshal(line, &payload); err != nil {
		return requestPayload{}, err
	}
	return payload, nil
}

func writeResponse(conn net.Conn, payload map[string]any) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = conn.Write(append(encoded, '\n'))
	return err
}

func ensureSameUID(conn net.Conn) error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	unixConn, ok := conn.(*net.UnixConn)
	if !ok {
		return fmt.Errorf("unexpected connection type")
	}
	file, err := unixConn.File()
	if err != nil {
		return err
	}
	defer file.Close()

	cred, err := unix.GetsockoptXucred(int(file.Fd()), unix.SOL_LOCAL, unix.LOCAL_PEERCRED)
	if err != nil {
		return err
	}
	if int(cred.Uid) != os.Getuid() {
		return fmt.Errorf("peer uid mismatch: %d", cred.Uid)
	}
	return nil
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
