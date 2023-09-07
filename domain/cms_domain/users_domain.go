package cms_domain

import (
	"context"
	"time"
)

type UserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserAdminResponse struct {
	Message string   `json:"msg"`
	Code    int64    `json:"code"`
	Total   int      `json:"total"`
	Data    []*Users `json:"data"`
}

type Users struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	Name      string     `json:"name"`
	IsLogin   int        `json:"is_login"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	CreatedBy int        `json:"created_by"`
	UpdatedBy int        `json:"updated_by"`
}

type IUsersAdminUseCase interface {
	GetUserAdmin(ctx context.Context, page int, limit int, search string) (*UserAdminResponse, error)
	CreateUserAdmin(ctx context.Context, req UserRequest) (*UserAdminResponse, error)
	UpdateUserAdmin(ctx context.Context, req UserRequest, id int) (*UserAdminResponse, error)
	DeleteUserAdmin(ctx context.Context, id int) (*UserAdminResponse, error)
}

type IUserAdminRepository interface {
	GetUserAdmin(ctx context.Context, page int, limit int, search string) (user []*Users, err error)
	GetTotalUserAdmin(ctx context.Context) (total int, err error)
	CreateUserAdmin(ctx context.Context, payload *Users) (lastID int64, err error)
	UpdateUserAdmin(ctx context.Context, payload *Users, id int) error
	DeleteUserAdmin(ctx context.Context, id int) error

	GetUserAdminByEmail(ctx context.Context, email string) (user *Users, err error)
}
