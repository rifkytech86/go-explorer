package domain

import (
	"context"
)

type LoginRequest struct {
	Email       string `json:"user_email"`
	Password    string `json:"user_password"`
	PhoneNumber string `json:"user_phone" validate:"required,user_phone"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginUsecase interface {
	GetUserByEmail(c context.Context, email string) (User, error)
	GetUserByPhone(c context.Context, phoneNumber string) (User, error)
	UpdateOTPUser(c context.Context, phoneNumber string, otp string) error
	CreateAccessToken(user *User, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user *User, secret string, expiry int) (refreshToken string, err error)
}
