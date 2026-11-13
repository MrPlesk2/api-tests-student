package myAdvertisements

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/myAdvertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

const pageSize = 50

// GetMyAdvertisements returns one page as JSON and checks the expected status.
func GetMyAdvertisements(t *testing.T, token string, expectedStatusCode, limit int) string {
	t.Helper()
	response := myAdvertisements.HttpGetMyAdvertisements(t, token, limit, "")
	body := utils.ReadResponseBody(t, response)
	require.Equal(t, expectedStatusCode, response.StatusCode, "Unexpected profile status; response: %s", body)
	return string(body)
}

// GetAllMyAdvertisements follows cursors so test data can be found in accounts with more than one page.
func GetAllMyAdvertisements(t *testing.T, token string) []models.Advertisement {
	t.Helper()
	var items []models.Advertisement
	cursor := ""
	seenCursors := make(map[string]bool)
	for {
		page := utils.DecodeResponse[models.AdvertisementsList](t,
			myAdvertisements.HttpGetMyAdvertisements(t, token, pageSize, cursor), http.StatusOK)
		items = append(items, page.Items...)
		if !page.HasNext {
			return items
		}
		require.NotEmpty(t, page.NextCursor, "Profile has another page but returned no cursor")
		require.False(t, seenCursors[page.NextCursor], "Profile pagination repeated cursor %q", page.NextCursor)
		seenCursors[page.NextCursor] = true
		cursor = page.NextCursor
	}
}
