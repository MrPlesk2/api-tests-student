package favorites

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// HttpGetFavorites reads one page using a stable order for pagination.
func HttpGetFavorites(t *testing.T, token string, page, pageSize int) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("page", strconv.Itoa(page)).
		Query("page_size", strconv.Itoa(pageSize)).
		Query("sort_by", "title").
		Query("sort_order", "asc").
		Expect(t).End().Response
}

// HttpPostFavorite sends the supplied JSON, including invalid payloads used by negative tests.
func HttpPostFavorite(t *testing.T, token string, request any) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").JSON(request).
		Expect(t).End().Response
}

// HttpDeleteFavorite removes a favorite by advertisement ID.
func HttpDeleteFavorite(t *testing.T, token, advertisementID string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.FavoritesPath + "/" + advertisementID).
		Expect(t).End().Response
}

// HttpPostFavoritesBulk merges a batch of advertisement IDs into the user's favorites.
func HttpPostFavoritesBulk(t *testing.T, token string, request models.BulkAddFavoritesRequest) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesBulkPath).
		ContentType("application/json").JSON(request).
		Expect(t).End().Response
}
