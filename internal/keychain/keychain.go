package keychain

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	ServiceAccountTokenService = "aiagents/op-service-account-token"
)

func Available() bool {
	return runtime.GOOS == "darwin"
}

func ReadGenericPassword(ctx context.Context, service string, account string) (string, error) {
	if !Available() {
		return "", fmt.Errorf("macos keychain is only available on darwin")
	}
	args := []string{"find-generic-password", "-w", "-s", service}
	if strings.TrimSpace(account) != "" {
		args = append(args, "-a", account)
	}
	output, err := exec.CommandContext(ctx, "security", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("security %s: %w", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(output)), nil
}

func UpsertGenericPassword(ctx context.Context, service string, account string, value string) error {
	if !Available() {
		return fmt.Errorf("macos keychain is only available on darwin")
	}
	args := []string{"add-generic-password", "-U", "-s", service, "-w", value}
	if strings.TrimSpace(account) != "" {
		args = append(args, "-a", account)
	}
	output, err := exec.CommandContext(ctx, "security", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("security %s: %w", strings.Join(args, " "), err)
	}
	_ = output
	return nil
}

func CurrentUser() string {
	return strings.TrimSpace(os.Getenv("USER"))
}

func ParseReference(reference string) (service string, account string, err error) {
	const prefix = "keychain://"
	if !strings.HasPrefix(reference, prefix) {
		return "", "", fmt.Errorf("unsupported macos keychain reference: %s", reference)
	}

	parsed, err := url.Parse(reference)
	if err != nil {
		return "", "", fmt.Errorf("parse keychain reference: %w", err)
	}

	service = strings.TrimSpace(parsed.Host)
	if path := strings.TrimPrefix(parsed.EscapedPath(), "/"); path != "" {
		if service != "" {
			service = service + "/" + path
		} else {
			service = path
		}
	}
	service, err = url.PathUnescape(service)
	if err != nil {
		return "", "", fmt.Errorf("decode keychain service: %w", err)
	}
	if service == "" {
		return "", "", fmt.Errorf("keychain reference is missing service name")
	}
	account = strings.TrimSpace(parsed.Query().Get("account"))
	if account == "" {
		account = CurrentUser()
	}
	return service, account, nil
}
