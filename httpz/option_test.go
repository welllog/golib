package httpz

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestWithHttpClient(t *testing.T) {
	t.Run("set custom http client", func(t *testing.T) {
		customClient := &http.Client{
			Timeout: 10 * time.Second,
		}

		client := NewClient(WithHttpClient(customClient))
		if client.client != customClient {
			t.Error("expected custom http client to be set")
		}
		if client.client.Timeout != 10*time.Second {
			t.Errorf("expected timeout 10s, got %v", client.client.Timeout)
		}
	})

	t.Run("set nil http client", func(t *testing.T) {
		client := NewClient(WithHttpClient(nil))
		if client.client != nil {
			t.Error("expected http client to be nil")
		}
	})
}

func TestWithRetryPolicy(t *testing.T) {
	t.Run("set all retry policy fields", func(t *testing.T) {
		customRetry := func(resp *http.Response, err error) bool {
			return false
		}

		policy := RetryPolicy{
			MaxRetries:    5,
			MinRetryDelay: 200 * time.Millisecond,
			ShouldRetry:   customRetry,
		}

		client := NewClient(WithRetryPolicy(policy))

		if client.retryPolicy.MaxRetries != 5 {
			t.Errorf("expected MaxRetries 5, got %d", client.retryPolicy.MaxRetries)
		}
		if client.retryPolicy.MinRetryDelay != 200*time.Millisecond {
			t.Errorf("expected MinRetryDelay 200ms, got %v", client.retryPolicy.MinRetryDelay)
		}
		if client.retryPolicy.ShouldRetry == nil {
			t.Error("expected ShouldRetry to be set")
		}
	})

	t.Run("set only MaxRetries", func(t *testing.T) {
		policy := RetryPolicy{
			MaxRetries: 3,
		}

		client := NewClient(WithRetryPolicy(policy))

		if client.retryPolicy.MaxRetries != 3 {
			t.Errorf("expected MaxRetries 3, got %d", client.retryPolicy.MaxRetries)
		}
		if client.retryPolicy.MinRetryDelay != 0 {
			t.Errorf("expected MinRetryDelay 0, got %v", client.retryPolicy.MinRetryDelay)
		}
		if client.retryPolicy.ShouldRetry == nil {
			t.Error("expected default ShouldRetry to be set")
		}
	})

	t.Run("set only MinRetryDelay", func(t *testing.T) {
		policy := RetryPolicy{
			MinRetryDelay: 100 * time.Millisecond,
		}

		client := NewClient(WithRetryPolicy(policy))

		if client.retryPolicy.MaxRetries != 0 {
			t.Errorf("expected MaxRetries 0, got %d", client.retryPolicy.MaxRetries)
		}
		if client.retryPolicy.MinRetryDelay != 100*time.Millisecond {
			t.Errorf("expected MinRetryDelay 100ms, got %v", client.retryPolicy.MinRetryDelay)
		}
	})

	t.Run("set only ShouldRetry", func(t *testing.T) {
		customRetry := func(resp *http.Response, err error) bool {
			return true
		}

		policy := RetryPolicy{
			ShouldRetry: customRetry,
		}

		client := NewClient(WithRetryPolicy(policy))

		if client.retryPolicy.MaxRetries != 0 {
			t.Errorf("expected MaxRetries 0, got %d", client.retryPolicy.MaxRetries)
		}
		if client.retryPolicy.ShouldRetry == nil {
			t.Error("expected ShouldRetry to be set")
		}
	})

	t.Run("zero values are not applied", func(t *testing.T) {
		policy := RetryPolicy{
			MaxRetries:    0,
			MinRetryDelay: 0,
		}

		client := NewClient(WithRetryPolicy(policy))

		// Zero values should not override defaults
		if client.retryPolicy.MaxRetries != 0 {
			t.Errorf("expected MaxRetries 0, got %d", client.retryPolicy.MaxRetries)
		}
		if client.retryPolicy.MinRetryDelay != 0 {
			t.Errorf("expected MinRetryDelay 0, got %v", client.retryPolicy.MinRetryDelay)
		}
	})

	t.Run("negative values are ignored", func(t *testing.T) {
		policy := RetryPolicy{
			MaxRetries:    -1,
			MinRetryDelay: -100 * time.Millisecond,
		}

		client := NewClient(WithRetryPolicy(policy))

		// Negative values should not be applied
		if client.retryPolicy.MaxRetries < 0 {
			t.Errorf("expected MaxRetries to not be negative, got %d", client.retryPolicy.MaxRetries)
		}
		if client.retryPolicy.MinRetryDelay < 0 {
			t.Errorf("expected MinRetryDelay to not be negative, got %v", client.retryPolicy.MinRetryDelay)
		}
	})
}

