package controllers

import (
	"JwtWithGo/pkg/dto"
	_ "JwtWithGo/pkg/dto"
	"JwtWithGo/pkg/models"
	"JwtWithGo/pkg/utils"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

var JwtKey = []byte("my_secret_key")
var RefreshTokenKey = []byte("my_refresh_token_key")

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// Register godoc
// @Summary Register a new User
// @Description Register a new User with the system.
// @Tags User
// @Accept  json
// @Produce  json
// @Param   user  body  dto.RegisterRequest  true  "User Data"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {string} string "Invalid input"
// @Router /register [post]
func Register(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		now := time.Now()

		user.Password = hashedPassword
		user.CreatedAt = &now
		if err := db.Create(&user).Error; err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		response := dto.RegisterResponse{
			Username: user.Username,
			Email:    user.Email,
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// Login godoc
// @Summary Login user
// @Description Authenticate a user and return a JWT token.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   login  body  dto.LoginRequest  true  "Login Data"
// @Success 200 {object} dto.TokenResponse
// @Failure 401 {string} string "Unauthorized"
// @Router /login [post]
func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.Where("email = ?", creds.Email).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if !utils.CheckPasswordHash(creds.Password, user.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		tokenString, err := utils.GenerateJWT(user.Email, int(user.ID))
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		refreshTokenString, err := utils.GenerateRefreshToken(user.Email)
		if err != nil {
			http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
			return
		}

		response := dto.TokenResponse{
			Token:        tokenString,
			RefreshToken: refreshTokenString,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Refresh the JWT token using a valid refresh token.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   token  body  dto.RefreshTokenRequest  true  "Refresh Token"
// @Success 200 {object} dto.TokenResponse
// @Failure 401 {string} string "Unauthorized"
// @Router /refresh [post]
func RefreshToken(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		claims, err := utils.ParseRefreshToken(req.RefreshToken)
		if err != nil {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		var user models.User
		if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
			http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
			return
		}

		claims.ExpiresAt = time.Now().Add(time.Hour * -24).Unix()

		newTokenString, err := utils.GenerateJWT(user.Email, int(user.ID))
		if err != nil {
			http.Error(w, "Error generating new token", http.StatusInternalServerError)
			return
		}

		newRefreshToken, err := utils.GenerateRefreshToken(user.Email)
		if err != nil {
			http.Error(w, "Error generating new refresh token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"token":         newTokenString,
			"refresh_token": newRefreshToken,
		})
	}
}

//func Logout(db *gorm.DB) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		claims, err := utils.GetUserIdByToken(r)
//		if err != nil {
//			http.Error(w, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		var user models.User
//		if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
//			http.Error(w, "User not found", http.StatusUnauthorized)
//			return
//		}
//
//		user.RefreshToken = ""
//		user.RefreshTokenExpiry = time.Time{}
//		db.Save(&user)
//
//		w.WriteHeader(http.StatusOK)
//		json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
//	}
//}
