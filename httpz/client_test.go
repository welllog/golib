package httpz

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type jsonCodec struct{}

func (j jsonCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsonCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func TestNewClient(t *testing.T) {
	t.Run("default client", func(t *testing.T) {
		client := NewClient()
		if client == nil {
			t.Fatal("expected client to be created")
		}
		if client.client == nil {
			t.Error("expected http client to be set")
		}
		if client.retryPolicy.ShouldRetry == nil {
			t.Error("expected default retry policy to be set")
		}
	})

	t.Run("with options", func(t *testing.T) {
		customClient := &http.Client{Timeout: 5 * time.Second}
		client := NewClient(
			WithHttpClient(customClient),
			WithRetryPolicy(RetryPolicy{
				MaxRetries:    3,
				MinRetryDelay: 100 * time.Millisecond,
			}),
		)
		if client.client != customClient {
			t.Error("expected custom http client to be set")
		}
		if client.retryPolicy.MaxRetries != 3 {
			t.Errorf("expected MaxRetries to be 3, got %d", client.retryPolicy.MaxRetries)
		}
		if client.retryPolicy.MinRetryDelay != 100*time.Millisecond {
			t.Errorf("expected MinRetryDelay to be 100ms, got %v", client.retryPolicy.MinRetryDelay)
		}
	})

	t.Run("with middleware", func(t *testing.T) {
		mw := func(next DoFunc) DoFunc {
			return func(ctx context.Context, req *http.Request) (*http.Response, error) {
				return next(ctx, req)
			}
		}

		client := NewClient(WithMiddleware(mw))
		if len(client.middlewares) != 1 {
			t.Errorf("expected 1 middleware, got %d", len(client.middlewares))
		}
	})
}

func TestClient_Do(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("success"))
		}))
		defer server.Close()

		client := NewClient()
		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		if string(body) != "success" {
			t.Errorf("expected body 'success', got '%s'", string(body))
		}
	})

	t.Run("failed request", func(t *testing.T) {
		client := NewClient()
		req, _ := http.NewRequest(http.MethodGet, "http://invalid-url-that-does-not-exist.local", nil)
		_, err := client.Do(req)
		if err == nil {
			t.Error("expected error for invalid URL")
		}
	})
}