func TestWithMiddleware(t *testing.T) {
	t.Run("add single middleware", func(t *testing.T) {
		mw := func(next DoFunc) DoFunc {
			return next
		}

		client := NewClient(WithMiddleware(mw))

		if len(client.middlewares) != 1 {
			t.Errorf("expected 1 middleware, got %d", len(client.middlewares))
		}
	})

	t.Run("add multiple middlewares", func(t *testing.T) {
		mw1 := func(next DoFunc) DoFunc { return next }
		mw2 := func(next DoFunc) DoFunc { return next }
		mw3 := func(next DoFunc) DoFunc { return next }

		client := NewClient(WithMiddleware(mw1, mw2, mw3))

		if len(client.middlewares) != 3 {
			t.Errorf("expected 3 middlewares, got %d", len(client.middlewares))
		}
	})

	t.Run("add middlewares in multiple calls", func(t *testing.T) {
		mw1 := func(next DoFunc) DoFunc { return next }
		mw2 := func(next DoFunc) DoFunc { return next }

		client := NewClient(
			WithMiddleware(mw1),
			WithMiddleware(mw2),
		)

		if len(client.middlewares) != 2 {
			t.Errorf("expected 2 middlewares, got %d", len(client.middlewares))
		}
	})

	t.Run("empty middleware list", func(t *testing.T) {
		client := NewClient(WithMiddleware())

		if len(client.middlewares) != 0 {
			t.Errorf("expected 0 middlewares, got %d", len(client.middlewares))
		}
	})
}

func TestDefaultRetryableFunc(t *testing.T) {
	t.Run("retry on generic error", func(t *testing.T) {
		result := DefaultRetryableFunc(nil, errors.New("connection error"))
		if !result {
			t.Error("expected to retry on generic error")
		}
	})

	t.Run("no retry on context canceled", func(t *testing.T) {
		result := DefaultRetryableFunc(nil, context.Canceled)
		if result {
			t.Error("expected not to retry on context canceled")
		}
	})

	t.Run("no retry on context deadline exceeded", func(t *testing.T) {
		result := DefaultRetryableFunc(nil, context.DeadlineExceeded)
		if result {
			t.Error("expected not to retry on context deadline exceeded")
		}
	})

	t.Run("retry on 500 status", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusInternalServerError}
		result := DefaultRetryableFunc(resp, nil)
		if !result {
			t.Error("expected to retry on 500 status")
		}
	})

	t.Run("retry on 502 status", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusBadGateway}
		result := DefaultRetryableFunc(resp, nil)
		if !result {
			t.Error("expected to retry on 502 status")
		}
	})

	t.Run("retry on 503 status", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusServiceUnavailable}
		result := DefaultRetryableFunc(resp, nil)
		if !result {
			t.Error("expected to retry on 503 status")
		}
	})

	t.Run("retry on 504 status", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusGatewayTimeout}
		result := DefaultRetryableFunc(resp, nil)
		if !result {
			t.Error("expected to retry on 504 status")
		}
	})

	t.Run("no retry on 200 status", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusOK}
		result := DefaultRetryableFunc(resp, nil)
		if result {
			t.Error("expected not to retry on 200 status")
		}
	})

	t.Run("no retry on 400 status", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusBadRequest}
		result := DefaultRetryableFunc(resp, nil)
		if result {
			t.Error("expected not to retry on 400 status")
		}
	})

	t.Run("no retry on 404 status", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusNotFound}
		result := DefaultRetryableFunc(resp, nil)
		if result {
			t.Error("expected not to retry on 404 status")
		}
	})

	t.Run("no retry on 201 status", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusCreated}
		result := DefaultRetryableFunc(resp, nil)
		if result {
			t.Error("expected not to retry on 201 status")
		}
	})

	t.Run("no retry on 204 status", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusNoContent}
		result := DefaultRetryableFunc(resp, nil)
		if result {
			t.Error("expected not to retry on 204 status")
		}
	})

	t.Run("error takes precedence over response", func(t *testing.T) {
		// When both error and response exist, error handling should take precedence
		resp := &http.Response{StatusCode: http.StatusOK}
		result := DefaultRetryableFunc(resp, errors.New("network error"))
		if !result {
			t.Error("expected to retry when error is present, regardless of response")
		}
	})

	t.Run("wrapped context canceled error", func(t *testing.T) {
		wrappedErr := errors.New("failed: " + context.Canceled.Error())
		result := DefaultRetryableFunc(nil, wrappedErr)
		// This should retry because it's not directly context.Canceled
		if !result {
			t.Error("expected to retry on wrapped error (not directly context.Canceled)")
		}
	})
}

