package domain

import (
	"context"
)

type SignupRequest struct {
	Name        string `json:"user_name" validate:"required,name"`
	PhoneNumber string `json:"user_phone" validate:"required,phone"`
}

type VerifyOTP struct {
	UserPhone string `json:"user_phone" validate:"required,phone"`
	OTP       string `json:"otp" validate:"required,otp"`
}

type SignupResponse struct {
	Code         int    `json:"code"`
	Message      string `json:"msg"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type VerifyResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type SignupUsecase interface {
	Create(c context.Context, user *User) error
	GetUserByEmail(c context.Context, email string) (User, error)
	GetUserByPhone(c context.Context, phoneNumber string) (User, error)
	CreateAccessToken(user *User, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user *User, secret string, expiry int) (refreshToken string, err error)
	SetOTPExpired(ctx context.Context, params CacheOTP) error
	GetCacheVerify(ctx context.Context, params VerifyOTP) (*CacheOTP, error)
	SetCacheVerify(ctx context.Context, params CacheOTP) error
	DelCache(ctx context.Context, email string) error
	UpdateUserVerify(ctx context.Context, userID int64, isVerify int) (int64, error)
	UpdateOTPUser(ctx context.Context, phoneNumber string, otp int) error
}
