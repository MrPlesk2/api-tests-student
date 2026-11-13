package favorites

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	client "api-tests-template/internal/client/http/favorites"
	apiModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"
)

// MaxFavorites is the per-user business limit documented for POST /favorites.
const MaxFavorites = 50

const favoritesPageSize = 20

// GetAllFavorites follows page-based pagination and checks that no entries were lost.
func GetAllFavorites(t *testing.T, token string) []models.Favorite {
	t.Helper()
	var items []models.Favorite
	for page := 1; ; page++ {
		response := utils.DecodeResponse[models.FavoritesList](t,
			client.HttpGetFavorites(t, token, page, favoritesPageSize), http.StatusOK)
		require.Equal(t, page, response.Page, "Favorites returned the wrong page")
		require.Equal(t, favoritesPageSize, response.PageSize, "Favorites returned the wrong page size")
		require.GreaterOrEqual(t, response.TotalPages, 0, "Negative number of favorites pages")
		items = append(items, response.Items...)
		if page >= response.TotalPages {
			require.Len(t, items, response.Total, "Favorites pagination lost or duplicated entries")
			return items
		}
		require.NotEmpty(t, response.Items, "Empty favorites page before the last page")
	}
}

// AddFavorite saves an advertisement and registers removal of that favorite during cleanup.
func AddFavorite(t *testing.T, token, advertisementID string) models.Favorite {
	t.Helper()
	body := postFavorite(t, token, models.AddFavoriteRequest{AdvertisementID: advertisementID}, http.StatusCreated)
	var result models.Favorite
	require.NoError(t, json.Unmarshal(body, &result), "Invalid favorite JSON: %s", body)
	require.NotEmpty(t, result.FavoriteID, "Created favorite has no favorite_id")
	require.NotEmpty(t, result.AddedAt, "Created favorite has no added_at")
	require.Equal(t, advertisementID, result.Advertisement.ID, "POST saved another advertisement")
	return result
}

// AddFavoriteExpectError checks a rejected POST and returns its structured error.
func AddFavoriteExpectError(t *testing.T, token string, request any, expectedStatus int) apiModels.ErrorResponse {
	t.Helper()
	body := postFavorite(t, token, request, expectedStatus)
	var result apiModels.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result), "Invalid favorite error JSON: %s", body)
	return result
}

// postFavorite registers cleanup before assertions, including an unexpected HTTP 201.
func postFavorite(t *testing.T, token string, request any, expectedStatus int) []byte {
	t.Helper()
	response := client.HttpPostFavorite(t, token, request)
	body := utils.ReadResponseBody(t, response)
	if response.StatusCode == http.StatusCreated {
		var identity struct {
			Advertisement struct {
				ID string `json:"id"`
			} `json:"advertisement"`
		}
		require.NoError(t, json.Unmarshal(body, &identity), "Invalid favorite response: %s", body)
		if identity.Advertisement.ID != "" {
			t.Cleanup(func() {
				t.Logf("[Cleanup] Remove favorite for advertisement %s", identity.Advertisement.ID)
				RemoveFavoriteIfExists(t, token, identity.Advertisement.ID)
			})
		}
	}
	require.Equal(t, expectedStatus, response.StatusCode, "Unexpected POST /favorites status; response: %s", body)
	return body
}

// RemoveFavorite requires successful removal of an existing favorite.
func RemoveFavorite(t *testing.T, token, advertisementID string) {
	t.Helper()
	response := client.HttpDeleteFavorite(t, token, advertisementID)
	body := utils.ReadResponseBody(t, response)
	require.Equal(t, http.StatusNoContent, response.StatusCode, "Cannot remove favorite %s; response: %s", advertisementID, body)
}

// RemoveFavoriteIfExists allows cleanup after explicit removal or advertisement cascade deletion.
func RemoveFavoriteIfExists(t *testing.T, token, advertisementID string) {
	t.Helper()
	response := client.HttpDeleteFavorite(t, token, advertisementID)
	body := utils.ReadResponseBody(t, response)
	require.Contains(t, []int{http.StatusNoContent, http.StatusNotFound}, response.StatusCode,
		"Cannot clean up favorite %s; response: %s", advertisementID, body)
}

// RemoveFavoriteExpectError checks the exact status of a rejected DELETE.
func RemoveFavoriteExpectError(t *testing.T, token, advertisementID string, expectedStatus int) apiModels.ErrorResponse {
	t.Helper()
	return utils.DecodeResponse[apiModels.ErrorResponse](t,
		client.HttpDeleteFavorite(t, token, advertisementID), expectedStatus)
}

// GetFavoritesExpectError checks the exact status of a rejected GET.
func GetFavoritesExpectError(t *testing.T, token string, expectedStatus int) apiModels.ErrorResponse {
	t.Helper()
	return utils.DecodeResponse[apiModels.ErrorResponse](t,
		client.HttpGetFavorites(t, token, 1, favoritesPageSize), expectedStatus)
}

// MergeFavorites adds a batch of test-owned advertisements whose fixtures handle cascade cleanup.
func MergeFavorites(t *testing.T, token string, advertisementIDs []string) models.BulkAddFavoritesResponse {
	t.Helper()
	return utils.DecodeResponse[models.BulkAddFavoritesResponse](t,
		client.HttpPostFavoritesBulk(t, token, models.BulkAddFavoritesRequest{AdvertisementIDs: advertisementIDs}), http.StatusOK)
}

// MergeFavoritesExpectError checks a rejected bulk request, including authentication failures.
func MergeFavoritesExpectError(t *testing.T, token string, advertisementIDs []string, expectedStatus int) apiModels.ErrorResponse {
	t.Helper()
	return utils.DecodeResponse[apiModels.ErrorResponse](t,
		client.HttpPostFavoritesBulk(t, token, models.BulkAddFavoritesRequest{AdvertisementIDs: advertisementIDs}), expectedStatus)
}
