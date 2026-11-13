package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	apiModels "api-tests-template/internal/managers/advertisements/models"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

const (
	invalidAdvertisementID = "not-a-uuid"
	invalidRequestError    = "invalid_request"
)

// TestFavoritesInvalidParameters checks malformed and missing IDs without confusing them with absent advertisements.
func TestFavoritesInvalidParameters(t *testing.T) {
	_, viewer := fixtures.FavoritesUsers(t)
	cases := []struct {
		name    string
		request any
		remove  bool
		message string
	}{
		{name: "missing_id", request: map[string]any{},
			message: "Key: 'AddFavoriteRequest.AdvertisementID' Error:Field validation for 'AdvertisementID' failed on the 'required' tag"},
		{name: "empty_id", request: models.AddFavoriteRequest{}, message: "invalid UUID length: 0"},
		{name: "invalid_id", request: models.AddFavoriteRequest{AdvertisementID: invalidAdvertisementID},
			message: "invalid UUID length: 10"},
		{name: "numeric_id", request: map[string]any{"advertisement_id": 123},
			message: "json: cannot unmarshal number into Go struct field AddFavoriteRequest.advertisement_id of type uuid.UUID"},
		{name: "remove_invalid_id", remove: true, message: "Invalid advertisement ID"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("[Step] Send invalid favorites parameter: %s", tc.name)
			var response apiModels.ErrorResponse
			if tc.remove {
				response = manager.RemoveFavoriteExpectError(t, viewer.Token, invalidAdvertisementID, http.StatusBadRequest)
			} else {
				response = manager.AddFavoriteExpectError(t, viewer.Token, tc.request, http.StatusBadRequest)
			}
			require.Equal(t, invalidRequestError, response.Error, "Wrong validation error type")
			require.Equal(t, tc.message, response.Message, "Wrong validation error message")
			require.Empty(t, manager.GetAllFavorites(t, viewer.Token), "Rejected request changed favorites")
			t.Log("[Result] HTTP 400 and the exact validation message were returned")
		})
	}
}
