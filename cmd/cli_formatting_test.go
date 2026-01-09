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

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestListActionsOutputFormat(t *testing.T) {
	// Sample JSON that we expect to be pretty-printed
	rawJSON := `[{"uuid":"123","kind":"email"}]`

	// Expected pretty-printed output
	expectedOutput := "[\n  {\n    \"kind\": \"email\",\n    \"uuid\": \"123\"\n  }\n]\n"

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rawJSON))
	}

	fakeServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer fakeServer.Close()

	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	t.Run("Default Pretty Print", func(t *testing.T) {
		cmd := createCLI()
		params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}

		output := captureOutput(func() {
			err := cmd.Run(context.Background(), params)
			require.Nil(t, err)
		})

		require.Equal(t, expectedOutput, output)
	})

	t.Run("Raw Output", func(t *testing.T) {
		cmd := createCLI()
		params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL, "--raw"}

		output := captureOutput(func() {
			err := cmd.Run(context.Background(), params)
			require.Nil(t, err)
		})

		require.Equal(t, rawJSON, output)
	})
}
