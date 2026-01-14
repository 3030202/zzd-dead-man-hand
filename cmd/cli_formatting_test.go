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
	jsonResponse := `[{"kind":"email","data":"{\"to\":\"test@example.com\"}"}]`

	// Mock server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}
	fakeServer := httptest.NewServer(http.HandlerFunc(handler))
	defer fakeServer.Close()

	// Mock client
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	// Test Case 1: Non-terminal (Raw Output)
	// Mock isTerminal to return false
	originalIsTerminal := isTerminal
	isTerminal = func(fd uintptr) bool { return false }
	defer func() { isTerminal = originalIsTerminal }()

	cmd := createCLI()
	params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}

	output, err := captureOutput(func() error {
		return cmd.Run(context.Background(), params)
	})
	require.Nil(t, err)
	require.Equal(t, jsonResponse, output, "Output should be raw JSON when not a terminal")

	// Test Case 2: Terminal (Pretty Output)
	// Mock isTerminal to return true
	isTerminal = func(fd uintptr) bool { return true }

	output, err = captureOutput(func() error {
		return cmd.Run(context.Background(), params)
	})
	require.Nil(t, err)

	// Use explicit requirement now that we implemented it
	require.Contains(t, output, "\n", "Output should contain newlines (pretty printed)")
	require.Contains(t, output, "\t", "Output should contain tabs (pretty printed)")
	require.Greater(t, len(output), len(jsonResponse), "Pretty printed output should be longer than raw")

	// Let's rely on structural check
	require.NotEqual(t, jsonResponse, output)
}
