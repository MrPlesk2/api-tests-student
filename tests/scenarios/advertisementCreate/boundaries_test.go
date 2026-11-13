package advertisementCreate

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/utils"
	"api-tests-template/tests/fixtures"
)

// TestCreateAdvertisementValidBoundaries checks accepted limits and persistence using independent advertisements.
func TestCreateAdvertisementValidBoundaries(t *testing.T) {
	user := fixtures.Login(t)
	cases := []struct{ name, field, value string }{
		{name: "maximum_title", field: "title", value: utils.RandomString(maximumTitleLength)},
		{name: "maximum_description", field: "description", value: strings.Repeat("a", maximumDescriptionLength)},
		{name: "zero_price", field: "price", value: "0"},
		{name: "maximum_price", field: "price", value: fmt.Sprint(maximumPrice)},
		{name: "minimum_quantity", field: "quantity", value: "1"},
		{name: "maximum_quantity", field: "quantity", value: fmt.Sprint(maximumQuantity)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("[Precondition] Prepare valid independent data and change one field to its accepted boundary")
			request := fixtures.AdvertisementRequest(t)
			request.Fields[tc.field] = tc.value
			t.Logf("[Step] Create an advertisement with %s=%q", tc.field, tc.value)
			created := advertisements.CreateAdvertisement(t, user.Token, request)
			assertions.AdvertisementMatchesRequest(t, created, request, user.User.Id)
			t.Log("[Step] Read the advertisement and compare its persisted fields")
			actual := advertisements.GetAdvertisement(t, user.Token, created.ID)
			require.Equal(t, created.ID, actual.ID, "GET returned another advertisement")
			assertions.AdvertisementMatchesRequest(t, actual, request, user.User.Id)
			t.Log("[Result] Accepted boundary was created and persisted without changes")
		})
	}
}
