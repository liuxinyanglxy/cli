// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package cmdutil

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDebugTransport_LogsRequestAndResponse(t *testing.T) {
	logBuf := &bytes.Buffer{}

	// Mock transport that returns a 200 response
	transport := &DebugTransport{
		Base: mockRoundTripper(func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				Status:        "200 OK",
				StatusCode:    200,
				Header:        http.Header{"Content-Type": []string{"application/json"}},
				Body:          io.NopCloser(bytes.NewBufferString(`{}`)),
				ContentLength: 2,
				Request:       req,
			}
			return resp, nil
		}),
		Out: logBuf,
	}

	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	req.Header.Set("Authorization", "Bearer token123")

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}
	defer resp.Body.Close()

	logOutput := logBuf.String()

	// Should log the request method and URL
	if !strings.Contains(logOutput, "GET") || !strings.Contains(logOutput, "example.com") {
		t.Errorf("log should contain request method and URL, got: %s", logOutput)
	}

	// Should log response status
	if !strings.Contains(logOutput, "200") {
		t.Errorf("log should contain response status code, got: %s", logOutput)
	}

	// Should mask sensitive headers
	if strings.Contains(logOutput, "token123") {
		t.Errorf("log should mask Authorization header value, got: %s", logOutput)
	}
	if strings.Contains(logOutput, "Authorization: Bearer") {
		t.Errorf("log should show Authorization key but mask value, got: %s", logOutput)
	}
}

// mockRoundTripper is a helper to create a mock http.RoundTripper for testing.
type mockRoundTripper func(*http.Request) (*http.Response, error)

func (m mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m(req)
}

func TestDebugTransport_WithoutBase(t *testing.T) {
	logBuf := &bytes.Buffer{}
	transport := &DebugTransport{
		Out: logBuf,
	}

	if transport.base() == nil {
		t.Error("base() should return a fallback transport")
	}
}
