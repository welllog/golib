package httpz

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMultipleMiddlewares(t *testing.T) {
	t.Run("middleware execution order", func(t *testing.T) {
		var executionOrder []string

		mw1 := func(next DoFunc) DoFunc {
			return func(ctx context.Context, req *http.Request) (*http.Response, error) {
				executionOrder = append(executionOrder, "mw1-before")
				resp, err := next(ctx, req)
				executionOrder = append(executionOrder, "mw1-after")
				return resp, err
			}
		}

		mw2 := func(next DoFunc) DoFunc {
			return func(ctx context.Context, req *http.Request) (*http.Response, error) {
				executionOrder = append(executionOrder, "mw2-before")
				resp, err := next(ctx, req)
				executionOrder = append(executionOrder, "mw2-after")
				return resp, err
			}
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "handler")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewClient(WithMiddleware(mw1, mw2))
		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
		resp, err := client.Do(req)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer resp.Body.Close()

		expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
		if len(executionOrder) != len(expected) {
			t.Fatalf("expected %d execution steps, got %d", len(expected), len(executionOrder))
		}

		for i, step := range expected {
			if executionOrder[i] != step {
				t.Errorf("at position %d: expected '%s', got '%s'", i, step, executionOrder[i])
			}
		}
	})
}