func TestRetryPolicy(t *testing.T) {
	t.Run("default retry policy", func(t *testing.T) {
		policy := RetryPolicy{}

		if policy.MaxRetries != 0 {
			t.Errorf("expected default MaxRetries 0, got %d", policy.MaxRetries)
		}
		if policy.MinRetryDelay != 0 {
			t.Errorf("expected default MinRetryDelay 0, got %v", policy.MinRetryDelay)
		}
		if policy.ShouldRetry != nil {
			t.Error("expected default ShouldRetry to be nil")
		}
	})

	t.Run("custom retry policy", func(t *testing.T) {
		customFunc := func(resp *http.Response, err error) bool {
			return resp != nil && resp.StatusCode == http.StatusTooManyRequests
		}

		policy := RetryPolicy{
			MaxRetries:    10,
			MinRetryDelay: 500 * time.Millisecond,
			ShouldRetry:   customFunc,
		}

		if policy.MaxRetries != 10 {
			t.Errorf("expected MaxRetries 10, got %d", policy.MaxRetries)
		}
		if policy.MinRetryDelay != 500*time.Millisecond {
			t.Errorf("expected MinRetryDelay 500ms, got %v", policy.MinRetryDelay)
		}
		if policy.ShouldRetry == nil {
			t.Error("expected ShouldRetry to be set")
		}

		// Test custom retry function
		resp := &http.Response{StatusCode: http.StatusTooManyRequests}
		if !policy.ShouldRetry(resp, nil) {
			t.Error("expected custom retry function to return true for 429 status")
		}

		resp = &http.Response{StatusCode: http.StatusInternalServerError}
		if policy.ShouldRetry(resp, nil) {
			t.Error("expected custom retry function to return false for 500 status")
		}
	})
}

func TestMultipleOptions(t *testing.T) {
	t.Run("combine all options", func(t *testing.T) {
		customClient := &http.Client{Timeout: 30 * time.Second}
		mw1 := func(next DoFunc) DoFunc { return next }
		mw2 := func(next DoFunc) DoFunc { return next }

		client := NewClient(
			WithHttpClient(customClient),
			WithRetryPolicy(RetryPolicy{
				MaxRetries:    5,
				MinRetryDelay: 100 * time.Millisecond,
			}),
			WithMiddleware(mw1, mw2),
		)

		if client.client != customClient {
			t.Error("expected custom http client to be set")
		}
		if client.retryPolicy.MaxRetries != 5 {
			t.Errorf("expected MaxRetries 5, got %d", client.retryPolicy.MaxRetries)
		}
		if client.retryPolicy.MinRetryDelay != 100*time.Millisecond {
			t.Errorf("expected MinRetryDelay 100ms, got %v", client.retryPolicy.MinRetryDelay)
		}
		if len(client.middlewares) != 2 {
			t.Errorf("expected 2 middlewares, got %d", len(client.middlewares))
		}
	})
}
