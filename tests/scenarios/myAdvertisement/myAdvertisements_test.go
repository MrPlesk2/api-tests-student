package myAdvertisement

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/myAdvertisements"
	"api-tests-template/tests/fixtures"
)

const (
	profilePageSize = 50
	invalidToken    = "incorrect_token"
)

// TestGetMyAdvertisementsPositive checks only the advertisement prepared by this test.
func TestGetMyAdvertisementsPositive(t *testing.T) {
	user := fixtures.Login(t)
	t.Log("[Precondition] Create an advertisement for the configured account")
	request := fixtures.AdvertisementRequest(t)
	created := advertisements.CreateAdvertisement(t, user.Token, request)

	t.Log("[Step] Read all pages of the account's own advertisements")
	items := myAdvertisements.GetAllMyAdvertisements(t, user.Token)
	found := assertions.FindAdvertisementByID(t, items, created.ID)
	assertions.AdvertisementMatchesRequest(t, found, request, user.User.Id)
	for _, item := range items {
		require.Equal(t, user.User.Id, item.UserID, "Profile returned another user's advertisement %s", item.ID)
	}
	t.Log("[Result] The created advertisement is in its owner's profile")
}

// TestGetMyAdvertisementsIncorrectToken checks the exact authorization error.
func TestGetMyAdvertisementsIncorrectToken(t *testing.T) {
	fixtures.LoadEnvironment(t)
	t.Log("[Step] Request the user's advertisements with an invalid token")
	body := myAdvertisements.GetMyAdvertisements(t, invalidToken, http.StatusUnauthorized, profilePageSize)
	var response models.ErrorResponse
	require.NoError(t, json.Unmarshal([]byte(body), &response), "Invalid error JSON: %s", body)
	require.Equal(t, "unauthorized", response.Error, "Wrong authorization error type")
	require.Equal(t, "Invalid or expired token", response.Message, "Wrong authorization error message")
	t.Log("[Result] Invalid token was rejected with HTTP 401 and the expected message")
}
