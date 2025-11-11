package httpz

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Option func(*Client)

type RetryableFunc func(resp *http.Response, err error) bool

type RetryPolicy struct {
	// MaxRetries maximum number of retries. Default is 0
	MaxRetries int
	// MinRetryDelay minimum retry delay. The retry delay will grow exponentially. Default is 0ms
	MinRetryDelay time.Duration
	// ShouldRetry retry judgment function. Default is DefaultRetryableFunc
	ShouldRetry RetryableFunc
}

func WithHttpClient(hc *http.Client) Option {
	return func(c *Client) { c.client = hc }
}

// WithRetryPolicy sets the retry policy.
// The retry policy will only take effect if the request body can be rebuilt through req.GetBody.
// Only request bodies of type *bytes.Buffer, *bytes.Reader, and *strings.Reader will automatically have the
// GetBody method set. For others, please define the GetBody method yourself.
func WithRetryPolicy(policy RetryPolicy) Option {
	return func(c *Client) {
		if policy.MaxRetries > 0 {
			c.retryPolicy.MaxRetries = policy.MaxRetries
		}
		if policy.MinRetryDelay > 0 {
			c.retryPolicy.MinRetryDelay = policy.MinRetryDelay
		}
		if policy.ShouldRetry != nil {
			c.retryPolicy.ShouldRetry = policy.ShouldRetry
		}
	}
}

func WithMiddleware(mws ...Middleware) Option {
	return func(c *Client) {
		c.middlewares = append(c.middlewares, mws...)
	}
}

func DefaultRetryableFunc(resp *http.Response, err error) bool {
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return false
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return false
		}

		return true
	}

	if resp.StatusCode >= 500 {
		return true
	}

	return false
}
