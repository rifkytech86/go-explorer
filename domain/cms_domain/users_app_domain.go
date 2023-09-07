package cms_domain

import (
	"context"
	"database/sql"
	"encoding/json"
)

type AuthAppRequest struct {
	UserName  string `json:"user_name"`
	UserPhone string `json:"user_phone"`
}

type UserAppRequest struct {
	Username     string `json:"user_name"`
	UserPhone    string `json:"user_phone"`
	UserIsVerify int64  `json:"user_is_verify"`
}

type UserAppResponse struct {
	Message string      `json:"msg"`
	Code    int64       `json:"code"`
	Total   int         `json:"total"`
	Data    []*UsersApp `json:"data"`
}

type UsersApp struct {
	UserID        int64          `json:"user_id"`
	UserName      sql.NullString `json:"user_name"`
	UserPhone     sql.NullString `json:"user_phone"`
	UserISVerify  sql.NullInt64  `json:"user_is_verify"`
	UserOTP       sql.NullInt64  `json:"user_otp"`
	UserCreatedAt sql.NullTime   `json:"created_at"`
	UserUpdatedAt sql.NullTime   `json:"updated_at"`
	UserCreatedBy sql.NullInt64  `json:"created_by"`
	UserUpdatedBy sql.NullInt64  `json:"updated_by"`
}

func (u UsersApp) MarshalJSON() ([]byte, error) {
	type Alias UsersApp
	return json.Marshal(&struct {
		Alias
		UserName      interface{} `json:"user_name,omitempty"`
		UserPhone     interface{} `json:"user_phone,omitempty"`
		UserISVerify  interface{} `json:"user_is_verify,omitempty"`
		UserOTP       interface{} `json:"user_otp,omitempty"`
		UserCreatedAt interface{} `json:"created_at,omitempty"`
		UserUpdatedAt interface{} `json:"updated_at,omitempty"`
		UserCreatedBy interface{} `json:"created_by,omitempty"`
		UserUpdatedBy interface{} `json:"updated_by,omitempty"`
	}{
		Alias:         (Alias)(u),
		UserName:      u.UserName.String,
		UserPhone:     u.UserPhone.String,
		UserISVerify:  u.UserISVerify.Int64,
		UserOTP:       u.UserOTP.Int64,
		UserCreatedAt: u.UserCreatedAt.Time,
		UserCreatedBy: u.UserCreatedBy.Int64,
		UserUpdatedBy: u.UserUpdatedBy.Int64,
	})
}

type IUsersAppUseCase interface {
	GetUserApp(ctx context.Context, page int, limit int, search string) (*UserAppResponse, error)
	CreateUserApp(ctx context.Context, req UserAppRequest) (*UserAppResponse, error)
	UpdateUserApp(ctx context.Context, req UserAppRequest, id int) (*UserAppResponse, error)
	DeleteUserApp(ctx context.Context, id int) (*UserAppResponse, error)
}

type IUserAppRepository interface {
	GetUserApp(ctx context.Context, page int, limit int, search string) (user []*UsersApp, err error)
	GetTotalUserApp(ctx context.Context) (total int, err error)
	CreateUserApp(ctx context.Context, payload *UsersApp) (lastID int64, err error)
	UpdateUserApp(ctx context.Context, payload *UsersApp, id int) error
	DeleteUserApp(ctx context.Context, id int) error

	GetUserAppByPhone(ctx context.Context, phone string) (user *UsersApp, err error)
}
