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

// TestFavoritesRepeatedRemoval checks the second DELETE and preserves an unrelated favorite.
func TestFavoritesRepeatedRemoval(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)
	t.Log("[Precondition] Create and save two foreign advertisements")
	target := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
	control := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
	targetFavorite := manager.AddFavorite(t, viewer.Token, target.ID)
	controlFavorite := manager.AddFavorite(t, viewer.Token, control.ID)
	require.ElementsMatch(t, []models.Favorite{targetFavorite, controlFavorite}, manager.GetAllFavorites(t, viewer.Token), "Both favorites must exist before removal")
	t.Log("[Step] Remove the target favorite once")
	manager.RemoveFavorite(t, viewer.Token, target.ID)
	expected := []models.Favorite{controlFavorite}
	require.Equal(t, expected, manager.GetAllFavorites(t, viewer.Token), "First removal must preserve the control favorite")
	t.Log("[Step] Remove the same favorite a second time")
	response := manager.RemoveFavoriteExpectError(t, viewer.Token, target.ID, http.StatusNotFound)
	require.Equal(t, "not_found", response.Error, "Wrong repeated-removal error type")
	require.Equal(t, "Favorite entry not found", response.Message, "Wrong repeated-removal message")
	require.Equal(t, expected, manager.GetAllFavorites(t, viewer.Token), "Repeated removal changed the other favorite")
	require.Equal(t, target.ID, advertisements.GetAdvertisement(t, owner.Token, target.ID).ID, "Favorite removal deleted the advertisement")
	t.Log("[Result] Second DELETE returned HTTP 404 and preserved the advertisement and the other favorite")
}
