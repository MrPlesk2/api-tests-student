package advertisementsSearch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"
)

const searchFragmentLength = 12

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
	request   advertisementModels.CreateAdvertisementRequest
	created   advertisementModels.Advertisement
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

// SetupSuite authenticates the configured test account.
func (s *TestSuite) SetupSuite() {
	s.loginData = fixtures.Login(s.T())
}

// SetupTest creates data owned by this test and registers its cleanup.
func (s *TestSuite) SetupTest() {
	s.T().Log("[Precondition] Create an advertisement with a unique title for this test")
	s.request = fixtures.AdvertisementRequest(s.T())
	s.created = advertisements.CreateAdvertisement(s.T(), s.loginData.Token, s.request)
}

// TearDownSuite logs completion of the search suite.
func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// TestSearchByTitle checks the full title, each substring position and case-insensitive search.
func (s *TestSuite) TestSearchByTitle() {
	title := s.created.Title
	middle := title[len(title)/2-searchFragmentLength/2 : len(title)/2+searchFragmentLength/2]
	cases := []struct {
		name   string
		search string
	}{
		{"full_title", title},
		{"beginning", title[:searchFragmentLength]},
		{"middle", middle},
		{"end", title[len(title)-searchFragmentLength:]},
		{"lowercase", strings.ToLower(middle)},
		{"uppercase", strings.ToUpper(middle)},
		{"mixed_case", utils.MixCase(middle)},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.T().Logf("[Step] Search for %q", tc.search)
			list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			found := assertions.FindAdvertisementByID(s.T(), list.Items, s.created.ID)
			assertions.AdvertisementMatchesRequest(s.T(), found, s.request, s.loginData.User.Id)
			s.T().Log("[Result] Search returned the advertisement created by this test")
		})
	}
}
