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

// captureOutput captures stdout output of a function
func captureOutput(f func()) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}

	oldStdout := os.Stdout
	os.Stdout = w

	outC := make(chan string)

	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	f()

	// close the writer, restore stdout
	w.Close()
	os.Stdout = oldStdout

	return <-outC, nil
}

func TestListActionsOutputFormat(t *testing.T) {
	// Raw JSON response
	rawJSON := `[{"id":"123","kind":"email","data":"test"}]`

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rawJSON))
	}

	fakeServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer fakeServer.Close()

	// Mock the client
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	cmd := createCLI()
	params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}

	output, err := captureOutput(func() {
		err := cmd.Run(context.Background(), params)
		require.Nil(t, err)
	})

	require.Nil(t, err)

	// Expected pretty output
	expectedPretty := `[
  {
    "data": "test",
    "id": "123",
    "kind": "email"
  }
]
`
	require.Equal(t, expectedPretty, output)
}
