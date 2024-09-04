package dto

// RegisterRequest represents the login request payload
// swagger:model RegisterRequest
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
