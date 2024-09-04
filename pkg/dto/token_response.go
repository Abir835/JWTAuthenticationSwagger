package dto

// TokenResponse represents the response payload containing tokens
// swagger:model TokenResponse
type TokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