func TestClient_DoWithRetry(t *testing.T) {
	t.Run("retry on failure", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			if attempts < 3 {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("success"))
			}
		}))
		defer server.Close()

		client := NewClient()
		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

		resp, err := client.DoWithRetry(req, RetryPolicy{
			MaxRetries:    3,
			MinRetryDelay: 10 * time.Millisecond,
			ShouldRetry:   DefaultRetryableFunc,
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		if attempts != 3 {
			t.Errorf("expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("custom retry policy", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		client := NewClient()
		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

		customRetry := func(resp *http.Response, err error) bool {
			return resp != nil && resp.StatusCode == http.StatusBadRequest
		}

		resp, err := client.DoWithRetry(req, RetryPolicy{
			MaxRetries:    2,
			MinRetryDelay: 10 * time.Millisecond,
			ShouldRetry:   customRetry,
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if attempts != 3 { // 1 initial + 2 retries
			t.Errorf("expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("with request body", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			body, _ := io.ReadAll(r.Body)
			if string(body) != "test body" {
				t.Errorf("expected body 'test body', got '%s'", string(body))
			}
			if attempts < 2 {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}))
		defer server.Close()

		client := NewClient()
		req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewBufferString("test body"))

		resp, err := client.DoWithRetry(req, RetryPolicy{
			MaxRetries:    2,
			MinRetryDelay: 10 * time.Millisecond,
			ShouldRetry:   DefaultRetryableFunc,
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if attempts != 2 {
			t.Errorf("expected 2 attempts, got %d", attempts)
		}
	})
}

func TestClient_DoWithoutRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(WithRetryPolicy(RetryPolicy{
		MaxRetries:    3,
		MinRetryDelay: 10 * time.Millisecond,
	}))

	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	resp, err := client.DoWithoutRetry(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if attempts != 1 {
		t.Errorf("expected 1 attempt (no retry), got %d", attempts)
	}
}

func TestClient_Request(t *testing.T) {
	t.Run("with string body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			if string(body) != "test string" {
				t.Errorf("expected body 'test string', got '%s'", string(body))
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"success"}`))
		}))
		defer server.Close()

		client := NewClient()
		var result map[string]string
		err := client.Request(context.Background(), http.MethodPost, server.URL,
			map[string]string{"Content-Type": "text/plain"},
			"test string", &result, jsonCodec{})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result["message"] != "success" {
			t.Errorf("expected message 'success', got '%s'", result["message"])
		}
	})

	t.Run("with bytes body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			if string(body) != "test bytes" {
				t.Errorf("expected body 'test bytes', got '%s'", string(body))
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		}))
		defer server.Close()

		client := NewClient()
		var result map[string]string
		err := client.Request(context.Background(), http.MethodPost, server.URL,
			nil, []byte("test bytes"), &result, jsonCodec{})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result["status"] != "ok" {
			t.Errorf("expected status 'ok', got '%s'", result["status"])
		}
	})

	t.Run("with io.Reader body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			if string(body) != "reader body" {
				t.Errorf("expected body 'reader body', got '%s'", string(body))
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"result":"pass"}`))
		}))
		defer server.Close()

		client := NewClient()
		var result map[string]string
		err := client.Request(context.Background(), http.MethodPost, server.URL,
			nil, strings.NewReader("reader body"), &result, jsonCodec{})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result["result"] != "pass" {
			t.Errorf("expected result 'pass', got '%s'", result["result"])
		}
	})

	t.Run("with struct body (codec)", func(t *testing.T) {
		type RequestBody struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}
		type ResponseBody struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req RequestBody
			body, _ := io.ReadAll(r.Body)
			_ = jsonCodec{}.Unmarshal(body, &req)

			if req.Name != "test" || req.Value != 123 {
				t.Errorf("expected name='test' value=123, got name='%s' value=%d", req.Name, req.Value)
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"success":true,"message":"received"}`))
		}))
		defer server.Close()

		client := NewClient()
		var result ResponseBody
		err := client.Request(context.Background(), http.MethodPost, server.URL,
			map[string]string{"Content-Type": "application/json"},
			RequestBody{Name: "test", Value: 123}, &result, jsonCodec{})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !result.Success {
			t.Error("expected success to be true")
		}
		if result.Message != "received" {
			t.Errorf("expected message 'received', got '%s'", result.Message)
		}
	})

	t.Run("with nil body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":"empty"}`))
		}))
		defer server.Close()

		client := NewClient()
		var result map[string]string
		err := client.Request(context.Background(), http.MethodGet, server.URL,
			nil, nil, &result, jsonCodec{})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result["data"] != "empty" {
			t.Errorf("expected data 'empty', got '%s'", result["data"])
		}
	})

	t.Run("server error 500", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"internal error"}`))
		}))
		defer server.Close()

		client := NewClient()
		var result map[string]string
		err := client.Request(context.Background(), http.MethodGet, server.URL,
			nil, nil, &result, jsonCodec{})

		if err == nil {
			t.Error("expected error for 500 status")
		}
		if !strings.Contains(err.Error(), "http server error") {
			t.Errorf("expected 'http server error' in error message, got: %v", err)
		}
	})

	t.Run("unmarshal error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`invalid json`))
		}))
		defer server.Close()

		client := NewClient()
		var result map[string]string
		err := client.Request(context.Background(), http.MethodGet, server.URL,
			nil, nil, &result, jsonCodec{})

		if err == nil {
			t.Error("expected unmarshal error")
		}
		if !strings.Contains(err.Error(), "unmarshal") {
			t.Errorf("expected 'unmarshal' in error message, got: %v", err)
		}
	})
}

func TestClient_doWithRetry(t *testing.T) {
	t.Run("no retry when MaxRetries is 0", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient()
		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
		resp, err := client.Do(req)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if attempts != 1 {
			t.Errorf("expected 1 attempt, got %d", attempts)
		}
	})

	t.Run("no retry when body cannot be copied", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient(WithRetryPolicy(RetryPolicy{
			MaxRetries:    3,
			MinRetryDelay: 10 * time.Millisecond,
		}))

		// io.NopCloser body doesn't have GetBody
		req, _ := http.NewRequest(http.MethodPost, server.URL, io.NopCloser(strings.NewReader("test")))
		resp, err := client.Do(req)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if attempts != 1 {
			t.Errorf("expected 1 attempt (no retry due to non-copyable body), got %d", attempts)
		}
	})

	t.Run("context canceled no retry", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewClient(WithRetryPolicy(RetryPolicy{
			MaxRetries:    3,
			MinRetryDelay: 10 * time.Millisecond,
		}))

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
		_, err := client.Do(req)

		if err == nil {
			t.Error("expected context deadline exceeded error")
		}

		// Should only try once because context cancellation should not be retried
		if attempts > 1 {
			t.Errorf("expected 1 attempt, got %d", attempts)
		}
	})

	t.Run("GetBody error prevents retry", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient(WithRetryPolicy(RetryPolicy{
			MaxRetries:    3,
			MinRetryDelay: 10 * time.Millisecond,
		}))

		req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewBufferString("test"))
		req.GetBody = func() (io.ReadCloser, error) {
			if attempts > 1 {
				return nil, errors.New("cannot recreate body")
			}
			return io.NopCloser(strings.NewReader("test")), nil
		}

		resp, _ := client.Do(req)
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}

		if attempts != 2 {
			t.Errorf("expected 2 attempts (stops when GetBody fails), got %d", attempts)
		}
	})
}
