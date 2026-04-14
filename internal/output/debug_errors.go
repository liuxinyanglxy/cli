// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package output

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime"
)

// WriteDebugStackTrace writes error information with stack trace to the given writer.
// Used in debug mode to provide detailed error diagnostics.
func WriteDebugStackTrace(w io.Writer, err error) {
	if err == nil {
		return
	}

	// Write the error message
	fmt.Fprintf(w, "[DEBUG] Error: %v\n", err)

	// Try to get stack trace information
	// For unwrapped errors, we'll show the error chain
	fmt.Fprintf(w, "[DEBUG] Stack trace:\n")

	// Print the error chain
	var errChain []error
	current := err
	for current != nil {
		errChain = append(errChain, current)
		current = errors.Unwrap(current)
	}

	for i, e := range errChain {
		fmt.Fprintf(w, "[DEBUG]   %d: %v (%T)\n", i, e, e)
	}

	// Print current goroutine stack
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	if n > 0 {
		fmt.Fprintf(w, "[DEBUG] Goroutine stack:\n")
		for _, line := range bytes.Split(buf[:n], []byte("\n")) {
			if len(line) > 0 {
				fmt.Fprintf(w, "[DEBUG]   %s\n", line)
			}
		}
	}
}
