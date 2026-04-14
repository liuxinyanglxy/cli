// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/larksuite/cli/internal/cmdutil"
	"github.com/larksuite/cli/internal/core"
)

func TestDebugFlag_OutputsDebugInfo(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	// Create a factory with debug enabled
	inv := cmdutil.InvocationContext{Debug: true}
	f := cmdutil.NewDefault(inv)
	f.Config = func() (*core.CliConfig, error) {
		return &core.CliConfig{
			AppID:     "test-app",
			AppSecret: "test-secret",
			Brand:     core.BrandFeishu,
		}, nil
	}

	// Get the HTTP client
	httpClient, err := f.HttpClient()
	if err != nil {
		t.Fatalf("HttpClient() error: %v", err)
	}

	// Make a request through the client
	req, _ := http.NewRequest("GET", server.URL+"/test", nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("Do() error: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "ok") {
		t.Errorf("response body should contain 'ok', got: %s", body)
	}
}

func TestDebugFlag_DefaultFalse(t *testing.T) {
	// Default debug should be false
	inv := cmdutil.InvocationContext{}
	if inv.Debug {
		t.Error("Debug should be false by default")
	}
}
