package broker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Client struct {
	SocketPath string
	Timeout    time.Duration
}

type StatusResponse struct {
	OK             bool   `json:"ok"`
	Error          string `json:"error,omitempty"`
	Socket         string `json:"socket,omitempty"`
	TokenAvailable bool   `json:"token_available"`
	TokenLength    int    `json:"token_length"`
	TokenSource    string `json:"token_source,omitempty"`
	PID            int    `json:"pid,omitempty"`
}

type TokenResponse struct {
	OK             bool   `json:"ok"`
	Error          string `json:"error,omitempty"`
	Token          string `json:"token,omitempty"`
	TokenAvailable bool   `json:"token_available"`
	TokenLength    int    `json:"token_length"`
	TokenSource    string `json:"token_source,omitempty"`
}

type MaterialResponse struct {
	OK             bool   `json:"ok"`
	Error          string `json:"error,omitempty"`
	Reference      string `json:"reference,omitempty"`
	Material       string `json:"material,omitempty"`
	MaterialLength int    `json:"material_length"`
	MaterialSource string `json:"material_source,omitempty"`
}

type CacheMaterialResponse struct {
	OK              bool   `json:"ok"`
	Error           string `json:"error,omitempty"`
	Reference       string `json:"reference,omitempty"`
	KeychainService string `json:"keychain_service,omitempty"`
	KeychainAccount string `json:"keychain_account,omitempty"`
	ValueLength     int    `json:"value_length"`
	MaterialSource  string `json:"material_source,omitempty"`
}

func (c Client) Status() (StatusResponse, error) {
	var response StatusResponse
	err := c.request("status", &response)
	return response, err
}

func (c Client) GetToken() (TokenResponse, error) {
	var response TokenResponse
	err := c.request("get_token", &response)
	return response, err
}

func (c Client) GetMaterial(reference string) (MaterialResponse, error) {
	var response MaterialResponse
	timeout := c.Timeout
	if timeout == 0 {
		timeout = 90 * time.Second
	}
	err := c.requestWithPayloadAndTimeout(map[string]string{
		"action":    "get_material",
		"reference": reference,
	}, timeout, &response)
	return response, err
}

func (c Client) CacheMaterial(reference string, keychainService string, keychainAccount string) (CacheMaterialResponse, error) {
	var response CacheMaterialResponse
	timeout := c.Timeout
	if timeout == 0 {
		timeout = 90 * time.Second
	}
	err := c.requestWithPayloadAndTimeout(map[string]string{
		"action":           "cache_material",
		"reference":        reference,
		"keychain_service": keychainService,
		"keychain_account": keychainAccount,
	}, timeout, &response)
	return response, err
}

func (c Client) request(action string, out any) error {
	timeout := c.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return c.requestWithPayloadAndTimeout(map[string]string{"action": action}, timeout, out)
}

func (c Client) requestWithPayload(payloadMap map[string]string, out any) error {
	timeout := c.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return c.requestWithPayloadAndTimeout(payloadMap, timeout, out)
}

func (c Client) requestWithPayloadAndTimeout(payloadMap map[string]string, timeout time.Duration, out any) error {
	conn, err := net.DialTimeout("unix", c.SocketPath, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(timeout))

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return err
	}
	if _, err := conn.Write(append(payload, '\n')); err != nil {
		return err
	}
	if unixConn, ok := conn.(*net.UnixConn); ok {
		_ = unixConn.CloseWrite()
	}

	line, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return err
	}
	if err := json.Unmarshal(line, out); err != nil {
		return fmt.Errorf("decode broker response: %w", err)
	}
	return nil
}
