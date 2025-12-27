package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestListActionsPrettyPrint(t *testing.T) {
	// Mock server returning minified JSON
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":"1","kind":"email","data":"..."}]`))
	}

	fakeServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer fakeServer.Close()

	// Mock client
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := createCLI()
	params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}

	err := cmd.Run(context.Background(), params)

	w.Close()
	os.Stdout = oldStdout

	require.Nil(t, err)

	out, _ := io.ReadAll(r)
	output := string(out)

	// Check if output is pretty printed (contains newlines and indentation)
	// The raw output is `[{"id":"1","kind":"email","data":"..."}]` which has no newlines (except maybe a trailing one from io.Copy if it was there, but it wasn't).
	// We expect indentation.

	require.Contains(t, output, "[\n", "Output should start with bracket and newline")
	require.Contains(t, output, "  {", "Output should be indented")
}
