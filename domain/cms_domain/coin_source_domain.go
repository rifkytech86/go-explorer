package cms_domain

import (
	"context"
	"database/sql"
	"encoding/json"
)

type CoinSourceRequest struct {
	CoinSourceName string `json:"coin_source_name"`
}

type CoinSourceResponse struct {
	Message string        `json:"msg"`
	Code    int64         `json:"code"`
	Total   int           `json:"total"`
	Data    []*CoinSource `json:"data"`
}

type CoinSource struct {
	CoinSourceID   sql.NullInt64  `json:"coin_source_id"`
	CoinSourceName sql.NullString `json:"coin_source_name"`
	CreatedAt      sql.NullTime   `json:"created_at"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`
	CreatedBy      sql.NullInt64  `json:"created_by"`
	UpdatedBy      sql.NullInt64  `json:"updated_by"`
}

func (u CoinSource) MarshalJSON() ([]byte, error) {
	type Alias CoinSource
	return json.Marshal(&struct {
		Alias
		CoinSourceID   interface{} `json:"coin_source_id,omitempty"`
		CoinSourceName interface{} `json:"coin_source_name,omitempty"`
		CreatedAt      interface{} `json:"created_at,omitempty"`
		UpdatedAt      interface{} `json:"updated_at,omitempty"`
		CreatedBy      interface{} `json:"created_by"`
		UpdatedBy      interface{} `json:"updated_by"`
	}{
		Alias:          (Alias)(u),
		CoinSourceID:   u.CoinSourceID.Int64,
		CoinSourceName: u.CoinSourceName.String,
		CreatedAt:      u.CreatedAt.Time,
		UpdatedAt:      u.UpdatedAt.Time,
		CreatedBy:      u.CreatedBy.Int64,
		UpdatedBy:      u.UpdatedBy.Int64,
	})
}

type ICoinSourceUseCase interface {
	Get(ctx context.Context, page int, limit int, search string) (*CoinSourceResponse, error)
	CreateCoinSource(ctx context.Context, req CoinSourceRequest) (*CoinSourceResponse, error)
	UpdateCoinSource(ctx context.Context, req CoinSourceRequest, id int) (*CoinSourceResponse, error)
	DeleteCoinSource(ctx context.Context, id int) (*CoinSourceResponse, error)
}

type ICoinSourceRepository interface {
	Get(ctx context.Context, page int, limit int, search string) (user []*CoinSource, err error)
	GetTotalCoinSource(ctx context.Context, search string) (total int, err error)
	CreateCoinSource(ctx context.Context, payload *CoinSource) (int64, error)
	UpdateCoinSource(ctx context.Context, payload *CoinSource, id int) error
	DeleteCoinSource(ctx context.Context, id int) error
}
