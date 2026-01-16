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

// Helper to capture stdout
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
	// Mock server returning JSON
	mockResponse := `[{"kind":"test","data":"foo"}]`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Mock client
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return server.Client()
	}

	// Mock isTerminal to force pretty printing
	originalIsTerminal := isTerminal
	defer func() { isTerminal = originalIsTerminal }()
	isTerminal = func(f *os.File) bool { return true }

	// Prepare command
	cmd := createCLI()
	args := []string{"dmh-client", "action", "list", "--server", server.URL}

	// Capture output
	output, err := captureOutput(func() error {
		return cmd.Run(context.Background(), args)
	})

	require.Nil(t, err)
	// Verify pretty printed output
	expectedOutput := "[\n  {\n    \"kind\": \"test\",\n    \"data\": \"foo\"\n  }\n]\n"
	require.Equal(t, expectedOutput, output)
}

func TestListActionsRaw(t *testing.T) {
	// Mock server returning JSON
	mockResponse := `[{"kind":"test","data":"foo"}]`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Mock client
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return server.Client()
	}

	// Mock isTerminal to force raw printing (non-terminal)
	originalIsTerminal := isTerminal
	defer func() { isTerminal = originalIsTerminal }()
	isTerminal = func(f *os.File) bool { return false }

	// Prepare command
	cmd := createCLI()
	args := []string{"dmh-client", "action", "list", "--server", server.URL}

	// Capture output
	output, err := captureOutput(func() error {
		return cmd.Run(context.Background(), args)
	})

	require.Nil(t, err)
	// Verify raw output
	require.Equal(t, mockResponse, output)
}
