package advertisements

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

// GetAdvertisementsList returns the advertisements matching the search string.
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	t.Helper()
	return utils.DecodeResponse[models.AdvertisementsList](t,
		advertisements.HttpGetAdvertisements(t, token, text), http.StatusOK)
}

// CreateAdvertisement creates an advertisement and registers its deletion even if later checks fail.
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	t.Helper()
	body := createAdvertisement(t, token, request, http.StatusCreated)
	var result models.Advertisement
	require.NoError(t, json.Unmarshal(body, &result), "Invalid advertisement response: %s", body)
	require.NotEmpty(t, result.ID, "Created advertisement has no ID")
	return result
}

// CreateAdvertisementExpectError checks the rejected request status and returns its structured error.
func CreateAdvertisementExpectError(
	t *testing.T, token string, request models.CreateAdvertisementRequest, expectedStatus int,
) models.ErrorResponse {
	t.Helper()
	body := createAdvertisement(t, token, request, expectedStatus)
	var result models.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &result), "Invalid error response: %s", body)
	return result
}

// createAdvertisement registers cleanup before asserting the status, including unexpected creations.
func createAdvertisement(
	t *testing.T, token string, request models.CreateAdvertisementRequest, expectedStatus int,
) []byte {
	t.Helper()
	response := advertisements.HttpPostAdvertisement(t, token, request)
	body := utils.ReadResponseBody(t, response)
	if response.StatusCode == http.StatusCreated {
		var identity struct {
			ID string `json:"id"`
		}
		require.NoError(t, json.Unmarshal(body, &identity), "Invalid creation response: %s", body)
		if identity.ID != "" {
			t.Cleanup(func() {
				t.Logf("[Cleanup] Delete advertisement %s", identity.ID)
				DeleteAdvertisementIfExists(t, token, identity.ID)
			})
		}
	}
	require.Equal(t, expectedStatus, response.StatusCode,
		"Unexpected creation status; response: %s", body)
	return body
}

// GetAdvertisement returns an advertisement by ID and requires HTTP 200.
func GetAdvertisement(t *testing.T, token, id string) models.Advertisement {
	t.Helper()
	return utils.DecodeResponse[models.Advertisement](t,
		advertisements.HttpGetAdvertisement(t, token, id), http.StatusOK)
}

// DeleteAdvertisement removes a test-owned advertisement and requires HTTP 204.
func DeleteAdvertisement(t *testing.T, token, id string) {
	t.Helper()
	response := advertisements.HttpDeleteAdvertisement(t, token, id)
	body := utils.ReadResponseBody(t, response)
	require.Equal(t, http.StatusNoContent, response.StatusCode,
		"Cannot delete advertisement %s; response: %s", id, body)
}

// GetAdvertisementPhotos returns the photo metadata for the specified advertisement.
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	t.Helper()
	return utils.DecodeResponse[[]models.Photo](t,
		advertisements.HttpGetAdvertisementPhotos(t, token, id), http.StatusOK)
}

// DownloadPhoto checks storage availability and returns the actual uploaded bytes.
func DownloadPhoto(t *testing.T, photoURL string) []byte {
	t.Helper()
	response := advertisements.HttpGetPhoto(t, photoURL)
	body := utils.ReadResponseBody(t, response)
	require.Equal(t, http.StatusOK, response.StatusCode,
		"Photo is unavailable at %s; response: %s", photoURL, body)
	require.Contains(t, response.Header.Get("Content-Type"), "image/",
		"Storage returned non-image content at %s", photoURL)
	return body
}

// DeleteAdvertisementIfExists confirms absence when DELETE reports forbidden for an already deleted ID.
func DeleteAdvertisementIfExists(t *testing.T, token, id string) {
	t.Helper()
	response := advertisements.HttpDeleteAdvertisement(t, token, id)
	body := utils.ReadResponseBody(t, response)
	if response.StatusCode == http.StatusForbidden {
		check := advertisements.HttpGetAdvertisement(t, token, id)
		checkBody := utils.ReadResponseBody(t, check)
		require.Equal(t, http.StatusNotFound, check.StatusCode,
			"Cleanup DELETE failed and advertisement %s is not confirmed absent; DELETE: %s; GET: %s", id, body, checkBody)
		return
	}
	require.Contains(t, []int{http.StatusNoContent, http.StatusNotFound}, response.StatusCode,
		"Cannot clean up advertisement %s; response: %s", id, body)
}
