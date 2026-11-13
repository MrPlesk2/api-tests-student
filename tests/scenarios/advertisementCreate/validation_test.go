package advertisementCreate

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/tests/fixtures"
)

const (
	maximumTitleLength         = 50
	maximumDescriptionLength   = 500
	maximumPrice               = 1000000
	maximumQuantity            = 100
	invalidRequestError        = "invalid_request"
	priceValidationMessage     = "price must be between 0 and 1000000 inclusive"
	quantityValidationMessage  = "quantity must be between 1 and 100 inclusive"
	quantityParsingMessage     = "cannot convert string value to int"
	invalidPhotoContentMessage = "photo 1: file type not allowed: detected content text/plain"
)

// TestCreateAdvertisementInvalidParameters checks one invalid input per independent subtest.
func TestCreateAdvertisementInvalidParameters(t *testing.T) {
	user := fixtures.Login(t)
	cases := []struct {
		name       string
		field      string
		value      string
		omit       bool
		photo      *models.PhotoUpload
		photoCount int
		message    string
	}{
		{name: "missing_title", field: "title", omit: true, photoCount: 1, message: "title is required"},
		{name: "empty_title", field: "title", value: "", photoCount: 1, message: "title is required"},
		{name: "title_above_limit", field: "title", value: strings.Repeat("a", maximumTitleLength+1), photoCount: 1,
			message: fmt.Sprintf("title must not exceed %d characters provided %d", maximumTitleLength, maximumTitleLength+1)},
		{name: "missing_description", field: "description", omit: true, photoCount: 1, message: "description is required"},
		{name: "empty_description", field: "description", value: "", photoCount: 1, message: "description is required"},
		{name: "description_above_limit", field: "description", value: strings.Repeat("a", maximumDescriptionLength+1), photoCount: 1,
			message: "description must be between 1 and 500 characters"},
		{name: "fractional_negative_price", field: "price", value: "-0.01", photoCount: 1, message: priceValidationMessage},
		{name: "fractional_price_above_limit", field: "price", value: "1000000.01", photoCount: 1, message: priceValidationMessage},
		{name: "negative_price", field: "price", value: "-1", photoCount: 1, message: priceValidationMessage},
		{name: "price_above_limit", field: "price", value: fmt.Sprint(maximumPrice + 1), photoCount: 1, message: priceValidationMessage},
		{name: "decimal_comma_price", field: "price", value: "1234,56", photoCount: 1,
			message: "cannot convert string value 1234,56 to float"},
		{name: "grouped_price", field: "price", value: "1 234.56", photoCount: 1,
			message: "cannot convert string value 1 234.56 to float"},
		{name: "non_numeric_price", field: "price", value: "abc", photoCount: 1, message: "cannot convert string value abc to float"},
		{name: "negative_quantity", field: "quantity", value: "-1", photoCount: 1, message: quantityValidationMessage},
		{name: "zero_quantity", field: "quantity", value: "0", photoCount: 1, message: quantityValidationMessage},
		{name: "quantity_above_limit", field: "quantity", value: fmt.Sprint(maximumQuantity + 1), photoCount: 1, message: quantityValidationMessage},
		{name: "fractional_quantity", field: "quantity", value: "1.5", photoCount: 1, message: quantityParsingMessage},
		{name: "scientific_quantity", field: "quantity", value: "1e2", photoCount: 1, message: quantityParsingMessage},
		{name: "overflowing_quantity", field: "quantity", value: "999999999999999999999999999999", photoCount: 1, message: quantityParsingMessage},
		{name: "non_numeric_quantity", field: "quantity", value: "abc", photoCount: 1, message: quantityParsingMessage},
		{name: "text_photo", photoCount: 1,
			photo:   &models.PhotoUpload{Name: "photo.txt", Content: []byte("not an image")},
			message: "photo 1: file type not allowed: extension \".txt\""},
		{name: "svg_photo", photoCount: 1,
			photo:   &models.PhotoUpload{Name: "photo.svg", Content: []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="1" height="1"/>`)},
			message: "photo 1: file type not allowed: extension \".svg\""},
		{name: "empty_photo", photoCount: 1,
			photo: &models.PhotoUpload{Name: "photo.png"}, message: invalidPhotoContentMessage},
		{name: "corrupt_photo", photoCount: 1,
			photo:   &models.PhotoUpload{Name: "photo.png", Content: []byte("not a valid PNG")},
			message: invalidPhotoContentMessage},
		{name: "missing_photos", photoCount: 0, message: "at least 1 photo is required"},
		{name: "too_many_photos", photoCount: maximumPhotos + 1, message: "maximum 3 photos allowed"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("[Precondition] Prepare valid independent data and change only the tested parameter")
			request := fixtures.AdvertisementRequest(t)
			if tc.omit {
				delete(request.Fields, tc.field)
			} else if tc.field != "" {
				request.Fields[tc.field] = tc.value
			}
			photo := request.Photos[0]
			if tc.photo != nil {
				photo = *tc.photo
			}
			request.Photos = nil
			for i := 0; i < tc.photoCount; i++ {
				request.Photos = append(request.Photos, photo)
			}

			t.Logf("[Step] Submit invalid request: %s", tc.name)
			response := advertisements.CreateAdvertisementExpectError(t, user.Token, request, http.StatusBadRequest)
			require.Equal(t, invalidRequestError, response.Error, "Wrong validation error type")
			require.Equal(t, tc.message, response.Message, "Wrong validation error message")
			t.Logf("[Result] HTTP 400, error=%s, message=%q", response.Error, response.Message)
		})
	}
}
