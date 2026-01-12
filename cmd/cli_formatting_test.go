package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

// captureOutput helps capture stdout for testing
func captureOutput(f func() error) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func TestListActionsFormatting(t *testing.T) {
	// Sample JSON response
	sampleResponse := `[{"id":"1","kind":"email","data":"test@example.com"}]`

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sampleResponse))
	}

	fakeServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer fakeServer.Close()

	// Mock getClient
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	cmd := createCLI()
	params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}

	output, err := captureOutput(func() error {
		return cmd.Run(context.Background(), params)
	})

	require.Nil(t, err)

	// Check if output is pretty printed (contains newlines and indentation)
    // The raw output would be `[{"id":"1","kind":"email","data":"test@example.com"}]`
    // The pretty printed output should span multiple lines.
	require.Contains(t, output, "\n")
	require.Contains(t, output, "  ") // Check for 2-space indentation
}
