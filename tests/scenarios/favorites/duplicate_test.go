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

// TestFavoritesDuplicateAdvertisement checks HTTP 409 and verifies that the original entry is preserved.
func TestFavoritesDuplicateAdvertisement(t *testing.T) {
	owner, viewer := fixtures.FavoritesUsers(t)
	t.Log("[Precondition] Create a foreign advertisement and save it to favorites")
	advertisement := advertisements.CreateAdvertisement(t, owner.Token, fixtures.AdvertisementRequest(t))
	saved := manager.AddFavorite(t, viewer.Token, advertisement.ID)

	t.Log("[Step] Add the same advertisement a second time")
	response := manager.AddFavoriteExpectError(t, viewer.Token,
		models.AddFavoriteRequest{AdvertisementID: advertisement.ID}, http.StatusConflict)
	require.Equal(t, "already_exists", response.Error, "Wrong duplicate error type")
	require.Equal(t, "Advertisement is already in favorites", response.Message, "Wrong duplicate error message")

	items := manager.GetAllFavorites(t, viewer.Token)
	require.Len(t, items, 1, "Repeated POST duplicated or removed the favorite")
	require.Equal(t, saved, items[0], "Repeated POST modified the original favorite")
	t.Log("[Result] Duplicate was rejected with HTTP 409; the original favorite is unchanged")
}
