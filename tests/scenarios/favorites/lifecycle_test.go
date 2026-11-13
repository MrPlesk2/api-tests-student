package favorites

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	manager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/tests/fixtures"
)

// TestFavoritesLifecycle checks add, remove, re-add and deletion of the underlying advertisement.
func TestFavoritesLifecycle(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)
	request := fixtures.AdvertisementRequest(t)
	t.Log("[Precondition] Create a foreign advertisement for the favorites user")
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, request)

	t.Log("[Step] Add the advertisement and find it in GET /favorites")
	added := manager.AddFavorite(t, viewer.Token, advertisement.ID)
	items := manager.GetAllFavorites(t, viewer.Token)
	require.Len(t, items, 1, "Added advertisement is missing or duplicated")
	require.Equal(t, added.FavoriteID, items[0].FavoriteID, "GET returned another favorite")
	require.Equal(t, added.AddedAt, items[0].AddedAt, "GET changed the favorite timestamp")
	require.Equal(t, advertisement.ID, items[0].Advertisement.ID, "GET returned another advertisement")
	assertions.AdvertisementMatchesRequest(t, items[0].Advertisement, request, owner.User.Id)
	t.Log("[Check] Created favorite and advertisement data are present in the list")

	original := items[0]
	t.Log("[Precondition] Save a second advertisement that must survive both deletion operations")
	control := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
	controlFavorite := manager.AddFavorite(t, viewer.Token, control.ID)
	require.ElementsMatch(t, []models.Favorite{original, controlFavorite}, manager.GetAllFavorites(t, viewer.Token),
		"Both favorites must be present before deletion")

	t.Log("[Step] Remove the advertisement from favorites")
	manager.RemoveFavorite(t, viewer.Token, advertisement.ID)
	require.Equal(t, []models.Favorite{controlFavorite}, manager.GetAllFavorites(t, viewer.Token),
		"Removing one favorite must preserve the other favorite unchanged")
	t.Log("[Check] Explicit favorite removal changed the list")

	t.Log("[Step] Add the same advertisement again")
	readded := manager.AddFavorite(t, viewer.Token, advertisement.ID)
	items = manager.GetAllFavorites(t, viewer.Token)
	require.ElementsMatch(t, []models.Favorite{readded, controlFavorite}, items,
		"Re-adding a favorite must restore it and preserve the other favorite")

	t.Log("[Step] Delete the advertisement as its owner and check cascade removal")
	advertisements.DeleteAdvertisement(t, owner.Token, advertisement.ID)
	require.Equal(t, []models.Favorite{controlFavorite}, manager.GetAllFavorites(t, viewer.Token),
		"Deleting an advertisement must remove only its favorite and preserve the other favorite unchanged")
	t.Log("[Result] Add, remove, re-add and cascade deletion all updated GET /favorites")
}
