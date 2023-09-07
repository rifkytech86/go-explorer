package domain

import (
	"context"
	"gitlab.com/naonweh-studio/bubbme-backend/dto"
	"time"
)

type AuthRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type AuthResponse struct {
	Message string   `json:"msg"`
	Code    int64    `json:"code"`
	Token   string   `json:"token"`
	Data    []*Users `json:"data"`
}

type Users struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IAuthorizationUsecase interface {
	Login(c context.Context, email string, password string) (*AuthResponse, error)
	LogOut(c context.Context, email string) (*AuthResponse, error)
}

type IUserRepository interface {
	GetUser(ctx context.Context, email string) (user dto.Users, err error)
	UpdateIsLogin(ctx context.Context, email string, isLogin int) (err error)
}
