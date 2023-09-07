package domain

import (
	"context"
	"database/sql"
)

const (
	CollectionUser = "users"
)

type User struct {
	UserID       sql.NullInt64  `json:"user_id"`
	UserName     sql.NullString `json:"user_name" validate:"required"`
	UserEmail    sql.NullString `json:"user_email" validate:"required"`
	UserPhone    sql.NullString `json:"user_phone" validate:"required"`
	UserPassword sql.NullString `json:"user_password" validate:"required"`
	UserDOB      sql.NullString `json:"user_dob" validate:"required"`
	UserSchool   sql.NullString `json:"user_school" validate:"required"`
	UserIsVerify sql.NullInt64  `json:"user_is_verify" validate:"required"`
	UserOTP      sql.NullInt64  `json:"user_otp" validate:"required"`
	CreatedAt    sql.NullTime   `json:"created_at"`
	UpdatedAt    sql.NullTime   `json:"updated_at"`
}

type CacheOTP struct {
	Phone   string `json:"phone"`
	OTP     string `json:"otp"`
	Attempt int64  `json:"attempt"`
}

type UserRepository interface {
	Create(c context.Context, user *User) (int64, error)
	Fetch(c context.Context) ([]User, error)
	GetByEmail(c context.Context, email string) (User, error)
	GetUserByPhone(c context.Context, phoneNumber string) (User, error)
	UpdateOTPUserByPhone(c context.Context, phoneNumber string, otp int) error
	GetByID(c context.Context, id string) (User, error)

	// cache
	SetExpiredToken(c context.Context, param CacheOTP) error
	GetCacheVerify(c context.Context, param VerifyOTP) (*CacheOTP, error)
	DeleteCacheVerify(c context.Context, phone string) error
	UpdateUserVerify(c context.Context, userID int64, isVerify int) (int64, error)
}
