package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	nethttp "net/http"
	"time"
)

type BackoffPolicy func(attempt int) time.Duration

type FakeHandler func(request *nethttp.Request) (*nethttp.Response, error)

// Client is a fluent HTTP client with JSON helper methods.
type Client struct {
	httpClient     *nethttp.Client
	defaultHeaders map[string]string
	token          string
	timeout        time.Duration
	retryCount     int
	backoff        BackoffPolicy
}

// NewClient returns a client configured with the default net/http client.
func NewClient() *Client {
	return &Client{httpClient: &nethttp.Client{}}
}

// Fake returns a client that uses a fake transport for tests.
func Fake(handler FakeHandler) *Client {
	return &Client{
		httpClient: &nethttp.Client{Transport: roundTripperFunc(func(req *nethttp.Request) (*nethttp.Response, error) {
			return handler(req)
		})},
	}
}

type roundTripperFunc func(*nethttp.Request) (*nethttp.Response, error)

func (r roundTripperFunc) RoundTrip(req *nethttp.Request) (*nethttp.Response, error) {
	return r(req)
}

// WithToken sets the default Bearer token for all requests.
func (c *Client) WithToken(token string) *Client {
	c.token = token
	return c
}

// WithHeader adds a default header for all requests.
func (c *Client) WithHeader(name, value string) *Client {
	if c.defaultHeaders == nil {
		c.defaultHeaders = make(map[string]string)
	}
	c.defaultHeaders[name] = value
	return c
}

// WithTimeout sets the default request timeout.
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.timeout = timeout
	return c
}

// WithRetry configures automatic retry for 5xx responses.
func (c *Client) WithRetry(attempts int, backoff BackoffPolicy) *Client {
	c.retryCount = attempts
	c.backoff = backoff
	return c
}

// Get creates a GET request builder.
func (c *Client) Get(url string) *Request {
	return c.newRequest(nethttp.MethodGet, url, nil)
}

// Post creates a POST request builder.
func (c *Client) Post(url string, body any) *Request {
	return c.newRequest(nethttp.MethodPost, url, body)
}

// Put creates a PUT request builder.
func (c *Client) Put(url string, body any) *Request {
	return c.newRequest(nethttp.MethodPut, url, body)
}

func (c *Client) newRequest(method, url string, body any) *Request {
	h := make(map[string]string, len(c.defaultHeaders))
	for k, v := range c.defaultHeaders {
		h[k] = v
	}
	return &Request{
		client:     c,
		method:     method,
		url:        url,
		body:       body,
		headers:    h,
		token:      c.token,
		timeout:    c.timeout,
		retryCount: c.retryCount,
		backoff:    c.backoff,
	}
}

// Request is a fluent HTTP request builder.
type Request struct {
	client     *Client
	method     string
	url        string
	body       any
	headers    map[string]string
	token      string
	timeout    time.Duration
	retryCount int
	backoff    BackoffPolicy
}

// WithToken sets a Bearer token for this request.
func (r *Request) WithToken(token string) *Request {
	r.token = token
	return r
}

// WithHeader adds a header for this request.
func (r *Request) WithHeader(name, value string) *Request {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers[name] = value
	return r
}

// WithTimeout sets a timeout for this request.
func (r *Request) WithTimeout(timeout time.Duration) *Request {
	r.timeout = timeout
	return r
}

// WithRetry configures retry for this request.
func (r *Request) WithRetry(attempts int, backoff BackoffPolicy) *Request {
	r.retryCount = attempts
	r.backoff = backoff
	return r
}

// Into executes the request and unmarshals the JSON response into dest.
func (r *Request) Into(dest any) error {
	resp, err := r.doRequest()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("http %d: %s", resp.StatusCode, resp.Status)
	}

	if dest == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(dest)
}

func (r *Request) doRequest() (*nethttp.Response, error) {
	var bodyReader *bytes.Reader
	var contentType string
	if r.body != nil {
		buf, err := json.Marshal(r.body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(buf)
		contentType = "application/json"
	} else {
		bodyReader = bytes.NewReader(nil)
	}

	req, err := nethttp.NewRequest(r.method, r.url, bodyReader)
	if err != nil {
		return nil, err
	}

	for name, value := range r.headers {
		req.Header.Set(name, value)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", "application/json")
	if r.token != "" {
		req.Header.Set("Authorization", "Bearer "+r.token)
	}

	client := r.client.httpClient
	if client == nil {
		client = &nethttp.Client{}
	}

	ctx := context.Background()
	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}

	rqCount := r.retryCount
	backoff := r.backoff
	if rqCount == 0 {
		rqCount = r.client.retryCount
		backoff = r.client.backoff
	}

	lastErr := error(nil)
	for attempt := 0; attempt < rqCount+1; attempt++ {
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < rqCount {
				wait := backoffDuration(backoff, attempt)
				time.Sleep(wait)
				continue
			}
			return nil, err
		}

		if shouldRetry(resp.StatusCode) {
			if attempt < rqCount {
				_ = resp.Body.Close()
				wait := backoffDuration(backoff, attempt)
				time.Sleep(wait)
				continue
			}
		}
		return resp, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.New("http request failed")
}

func shouldRetry(status int) bool {
	return status >= 500 && status < 600
}

func backoffDuration(policy BackoffPolicy, attempt int) time.Duration {
	if policy == nil {
		return 0
	}
	return policy(attempt)
}

// ExponentialBackoff returns a backoff policy with an increasing delay.
func ExponentialBackoff(base time.Duration) BackoffPolicy {
	return func(attempt int) time.Duration {
		if attempt <= 0 {
			return base
		}
		return base * (1 << attempt)
	}
}
