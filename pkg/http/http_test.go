package http

import (
	"bytes"
	"io"
	nethttp "net/http"
	"testing"
	"time"
)

type testPayload struct {
	Message string `json:"message"`
}

func TestClient_Post_Into(t *testing.T) {
	client := Fake(func(req *nethttp.Request) (*nethttp.Response, error) {
		if req.Method != nethttp.MethodPost {
			t.Fatalf("expected POST request, got %s", req.Method)
		}
		return &nethttp.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Header:     nethttp.Header{"Content-Type": {"application/json"}},
			Body:       ioNopCloser(bytes.NewReader([]byte(`{"message":"hello"}`))),
		}, nil
	})

	var response testPayload
	if err := client.Post("https://example.com/api", map[string]string{"name": "test"}).Into(&response); err != nil {
		t.Fatalf("Into failed: %v", err)
	}

	if response.Message != "hello" {
		t.Fatalf("unexpected response message: %s", response.Message)
	}
}

func TestClient_WithTokenAndHeader(t *testing.T) {
	client := Fake(func(req *nethttp.Request) (*nethttp.Response, error) {
		if req.Header.Get("Authorization") != "Bearer secret" {
			t.Fatalf("expected bearer token header, got %s", req.Header.Get("Authorization"))
		}
		if req.Header.Get("X-Custom") != "value" {
			t.Fatalf("expected custom header, got %s", req.Header.Get("X-Custom"))
		}
		return &nethttp.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Header:     nethttp.Header{"Content-Type": {"application/json"}},
			Body:       ioNopCloser(bytes.NewReader([]byte(`{"message":"ok"}`))),
		}, nil
	})

	var response testPayload
	err := client.WithToken("secret").WithHeader("X-Custom", "value").Post("https://example.com", nil).Into(&response)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestClient_RetryOnServerError(t *testing.T) {
	attempts := 0
	client := Fake(func(req *nethttp.Request) (*nethttp.Response, error) {
		attempts++
		if attempts == 1 {
			return &nethttp.Response{
				StatusCode: 502,
				Status:     "502 Bad Gateway",
				Header:     nethttp.Header{"Content-Type": {"application/json"}},
				Body:       ioNopCloser(bytes.NewReader([]byte(`{"message":"error"}`))),
			}, nil
		}
		return &nethttp.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Header:     nethttp.Header{"Content-Type": {"application/json"}},
			Body:       ioNopCloser(bytes.NewReader([]byte(`{"message":"success"}`))),
		}, nil
	})

	var response testPayload
	client.WithRetry(2, ExponentialBackoff(1*time.Millisecond))
	err := client.Post("https://example.com/retry", nil).Into(&response)
	if err != nil {
		t.Fatalf("retry request failed: %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
	if response.Message != "success" {
		t.Fatalf("unexpected response after retry: %s", response.Message)
	}
}

func TestFakeClient_DoesNotSendRealRequest(t *testing.T) {
	client := Fake(func(req *nethttp.Request) (*nethttp.Response, error) {
		return &nethttp.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Header:     nethttp.Header{"Content-Type": {"application/json"}},
			Body:       ioNopCloser(bytes.NewReader([]byte(`{"message":"fake"}`))),
		}, nil
	})
	var response testPayload
	if err := client.Get("https://example.com/test").Into(&response); err != nil {
		t.Fatalf("fake request failed: %v", err)
	}
	if response.Message != "fake" {
		t.Fatalf("unexpected fake response: %s", response.Message)
	}
}

type nopCloser struct {
	*bytes.Reader
}

func (n nopCloser) Close() error { return nil }

func ioNopCloser(r *bytes.Reader) io.ReadCloser {
	return nopCloser{Reader: r}
}
