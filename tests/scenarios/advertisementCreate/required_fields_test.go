package advertisementCreate

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

// TestCreateAdvertisementWithRequiredFields checks omitted optional fields and their default values.
func TestCreateAdvertisementWithRequiredFields(t *testing.T) {
	user := fixtures.Login(t)
	t.Log("[Precondition] Prepare an advertisement without price and quantity")
	request := fixtures.AdvertisementRequest(t)
	delete(request.Fields, "price")
	delete(request.Fields, "quantity")
	t.Log("[Step] Create an advertisement using only required fields")
	created := advertisements.CreateAdvertisement(t, user.Token, request)
	assertions.AdvertisementMatchesRequest(t, created, request, user.User.Id)
	require.Zero(t, created.Price, "Omitted price must default to zero")
	require.Zero(t, created.Quantity, "Omitted quantity must default to zero")
	t.Log("[Step] Read the advertisement and check required fields and default values")
	actual := advertisements.GetAdvertisement(t, user.Token, created.ID)
	require.Equal(t, created.ID, actual.ID, "GET returned another advertisement")
	assertions.AdvertisementMatchesRequest(t, actual, request, user.User.Id)
	require.Zero(t, actual.Price, "GET changed the default price")
	require.Zero(t, actual.Quantity, "GET changed the default quantity")
	t.Log("[Result] Required fields were persisted and omitted optional fields defaulted to zero")
}
