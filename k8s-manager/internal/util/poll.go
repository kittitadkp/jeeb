package util

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type PollConfig struct {
	Timeout  time.Duration
	Interval time.Duration // defaults to 5s if zero
}

// PollUntil calls check repeatedly until it returns nil or the deadline is reached.
// Returns the last error wrapped with a timeout message.
func PollUntil(ctx context.Context, cfg PollConfig, check func(ctx context.Context) error) error {
	interval := cfg.Interval
	if interval == 0 {
		interval = 5 * time.Second
	}
	deadline := time.Now().Add(cfg.Timeout)
	var lastErr error
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		lastErr = check(ctx)
		if lastErr == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out after %s: %w", cfg.Timeout, lastErr)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}
	}
}

// PollHTTP polls urlStr until the response status equals wantStatus.
// Pass nil for client to use a default 10s-timeout client.
// beforeReq is called before each request (use it for auth headers); pass nil if not needed.
func PollHTTP(ctx context.Context, cfg PollConfig, client *http.Client, method, urlStr string, wantStatus int, beforeReq func(*http.Request)) error {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return PollUntil(ctx, cfg, func(ctx context.Context) error {
		req, err := http.NewRequestWithContext(ctx, method, urlStr, nil)
		if err != nil {
			return err
		}
		if beforeReq != nil {
			beforeReq(req)
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if resp.StatusCode != wantStatus {
			return fmt.Errorf("HTTP %d (want %d)", resp.StatusCode, wantStatus)
		}
		return nil
	})
}
