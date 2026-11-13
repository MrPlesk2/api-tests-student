package myAdvertisements

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetMyAdvertisements retrieves one page of the authenticated user's advertisements.
func HttpGetMyAdvertisements(t *testing.T, token string, limit int, cursor string) *http.Response {
	t.Helper()
	request := apiRunner.GetRunner().Auth(token).Create().Get(path.MyAdvertisementsPath).
		Query("limit", strconv.Itoa(limit))
	if cursor != "" {
		request.Query("cursor", cursor)
	}
	return request.Expect(t).End().Response
}
