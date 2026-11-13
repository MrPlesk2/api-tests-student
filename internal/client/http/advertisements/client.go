package advertisements

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpGetAdvertisements searches advertisements by title.
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		End().Response
}

// HttpPostAdvertisement sends multipart fields and photo bytes, including deliberately invalid values.
func HttpPostAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) *http.Response {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for name, value := range request.Fields {
		require.NoError(t, writer.WriteField(name, value), "Cannot encode multipart field %q", name)
	}
	for _, photo := range request.Photos {
		part, err := writer.CreateFormFile("photos", photo.Name)
		require.NoError(t, err, "Cannot encode photo %q", photo.Name)
		_, err = part.Write(photo.Content)
		require.NoError(t, err, "Cannot write photo %q", photo.Name)
	}
	require.NoError(t, writer.Close(), "Cannot finish multipart request")

	return apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		ContentType(writer.FormDataContentType()).
		Body(body.String()).
		Expect(t).
		End().Response
}

// HttpGetAdvertisement reads a single advertisement using the documented id query parameter.
func HttpGetAdvertisement(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement removes an advertisement owned by the authenticated user.
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos returns the photos attached to an advertisement.
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).Create().
		Get(fmt.Sprintf("%s/%s/photos", path.AdvertisementsPath, id)).
		Expect(t).
		End().Response
}

// HttpGetPhoto downloads a photo from the storage URL returned by the API.
func HttpGetPhoto(t *testing.T, photoURL string) *http.Response {
	t.Helper()
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, photoURL, nil)
	require.NoError(t, err, "Invalid photo URL %q", photoURL)
	client := &http.Client{Timeout: apiRunner.HTTPTimeout}
	response, err := client.Do(request)
	require.NoError(t, err, "Cannot download photo %q", photoURL)
	return response
}
