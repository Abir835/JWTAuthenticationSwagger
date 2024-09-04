package dto

// RegisterResponse represents the login request payload
// swagger:model RegisterResponse
type RegisterResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
