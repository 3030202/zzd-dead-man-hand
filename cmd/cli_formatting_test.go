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
		responseBody   string
		expectedOutput string
	}{
		{
			name:         "Pretty print JSON",
			responseBody: `[{"id":"1","data":"test"}]`,
			expectedOutput: `[
  {
    "id": "1",
    "data": "test"
  }
]
`,
		},
		{
			name:           "Raw output for non-JSON",
			responseBody:   `invalid json`,
			expectedOutput: `invalid json`,
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

			// Mock getClient
			originalGetClient := getClient
			defer func() { getClient = originalGetClient }()
			getClient = func(*cli.Command) *http.Client {
				return ts.Client()
			}

			// Mock outputWriter
			var buf bytes.Buffer
			originalOutputWriter := outputWriter
			outputWriter = &buf
			defer func() { outputWriter = originalOutputWriter }()

			// Run command
			cmd := createCLI()
			err := cmd.Run(context.Background(), []string{"dmh-client", "action", "list", "--server", ts.URL})

			require.NoError(t, err)
			require.Equal(t, test.expectedOutput, buf.String())
		})
	}
}

func TestActionSuccessMessages(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		mockHandler    http.HandlerFunc
	}{
		{
			name:           "Add action success",
			args:           []string{"action", "add", "--data", "{}", "--kind", "test", "--process-after", "10"},
			expectedOutput: "Action added successfully\n",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
			},
		},
		{
			name:           "Test action success",
			args:           []string{"action", "test", "--data", "{}", "--kind", "test"},
			expectedOutput: "Action tested successfully\n",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			name:           "Delete action success",
			args:           []string{"action", "delete", "--uuid", "123"},
			expectedOutput: "Action deleted successfully\n",
			mockHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mock server
			ts := httptest.NewServer(test.mockHandler)
			defer ts.Close()

			// Mock getClient
			originalGetClient := getClient
			defer func() { getClient = originalGetClient }()
			getClient = func(*cli.Command) *http.Client {
				return ts.Client()
			}

			// Mock outputWriter
			var buf bytes.Buffer
			originalOutputWriter := outputWriter
			outputWriter = &buf
			defer func() { outputWriter = originalOutputWriter }()

			// Run command
			cmd := createCLI()
			fullArgs := append([]string{"dmh-client"}, test.args...)
			fullArgs = append(fullArgs, "--server", ts.URL)
			err := cmd.Run(context.Background(), fullArgs)

			require.NoError(t, err)
			require.Equal(t, test.expectedOutput, buf.String())
		})
	}
}
