package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// HTTPClient wraps Fiber.Test() for convenient HTTP testing.
type HTTPClient struct {
	app     *fiber.App
	token   string
	baseURL string
}

// NewHTTPClient creates a new HTTP test client for the given Fiber app.
func NewHTTPClient(app *fiber.App) *HTTPClient {
	return &HTTPClient{
		app:     app,
		baseURL: "",
	}
}

// WithToken sets the Authorization header for subsequent requests.
func (c *HTTPClient) WithToken(token string) *HTTPClient {
	c.token = token
	return c
}

// Get performs a GET request.
func (c *HTTPClient) Get(path string) *HTTPResponse {
	return c.request("GET", path, nil, nil)
}

// Post performs a POST request with JSON body.
func (c *HTTPClient) Post(path string, body interface{}) *HTTPResponse {
	return c.request("POST", path, body, nil)
}

// Put performs a PUT request with JSON body.
func (c *HTTPClient) Put(path string, body interface{}) *HTTPResponse {
	return c.request("PUT", path, body, nil)
}

// Delete performs a DELETE request.
func (c *HTTPClient) Delete(path string) *HTTPResponse {
	return c.request("DELETE", path, nil, nil)
}

// Patch performs a PATCH request with JSON body.
func (c *HTTPClient) Patch(path string, body interface{}) *HTTPResponse {
	return c.request("PATCH", path, body, nil)
}

// request performs the actual HTTP request.
func (c *HTTPClient) request(method, path string, body interface{}, headers map[string]string) *HTTPResponse {
	var reqBody io.Reader
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewReader(jsonData)
	}

	req := httptest.NewRequest(method, c.baseURL+path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.app.Test(req)
	if err != nil {
		panic(fmt.Sprintf("HTTP test failed: %v", err))
	}

	return &HTTPResponse{resp: resp}
}

// HTTPResponse wraps httptest.ResponseRecorder for assertions.
type HTTPResponse struct {
	resp *http.Response
	body []byte
}

// Status returns the HTTP status code.
func (r *HTTPResponse) Status() int {
	return r.resp.StatusCode
}

// Body returns the response body as string.
func (r *HTTPResponse) Body() string {
	if r.body == nil {
		body, _ := io.ReadAll(r.resp.Body)
		r.body = body
	}
	return string(r.body)
}

// JSON unmarshals the response body into the given interface.
func (r *HTTPResponse) JSON(v interface{}) error {
	return json.Unmarshal([]byte(r.Body()), v)
}

// AssertStatus asserts the response status code.
func (r *HTTPResponse) AssertStatus(expected int) *HTTPResponse {
	if r.Status() != expected {
		panic(fmt.Sprintf("expected status %d, got %d. Body: %s", expected, r.Status(), r.Body()))
	}
	return r
}

// AssertJSON asserts the response body matches the expected JSON.
func (r *HTTPResponse) AssertJSON(expected interface{}) *HTTPResponse {
	expectedJSON, _ := json.Marshal(expected)
	actualStr := strings.TrimSpace(r.Body())

	expectedStr := string(expectedJSON)

	if expectedStr != actualStr {
		panic(fmt.Sprintf("expected JSON %s, got %s", expectedStr, actualStr))
	}
	return r
}

// AssertContains asserts the response body contains the given substring.
func (r *HTTPResponse) AssertContains(substring string) *HTTPResponse {
	if !strings.Contains(r.Body(), substring) {
		panic(fmt.Sprintf("expected body to contain %q, got %s", substring, r.Body()))
	}
	return r
}
