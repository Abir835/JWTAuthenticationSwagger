package dto

// LoginRequest represents the login request payload
// swagger:model LoginRequest
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
