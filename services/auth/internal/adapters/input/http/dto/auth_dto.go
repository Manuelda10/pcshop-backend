package dto

import "auth-service/internal/core/domain"

// -------------------------
// Requests
// -------------------------

type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"fullName" validate:"required,min=2"`
}

type VerifyEmailRequest struct {
	UserID string `json:"userId" validate:"required,uuid"`
	Code   string `json:"code"   validate:"required,len=6"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// -------------------------
// Responses
// -------------------------

type UserResponse struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	FullName   string `json:"fullName"`
	Provider   string `json:"provider"`
	IsVerified bool   `json:"isVerified"`
	CreatedAt  string `json:"createdAt"`
}

type AuthResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	TokenType    string       `json:"tokenType"`
	ExpiresIn    int          `json:"expiresIn"`
	User         UserResponse `json:"user"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

// -------------------------
// Mappers
// -------------------------

func ToUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:         u.ID.String(),
		Email:      u.Email,
		FullName:   u.FullName,
		Provider:   string(u.Provider),
		IsVerified: u.IsVerified,
		CreatedAt:  u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
