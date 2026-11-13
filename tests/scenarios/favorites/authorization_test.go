package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	apiModels "api-tests-template/internal/managers/advertisements/models"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

const invalidToken = "incorrect_token"

// TestFavoritesWithoutAuthorization checks missing and invalid tokens for every documented favorites operation.
func TestFavoritesWithoutAuthorization(t *testing.T) {
	owner := fixtures.Login(t)
	t.Log("[Precondition] Create a real advertisement to keep all request IDs valid")
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
	operations := []struct {
		name string
		call func(*testing.T, string) apiModels.ErrorResponse
	}{
		{name: "get", call: func(t *testing.T, token string) apiModels.ErrorResponse {
			return manager.GetFavoritesExpectError(t, token, http.StatusUnauthorized)
		}},
		{name: "add", call: func(t *testing.T, token string) apiModels.ErrorResponse {
			return manager.AddFavoriteExpectError(t, token,
				models.AddFavoriteRequest{AdvertisementID: advertisement.ID}, http.StatusUnauthorized)
		}},
		{name: "remove", call: func(t *testing.T, token string) apiModels.ErrorResponse {
			return manager.RemoveFavoriteExpectError(t, token, advertisement.ID, http.StatusUnauthorized)
		}},
		{name: "bulk_add", call: func(t *testing.T, token string) apiModels.ErrorResponse {
			return manager.MergeFavoritesExpectError(t, token, []string{advertisement.ID}, http.StatusUnauthorized)
		}},
	}
	tokens := []struct {
		name    string
		token   string
		message string
	}{
		{name: "missing_token", message: "Authorization header required"},
		{name: "invalid_token", token: invalidToken, message: "Invalid or expired token"},
	}
	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			for _, token := range tokens {
				t.Run(token.name, func(t *testing.T) {
					t.Logf("[Step] Call favorites operation %s with %s", operation.name, token.name)
					response := operation.call(t, token.token)
					require.Equal(t, "unauthorized", response.Error, "Wrong authentication error type")
					require.Equal(t, token.message, response.Message, "Wrong authentication error message")
					t.Log("[Result] HTTP 401 and the exact authentication message were returned")
				})
			}
		})
	}
}
