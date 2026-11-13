package assertions

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

// FindAdvertisementByID requires the exact created advertisement in a result list.
func FindAdvertisementByID(t *testing.T, items []models.Advertisement, id string) models.Advertisement {
	t.Helper()
	for _, item := range items {
		if item.ID == id {
			return item
		}
	}
	require.FailNow(t, "Advertisement is missing from the result",
		"Expected ID %s among %d advertisements", id, len(items))
	return models.Advertisement{}
}

// AdvertisementMatchesRequest checks every submitted field, ownership and photo count.
func AdvertisementMatchesRequest(
	t *testing.T, actual models.Advertisement, request models.CreateAdvertisementRequest, userID string,
) {
	t.Helper()
	require.NotEmpty(t, actual.ID, "Advertisement ID is empty")
	require.Equal(t, request.Fields["title"], actual.Title, "Title differs from the submitted value")
	require.Equal(t, request.Fields["description"], actual.Description, "Description differs from the submitted value")
	if value, present := request.Fields["price"]; present {
		price, err := strconv.ParseFloat(value, 64)
		require.NoError(t, err, "Fixture price is invalid")
		require.Equal(t, price, actual.Price, "Price differs from the submitted value")
	}
	if value, present := request.Fields["quantity"]; present {
		quantity, err := strconv.Atoi(value)
		require.NoError(t, err, "Fixture quantity is invalid")
		require.Equal(t, quantity, actual.Quantity, "Quantity differs from the submitted value")
	}
	require.Equal(t, userID, actual.UserID, "Advertisement belongs to another user")
	require.NotEmpty(t, actual.CreatedAt, "Creation timestamp is empty")
	require.NotEmpty(t, actual.UpdatedAt, "Update timestamp is empty")
	require.Len(t, actual.Photos, len(request.Photos), "Unexpected number of photos")
	t.Logf("[Check] Advertisement %s matches every submitted field", actual.ID)
}
