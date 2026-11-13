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

// TestFavoritesNonexistentAdvertisement uses a deleted advertisement to guarantee that its valid UUID is absent.
func TestFavoritesNonexistentAdvertisement(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)
	t.Log("[Precondition] Create and delete an advertisement to obtain a guaranteed nonexistent ID")
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
	advertisements.DeleteAdvertisement(t, owner.Token, advertisement.ID)

	t.Run("add", func(t *testing.T) {
		t.Log("[Step] Add the deleted advertisement to favorites")
		response := manager.AddFavoriteExpectError(t, viewer.Token,
			models.AddFavoriteRequest{AdvertisementID: advertisement.ID}, http.StatusNotFound)
		require.Equal(t, "not_found", response.Error, "Wrong missing-advertisement error type")
		require.Equal(t, "Advertisement not found", response.Message, "Wrong missing-advertisement message")
		require.Empty(t, manager.GetAllFavorites(t, viewer.Token), "Rejected POST changed favorites")
		t.Log("[Result] Nonexistent advertisement was rejected with HTTP 404")
	})
	t.Run("remove", func(t *testing.T) {
		t.Log("[Step] Remove a nonexistent advertisement from favorites")
		response := manager.RemoveFavoriteExpectError(t, viewer.Token, advertisement.ID, http.StatusNotFound)
		require.Equal(t, "not_found", response.Error, "Wrong missing-favorite error type")
		require.Equal(t, "Favorite entry not found", response.Message, "Wrong missing-favorite message")
		require.Empty(t, manager.GetAllFavorites(t, viewer.Token), "Rejected DELETE changed favorites")
		t.Log("[Result] Nonexistent favorite was rejected with HTTP 404")
	})
}
