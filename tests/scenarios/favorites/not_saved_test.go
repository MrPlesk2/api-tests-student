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

// TestRemoveUnsavedAdvertisement rejects removal of an existing advertisement absent from favorites.
func TestRemoveUnsavedAdvertisement(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)
	t.Log("[Precondition] Create an unsaved advertisement and another advertisement saved to favorites")
	target := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
	control := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
	saved := manager.AddFavorite(t, viewer.Token, control.ID)
	before := manager.GetAllFavorites(t, viewer.Token)
	require.Equal(t, []models.Favorite{saved}, before, "Precondition must contain only the control favorite")
	t.Log("[Step] Remove an existing advertisement that was never saved")
	response := manager.RemoveFavoriteExpectError(t, viewer.Token, target.ID, http.StatusNotFound)
	require.Equal(t, "not_found", response.Error, "Wrong missing-favorite error type")
	require.Equal(t, "Favorite entry not found", response.Message, "Wrong missing-favorite message")
	require.Equal(t, before, manager.GetAllFavorites(t, viewer.Token), "Rejected removal changed existing favorites")
	require.Equal(t, target.ID, advertisements.GetAdvertisement(t, owner.Token, target.ID).ID, "Rejected removal deleted the advertisement")
	t.Log("[Result] HTTP 404 rejected removal and preserved favorites and the advertisement")
}
