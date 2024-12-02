package models

// AuthResponse represents the response for authentication endpoints
type AuthResponse struct {
	UserID string `json:"user_id" example:"uId123456"`
	Token  string `json:"token" example:"token123456"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid credentials"`
}