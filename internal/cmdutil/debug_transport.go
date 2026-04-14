// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package cmdutil

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/larksuite/cli/internal/util"
)

// DebugTransport is an http.RoundTripper that logs HTTP requests and responses.
// It wraps another transport and logs debug information to the given io.Writer.
type DebugTransport struct {
	Base http.RoundTripper
	Out  io.Writer
}

func (t *DebugTransport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return util.FallbackTransport()
}

// RoundTrip implements http.RoundTripper and logs request/response details.
func (t *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log request
	fmt.Fprintf(t.Out, "[DEBUG] %s %s\n", req.Method, req.URL.String())
	if req.Header != nil {
		for k, vals := range req.Header {
			for _, v := range vals {
				// Mask sensitive headers
				if isSensitiveHeader(k) {
					fmt.Fprintf(t.Out, "[DEBUG] %s: ***\n", k)
				} else {
					fmt.Fprintf(t.Out, "[DEBUG] %s: %s\n", k, v)
				}
			}
		}
	}

	// Record start time for timing info
	startTime := time.Now()

	// Execute the request
	resp, err := t.base().RoundTrip(req)

	// Log response
	elapsed := time.Since(startTime)
	if err != nil {
		fmt.Fprintf(t.Out, "[DEBUG] error: %v (took %v)\n", err, elapsed)
		return resp, err
	}

	if resp != nil {
		fmt.Fprintf(t.Out, "[DEBUG] %d %s (took %v)\n", resp.StatusCode, resp.Status, elapsed)
		if resp.Header != nil {
			for k, vals := range resp.Header {
				for _, v := range vals {
					if isSensitiveHeader(k) {
						fmt.Fprintf(t.Out, "[DEBUG] %s: ***\n", k)
					} else {
						fmt.Fprintf(t.Out, "[DEBUG] %s: %s\n", k, v)
					}
				}
			}
		}
	}

	return resp, err
}

// isSensitiveHeader returns true for headers that should be masked in debug output.
func isSensitiveHeader(name string) bool {
	switch name {
	case "Authorization", "X-Lark-Request-Timestamp", "X-Lark-Content-V2", "Cookie", "Set-Cookie":
		return true
	}
	return false
}
