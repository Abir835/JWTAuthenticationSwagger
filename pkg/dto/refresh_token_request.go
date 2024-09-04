package dto

// RefreshTokenRequest represents the request payload for refreshing a token
// swagger:model RefreshTokenRequest
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
