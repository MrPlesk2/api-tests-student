package models

// CreateAdvertisementRequest preserves absent, empty and invalid multipart fields for negative tests.
type CreateAdvertisementRequest struct {
	Fields map[string]string
	Photos []PhotoUpload
}

// PhotoUpload contains the original filename and bytes sent to the API.
type PhotoUpload struct {
	Name    string
	Content []byte
}
