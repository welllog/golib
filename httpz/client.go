package httpz

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/welllog/golib/strz"
)

const (
	ClientRetryPolicyKey = "httpz_retry_policy"
)

type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

type Client struct {
	client      *http.Client
	retryPolicy RetryPolicy
	middlewares []Middleware
	doChain     DoFunc
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		client: http.DefaultClient,
		retryPolicy: RetryPolicy{
			ShouldRetry: DefaultRetryableFunc,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	finalDo := func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return c.doWithRetry(ctx, req)
	}

	doChain := finalDo
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		doChain = c.middlewares[i](doChain)
	}
	c.doChain = doChain

	return c
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.doChain(req.Context(), req)
}

// DoWithRetry retryPolicy will override the default retry policy
// executing retries requires the ability to rebuild the request body via req.GetBody.
// Only request bodies of type *bytes.Buffer, *bytes.Reader, and *strings.Reader will automatically have the
// GetBody method set. For others, please define the GetBody method yourself.
func (c *Client) DoWithRetry(req *http.Request, retryPolicy RetryPolicy) (*http.Response, error) {
	ctx := context.WithValue(req.Context(), ClientRetryPolicyKey, retryPolicy)
	return c.doChain(ctx, req)
}

// DoWithoutRetry executes the request without any retry attempts.
func (c *Client) DoWithoutRetry(req *http.Request) (*http.Response, error) {
	ctx := context.WithValue(req.Context(), ClientRetryPolicyKey, RetryPolicy{})
	return c.doChain(ctx, req)
}

// Request
// body supports string, []byte, io.Reader, and other types serialized through codec
func (c *Client) Request(ctx context.Context, method, path string, headers map[string]string, body, out any,
	codec Codec) (err error) {

	var req *http.Request
	if body == nil {
		req, err = http.NewRequestWithContext(ctx, method, path, nil)
	} else {
		switch r := body.(type) {
		case string:
			req, err = http.NewRequestWithContext(ctx, method, path, strings.NewReader(r))
		case []byte:
			req, err = http.NewRequestWithContext(ctx, method, path, bytes.NewBuffer(r))
		case io.Reader:
			req, err = http.NewRequestWithContext(ctx, method, path, r)
		default:
			var b []byte
			b, err = codec.Marshal(body)
			if err != nil {
				return fmt.Errorf("marshal request body failed: %w", err)
			}

			req, err = http.NewRequestWithContext(ctx, method, path, bytes.NewBuffer(b))
		}
	}

	if err != nil {
		return fmt.Errorf("create http request failed: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("send http request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("http server error: status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read http response body failed: %w", err)
	}

	err = codec.Unmarshal(b, out)
	if err != nil {
		return fmt.Errorf("unmarshal http response body failed: %w, body: %s", err, strz.UnsafeString(b))
	}

	return nil
}

func (c *Client) doWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	policyVal := ctx.Value(ClientRetryPolicyKey)
	retryPolicy, ok := policyVal.(RetryPolicy)
	if !ok { // use global retry policy
		retryPolicy = c.retryPolicy
	} else if retryPolicy.ShouldRetry == nil {
		retryPolicy.ShouldRetry = c.retryPolicy.ShouldRetry
	}

	if retryPolicy.MaxRetries <= 0 {
		return c.do(ctx, req)
	}

	// no retry if body cannot be copied
	if req.Body != nil && req.GetBody == nil {
		return c.do(ctx, req)
	}

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= retryPolicy.MaxRetries; attempt++ {
		if attempt > 0 {
			if req.Body != nil {
				copyBody, copyErr := req.GetBody()
				if copyErr != nil {
					// copy failed, cannot retry
					break
				}
				req.Body = copyBody
			}

			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}

			backoff := retryPolicy.MinRetryDelay * time.Duration(1<<uint(attempt-1))

			if backoff > 0 {
				time.Sleep(backoff)
			}
		}

		resp, err = c.do(ctx, req)
		if !retryPolicy.ShouldRetry(resp, err) {
			return resp, err
		}
	}

	return resp, err
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return c.client.Do(req.WithContext(ctx))
}
