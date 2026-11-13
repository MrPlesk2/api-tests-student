package models

// ErrorResponse describes a rejected API request.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
