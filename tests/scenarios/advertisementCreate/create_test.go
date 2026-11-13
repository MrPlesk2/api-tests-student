package advertisementCreate

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/tests/fixtures"
)

const maximumPhotos = 3

// TestCreateAdvertisementWithAllParameters checks creation, retrieval, full-title search and stored photos.
func TestCreateAdvertisementWithAllParameters(t *testing.T) {
	user := fixtures.Login(t)
	request := fixtures.AdvertisementRequest(t)
	request.Photos = []models.PhotoUpload{
		fixtures.PNGPhoto(t, "red.png", color.RGBA{R: 255, A: 255}),
		fixtures.PNGPhoto(t, "green.png", color.RGBA{G: 255, A: 255}),
		fixtures.PNGPhoto(t, "blue.png", color.RGBA{B: 255, A: 255}),
	}

	t.Log("[Step] Create an advertisement with all fields and the maximum number of photos")
	created := advertisements.CreateAdvertisement(t, user.Token, request)
	assertions.AdvertisementMatchesRequest(t, created, request, user.User.Id)

	t.Log("[Step] Read the advertisement by ID and compare the submitted fields")
	actual := advertisements.GetAdvertisement(t, user.Token, created.ID)
	require.Equal(t, created.ID, actual.ID, "GET returned another advertisement")
	assertions.AdvertisementMatchesRequest(t, actual, request, user.User.Id)
	require.Equal(t, created.Photos, actual.Photos, "GET changed the uploaded photo metadata")

	t.Logf("[Step] Search by the complete title %q", request.Fields["title"])
	list := advertisements.GetAdvertisementsList(t, user.Token, request.Fields["title"])
	found := assertions.FindAdvertisementByID(t, list.Items, created.ID)
	assertions.AdvertisementMatchesRequest(t, found, request, user.User.Id)

	t.Log("[Step] Retrieve photo metadata and download each photo from storage")
	photos := advertisements.GetAdvertisementPhotos(t, user.Token, created.ID)
	require.Len(t, photos, len(request.Photos), "Photo endpoint lost uploaded photos")
	seenIDs := make(map[string]bool, len(photos))
	for i, stored := range photos {
		require.NotEmpty(t, stored.ID, "Photo %d has no ID", i)
		require.False(t, seenIDs[stored.ID], "Duplicate photo ID %s", stored.ID)
		seenIDs[stored.ID] = true
		require.NotEmpty(t, stored.URL, "Photo %s has no storage URL", stored.ID)
		require.NotEmpty(t, stored.CreatedAt, "Photo %s has no creation timestamp", stored.ID)
		require.Equal(t, i, stored.SortOrder, "Unexpected photo order")
		require.Equal(t, created.Photos[i], stored, "Photo endpoint returned different metadata")
		require.Equal(t, request.Photos[i].Content, advertisements.DownloadPhoto(t, stored.URL),
			"Stored photo %s differs from the uploaded PNG", stored.ID)
		t.Logf("[Check] Photo %s is available and matches the uploaded file", stored.ID)
	}
	t.Log("[Result] Created advertisement, full-title search and all uploaded photos are correct")
}
