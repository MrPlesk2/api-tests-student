package fixtures

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
)

// FavoritesUsers returns an advertiser and a distinct account dedicated to the current scenario.
func FavoritesUsers(t *testing.T) (authModels.LoginOkResponse, authModels.LoginOkResponse) {
	t.Helper()
	owner := Login(t)
	seed := owner.User.Id + "\x00" + t.Name()
	emailHash := sha256.Sum256([]byte(seed))
	passwordHash := sha256.Sum256([]byte(os.Getenv("TEST_PASSWORD") + "\x00" + seed))
	request := authModels.RegisterRequest{
		Email:     fmt.Sprintf("api-favorites-%x@example.com", emailHash[:12]),
		Password:  fmt.Sprintf("Aa1!%x", passwordHash),
		FirstName: "API Favorites",
	}

	t.Log("[Precondition] Register or reuse the dedicated favorites account")
	viewer := auth.RegisterOrLogin(t, request)
	require.NotEmpty(t, viewer.Token, "Favorites account has no token")
	require.NotEmpty(t, viewer.User.Id, "Favorites account has no user ID")
	require.NotEqual(t, owner.User.Id, viewer.User.Id, "Advertiser and favorites user must be different")
	require.Empty(t, favorites.GetAllFavorites(t, viewer.Token),
		"Dedicated favorites account contains entries left by another run")
	t.Cleanup(func() {
		t.Log("[Cleanup] Check that the dedicated account has no remaining favorites")
		require.Empty(t, favorites.GetAllFavorites(t, viewer.Token),
			"Scenario cleanup left entries in the dedicated favorites account")
	})
	return owner, viewer
}
