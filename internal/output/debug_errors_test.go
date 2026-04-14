// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package output

import (
	"bytes"
	"errors"
	"testing"
)

func TestWriteDebugStackTrace_WithError(t *testing.T) {
	buf := &bytes.Buffer{}
	err := errors.New("test error")

	WriteDebugStackTrace(buf, err)

	output := buf.String()
	// Should contain the error message
	if !bytes.Contains([]byte(output), []byte("test error")) {
		t.Errorf("output should contain error message, got: %s", output)
	}
	// Should indicate it's debug output
	if !bytes.Contains([]byte(output), []byte("[DEBUG]")) {
		t.Errorf("output should contain [DEBUG] prefix, got: %s", output)
	}
}

func TestWriteDebugStackTrace_NilError(t *testing.T) {
	buf := &bytes.Buffer{}
	WriteDebugStackTrace(buf, nil)

	output := buf.String()
	// Should handle nil gracefully
	if len(output) > 0 && !bytes.Contains([]byte(output), []byte("[DEBUG]")) {
		t.Errorf("output should not contain data for nil error, got: %s", output)
	}
}
