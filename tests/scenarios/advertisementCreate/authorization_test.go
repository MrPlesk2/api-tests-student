package advertisementCreate

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

const (
	invalidToken      = "incorrect_token"
	unauthorizedError = "unauthorized"
)

// TestCreateAdvertisementWithoutAuthorization distinguishes absent and invalid bearer tokens.
func TestCreateAdvertisementWithoutAuthorization(t *testing.T) {
	// Environment initialization does not require valid user credentials in these scenarios.
	cases := []struct {
		name    string
		token   string
		message string
	}{
		{name: "missing_token", message: "Authorization header required"},
		{name: "invalid_token", token: invalidToken, message: "Invalid or expired token"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("[Precondition] Load API configuration and prepare a valid advertisement")
			fixtures.LoadEnvironment(t)
			request := fixtures.AdvertisementRequest(t)
			t.Logf("[Step] Create an advertisement using %s", tc.name)
			response := advertisements.CreateAdvertisementExpectError(t, tc.token, request, http.StatusUnauthorized)
			require.Equal(t, unauthorizedError, response.Error, "Wrong authorization error type")
			require.Equal(t, tc.message, response.Message, "Wrong authorization error message")
			t.Logf("[Result] HTTP 401, error=%s, message=%q", response.Error, response.Message)
		})
	}
}
