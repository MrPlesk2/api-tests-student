package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ReadResponseBody reads and closes a response, logging its HTTP status.
func ReadResponseBody(t *testing.T, response *http.Response) []byte {
	t.Helper()
	require.NotNil(t, response, "API did not return an HTTP response")
	require.NotNil(t, response.Body, "HTTP response has no body")
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Cannot read HTTP response")
	t.Logf("[HTTP] status=%d", response.StatusCode)
	return body
}

// DecodeResponse checks the HTTP status and decodes a JSON response into the requested model.
func DecodeResponse[T any](t *testing.T, response *http.Response, expectedStatus int) T {
	t.Helper()
	body := ReadResponseBody(t, response)
	require.Equal(t, expectedStatus, response.StatusCode,
		"Unexpected HTTP status; response: %s", body)

	var result T
	require.NoError(t, json.Unmarshal(body, &result), "Invalid JSON response: %s", body)
	return result
}
