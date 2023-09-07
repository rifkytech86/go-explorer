package cms_domain

import (
	"context"
	"database/sql"
	"encoding/json"
)

type UserCoinRequest struct {
	UserID        int64 `json:"user_id"`
	UserCoinValue int64 `json:"user_coin_value"`
}

type UserCoinResponse struct {
	Message string      `json:"msg"`
	Code    int64       `json:"code"`
	Total   int         `json:"total"`
	Data    []*UserCoin `json:"data"`
}

type UserCoin struct {
	UserCoinID    sql.NullInt64 `json:"user_coin_id"`
	UserID        sql.NullInt64 `json:"user_id"`
	UserCoinValue sql.NullInt64 `json:"user_coin_value"`
	CreatedAt     sql.NullTime  `json:"created_at"`
	UpdatedAt     sql.NullTime  `json:"updated_at"`
	CreatedBy     sql.NullInt64 `json:"created_by"`
	UpdatedBy     sql.NullInt64 `json:"updated_by"`
}

func (u UserCoin) MarshalJSON() ([]byte, error) {
	type Alias UserCoin
	return json.Marshal(&struct {
		Alias
		UserCoinID    interface{} `json:"user_coin_id,omitempty"`
		UserID        interface{} `json:"user_id,omitempty"`
		UserCoinValue interface{} `json:"user_coin_value,omitempty"`
		CreatedAt     interface{} `json:"created_at,omitempty"`
		UpdatedAt     interface{} `json:"updated_at,omitempty"`
		CreatedBy     interface{} `json:"created_by"`
		UpdatedBy     interface{} `json:"updated_by"`
	}{
		Alias:         (Alias)(u),
		UserCoinID:    u.UserCoinID.Int64,
		UserID:        u.UserID.Int64,
		UserCoinValue: u.UserCoinValue.Int64,
		CreatedAt:     u.CreatedAt.Time,
		UpdatedAt:     u.UpdatedAt.Time,
		CreatedBy:     u.CreatedBy.Int64,
		UpdatedBy:     u.UpdatedBy.Int64,
	})
}

type IUserCoinUseCase interface {
	Get(ctx context.Context, page int, limit int, search string) (*UserCoinResponse, error)
	Create(ctx context.Context, req UserCoinRequest) (*UserCoinResponse, error)
	Update(ctx context.Context, req UserCoinRequest, id int) (*UserCoinResponse, error)
	Delete(ctx context.Context, id int) (*UserCoinResponse, error)
}

type IUserCoinRepository interface {
	Get(ctx context.Context, page int, limit int, search string) (user []*UserCoin, err error)
	GetTotal(ctx context.Context, search string) (total int, err error)
	Create(ctx context.Context, payload *UserCoin) (int64, error)
	Update(ctx context.Context, payload *UserCoin, id int) error
	Delete(ctx context.Context, id int) error
}
