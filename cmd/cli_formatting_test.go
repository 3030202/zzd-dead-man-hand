package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestListActionsFormatting(t *testing.T) {
	tests := []struct {
		name           string
		mockTerminal   bool
		responseBody   string
		expectedOutput string
	}{
		{
			name:         "Terminal - Pretty Print",
			mockTerminal: true,
			responseBody: `[{"kind":"email","data":"test"}]`,
			expectedOutput: `[
  {
    "kind": "email",
    "data": "test"
  }
]
`,
		},
		{
			name:         "Non-Terminal - Raw Output",
			mockTerminal: false,
			responseBody: `[{"kind":"email","data":"test"}]`,
			expectedOutput: `[{"kind":"email","data":"test"}]`,
		},
		{
			name:         "Terminal - Invalid JSON fallback",
			mockTerminal: true,
			responseBody: `invalid-json`,
			expectedOutput: `invalid-json`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mock server
			fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(test.responseBody))
			}))
			defer fakeServer.Close()

			// Mock getClient
			originalGetClient := getClient
			defer func() { getClient = originalGetClient }()
			getClient = func(*cli.Command) *http.Client {
				return fakeServer.Client()
			}

			// Mock outputWriter
			var outBuf bytes.Buffer
			originalOutputWriter := outputWriter
			outputWriter = &outBuf
			defer func() { outputWriter = originalOutputWriter }()

			// Mock isTerminal
			originalIsTerminal := isTerminal
			isTerminal = func() bool { return test.mockTerminal }
			defer func() { isTerminal = originalIsTerminal }()

			cmd := createCLI()
			params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}

			err := cmd.Run(context.Background(), params)
			require.Nil(t, err)

			require.Equal(t, test.expectedOutput, outBuf.String())
		})
	}
}
