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

// TestFavoritesLimit rejects the 51st entry and allows adding again after one slot is freed.
func TestFavoritesLimit(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)
	t.Logf("[Precondition] Create %d foreign advertisements with automatic cleanup", manager.MaxFavorites+1)
	advertisementIDs := make([]string, 0, manager.MaxFavorites+1)
	for i := 0; i <= manager.MaxFavorites; i++ {
		advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
		advertisementIDs = append(advertisementIDs, advertisement.ID)
	}
	expectedIDs := advertisementIDs[:manager.MaxFavorites]
	extraID := advertisementIDs[manager.MaxFavorites]

	t.Log("[Step] Fill favorites using the documented bulk endpoint")
	merged := manager.MergeFavorites(t, viewer.Token, expectedIDs)
	require.ElementsMatch(t, expectedIDs, merged.Added, "Bulk setup did not add all 50 foreign advertisements")
	require.Empty(t, merged.Skipped, "Bulk setup skipped new advertisements")
	require.Empty(t, merged.NotFound, "Bulk setup lost existing advertisements")
	require.False(t, merged.LimitReached, "Empty account reported that the limit was already reached")

	before := manager.GetAllFavorites(t, viewer.Token)
	require.Len(t, before, manager.MaxFavorites, "Favorites were not filled to the limit")
	actualIDs := make([]string, 0, len(before))
	for _, favorite := range before {
		actualIDs = append(actualIDs, favorite.Advertisement.ID)
	}
	require.ElementsMatch(t, expectedIDs, actualIDs, "Favorites pagination changed the 50 prepared entries")

	t.Log("[Step] Add the 51st existing foreign advertisement using POST /favorites")
	response := manager.AddFavoriteExpectError(t, viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: extraID}, http.StatusBadRequest)
	t.Logf("[Check] Limit error: error=%q, message=%q", response.Error, response.Message)
	require.Equal(t, "limit_reached", response.Error, "Wrong favorites limit error type")
	require.Equal(t, "Favorites limit of 50 reached", response.Message, "Wrong favorites limit error message")
	require.Equal(t, before, manager.GetAllFavorites(t, viewer.Token), "Limit rejection changed existing favorites")
	t.Log("[Result] HTTP 400 rejected the 51st entry and preserved the original 50 favorites")

	t.Log("[Step] Remove one favorite to free a slot at the limit")
	removedID := expectedIDs[0]
	manager.RemoveFavorite(t, viewer.Token, removedID)
	remaining := make([]models.Favorite, 0, manager.MaxFavorites-1)
	for _, favorite := range before {
		if favorite.Advertisement.ID != removedID {
			remaining = append(remaining, favorite)
		}
	}
	require.Len(t, remaining, manager.MaxFavorites-1, "Exactly one prepared favorite must be removed")
	require.ElementsMatch(t, remaining, manager.GetAllFavorites(t, viewer.Token), "Freeing a slot changed other favorites")
	t.Log("[Step] Add the previously rejected advertisement into the freed slot")
	added := manager.AddFavorite(t, viewer.Token, extraID)
	expected := append(remaining, added)
	actual := manager.GetAllFavorites(t, viewer.Token)
	require.Len(t, actual, manager.MaxFavorites, "Adding into the freed slot must restore exactly 50 favorites")
	require.ElementsMatch(t, expected, actual, "Adding into the freed slot lost, duplicated or changed favorites")
	t.Log("[Result] POST successfully added the 50th favorite after a slot was freed")
}
