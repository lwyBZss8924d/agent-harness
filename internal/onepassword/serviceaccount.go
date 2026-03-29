package onepassword

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	sdk "github.com/1password/onepassword-sdk-go"
)

const (
	defaultSDKTimeout = 30 * time.Second
	defaultSDKRetries = 3
)

type ServiceAccountClient struct {
	IntegrationName    string
	IntegrationVersion string
	TokenEnv           string
	Token              string
}

func (c ServiceAccountClient) NewClient(ctx context.Context) (*sdk.Client, error) {
	token := c.Token
	if token == "" {
		token = os.Getenv(c.TokenEnv)
	}
	if token == "" {
		return nil, errors.New("service account token env is empty")
	}

	return sdk.NewClient(
		ctx,
		sdk.WithServiceAccountToken(token),
		sdk.WithIntegrationInfo(c.IntegrationName, c.IntegrationVersion),
	)
}

func (c ServiceAccountClient) ResolveReference(ctx context.Context, reference string) (string, error) {
	var lastErr error
	for attempt := 1; attempt <= defaultSDKRetries; attempt++ {
		callCtx, cancel := ensureTimeout(ctx, defaultSDKTimeout)
		client, err := c.NewClient(callCtx)
		if err == nil {
			var value string
			value, err = client.Secrets().Resolve(callCtx, reference)
			cancel()
			if err == nil {
				return value, nil
			}
			lastErr = fmt.Errorf("resolve attempt %d failed: %w", attempt, err)
		} else {
			cancel()
			lastErr = fmt.Errorf("client init attempt %d failed: %w", attempt, err)
		}
		if attempt < defaultSDKRetries {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}
	return "", lastErr
}

func ensureTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, timeout)
}
