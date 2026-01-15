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

func TestListActionsFormatting(t *testing.T) {
	// Setup capture of stdout
	captureOutput := func(f func()) string {
		r, w, _ := os.Pipe()
		stdout := os.Stdout
		os.Stdout = w
		defer func() { os.Stdout = stdout }()

		f()
		w.Close()
		out, _ := io.ReadAll(r)
		return string(out)
	}

	tests := []struct {
		name          string
		isTerminal    bool
		responseBody  string
		expectedOutput string
	}{
		{
			name:          "Terminal: Pretty Print",
			isTerminal:    true,
			responseBody:  `[{"id":"1","kind":"email"}]`,
			expectedOutput: "[\n  {\n    \"id\": \"1\",\n    \"kind\": \"email\"\n  }\n]\n",
		},
		{
			name:          "Non-Terminal: Raw Output",
			isTerminal:    false,
			responseBody:  `[{"id":"1","kind":"email"}]`,
			expectedOutput: `[{"id":"1","kind":"email"}]`,
		},
		{
			name:          "Terminal: Invalid JSON",
			isTerminal:    true,
			responseBody:  `not json`,
			expectedOutput: `not json`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mock server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(test.responseBody))
			}))
			defer ts.Close()

			// Mock isTerminal
			originalIsTerminal := isTerminal
			defer func() { isTerminal = originalIsTerminal }()
			isTerminal = func(w *os.File) bool {
				return test.isTerminal
			}

			// Mock getClient
			originalGetClient := getClient
			defer func() { getClient = originalGetClient }()
			getClient = func(cmd *cli.Command) *http.Client {
				return ts.Client()
			}

			cmd := createCLI()
			params := []string{"dmh-cli", "action", "list", "--server", ts.URL}

			output := captureOutput(func() {
				err := cmd.Run(context.Background(), params)
				require.NoError(t, err)
			})

			require.Equal(t, test.expectedOutput, output)
		})
	}
}
