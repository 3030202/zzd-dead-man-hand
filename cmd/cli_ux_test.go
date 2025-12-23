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

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	f()

	w.Close()
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestListActionsOutputFormatting(t *testing.T) {
	mockResponse := `[{"id":"1","name":"action1"},{"id":"2","name":"action2"}]`

	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer fakeServer.Close()

	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	cmd := createCLI()
	params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}

	output := captureOutput(func() {
		err := cmd.Run(context.Background(), params)
		require.Nil(t, err)
	})

	// Expected pretty-printed output
	expectedOutput := `[
  {
    "id": "1",
    "name": "action1"
  },
  {
    "id": "2",
    "name": "action2"
  }
]
`
	require.Equal(t, expectedOutput, output)
}
