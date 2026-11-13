package models

// AddFavoriteRequest identifies the advertisement to add.
type AddFavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

// BulkAddFavoritesRequest identifies the advertisements to merge in one request.
type BulkAddFavoritesRequest struct {
	AdvertisementIDs []string `json:"advertisement_ids"`
}
