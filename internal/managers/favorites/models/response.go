package models

import advertisementModels "api-tests-template/internal/managers/advertisements/models"

// Favorite contains the saved advertisement and the favorite entry's own metadata.
type Favorite struct {
	FavoriteID    string                            `json:"favorite_id"`
	AddedAt       string                            `json:"added_at"`
	Advertisement advertisementModels.Advertisement `json:"advertisement"`
}

// FavoritesList describes page-based favorites pagination.
type FavoritesList struct {
	Items      []Favorite `json:"items"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	Total      int        `json:"total"`
	TotalPages int        `json:"total_pages"`
}

// BulkAddFavoritesResponse distinguishes newly added, skipped and nonexistent advertisements.
type BulkAddFavoritesResponse struct {
	Added        []string `json:"added"`
	Skipped      []string `json:"skipped"`
	NotFound     []string `json:"not_found"`
	LimitReached bool     `json:"limit_reached"`
}
