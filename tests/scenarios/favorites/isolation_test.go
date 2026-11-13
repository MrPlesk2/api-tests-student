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

// TestFavoritesUserIsolation checks independent favorites for two users saving the same foreign advertisement.
func TestFavoritesUserIsolation(t *testing.T) {
	owner, first := fixtures.FavoritesUsers(t)
	t.Run("second_user", func(t *testing.T) {
		_, second := fixtures.FavoritesUsers(t)
		require.NotEqual(t, first.User.Id, second.User.Id, "Favorites viewers must be different users")
		t.Log("[Precondition] Create a foreign advertisement and save it for the first user")
		advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
		firstFavorite := manager.AddFavorite(t, first.Token, advertisement.ID)
		require.Equal(t, []models.Favorite{firstFavorite}, manager.GetAllFavorites(t, first.Token), "First user must see their favorite")
		require.Empty(t, manager.GetAllFavorites(t, second.Token), "Second user can see another user's favorite")
		t.Log("[Step] Try removing the first user's favorite as the second user")
		response := manager.RemoveFavoriteExpectError(t, second.Token, advertisement.ID, http.StatusNotFound)
		require.Equal(t, "not_found", response.Error, "Wrong cross-user removal error type")
		require.Equal(t, "Favorite entry not found", response.Message, "Wrong cross-user removal message")
		require.Equal(t, []models.Favorite{firstFavorite}, manager.GetAllFavorites(t, first.Token), "Second user removed or changed the first user's favorite")
		require.Empty(t, manager.GetAllFavorites(t, second.Token), "Rejected removal changed the second user's favorites")
		t.Log("[Step] Save the same advertisement independently for the second user")
		secondFavorite := manager.AddFavorite(t, second.Token, advertisement.ID)
		require.NotEqual(t, firstFavorite.FavoriteID, secondFavorite.FavoriteID, "Different users must have different favorite records")
		require.Equal(t, []models.Favorite{firstFavorite}, manager.GetAllFavorites(t, first.Token), "Second user's POST changed the first user's favorite")
		require.Equal(t, []models.Favorite{secondFavorite}, manager.GetAllFavorites(t, second.Token), "Second user must see their own favorite")
		t.Log("[Step] Remove the first user's favorite and preserve the second user's record")
		manager.RemoveFavorite(t, first.Token, advertisement.ID)
		require.Empty(t, manager.GetAllFavorites(t, first.Token), "First user's favorite was not removed")
		require.Equal(t, []models.Favorite{secondFavorite}, manager.GetAllFavorites(t, second.Token), "First user's removal changed the second user's favorite")
		t.Log("[Result] Favorites are visible and removable only within the authenticated user's account")
	})
}
