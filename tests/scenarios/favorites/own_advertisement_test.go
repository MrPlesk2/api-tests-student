package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

// TestFavoritesOwnAdvertisement checks that a user cannot favorite their own advertisement.
func TestFavoritesOwnAdvertisement(t *testing.T) {
	_, viewer := fixtures.FavoritesUsers(t)
	t.Log("[Precondition] Create an advertisement owned by the favorites user")
	advertisement := advertisements.CreateAdvertisement(t, viewer.Token, fixtures.AdvertisementRequest(t))

	t.Log("[Step] Add the user's own advertisement to favorites")
	response := manager.AddFavoriteExpectError(t, viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: advertisement.ID}, http.StatusForbidden)
	require.Equal(t, "forbidden", response.Error, "Wrong ownership error type")
	require.Equal(t, "Cannot add own advertisement to favorites", response.Message, "Wrong ownership error message")
	require.Empty(t, manager.GetAllFavorites(t, viewer.Token), "Rejected own advertisement appeared in favorites")
	t.Log("[Result] Own advertisement was rejected with HTTP 403; favorites remain empty")
}
