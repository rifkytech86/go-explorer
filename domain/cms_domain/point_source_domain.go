package cms_domain

import (
	"context"
	"database/sql"
	"encoding/json"
)

type PointSourceRequest struct {
	PointSourceName string `json:"point_source_name"`
}

type PointSourceResponse struct {
	Message string         `json:"msg"`
	Code    int64          `json:"code"`
	Total   int            `json:"total"`
	Data    []*PointSource `json:"data"`
}

type PointSource struct {
	PointSourceID   sql.NullInt64  `json:"point_source_id"`
	PointSourceName sql.NullString `json:"point_source_name"`
	CreatedAt       sql.NullTime   `json:"created_at"`
	UpdatedAt       sql.NullTime   `json:"updated_at"`
	CreatedBy       sql.NullInt64  `json:"created_by"`
	UpdatedBy       sql.NullInt64  `json:"updated_by"`
}

func (u PointSource) MarshalJSON() ([]byte, error) {
	type Alias PointSource
	return json.Marshal(&struct {
		Alias
		PointSourceID   interface{} `json:"point_source_id,omitempty"`
		PointSourceName interface{} `json:"point_source_name,omitempty"`
		CreatedAt       interface{} `json:"created_at,omitempty"`
		UpdatedAt       interface{} `json:"updated_at,omitempty"`
		CreatedBy       interface{} `json:"created_by"`
		UpdatedBy       interface{} `json:"updated_by"`
	}{
		Alias:           (Alias)(u),
		PointSourceID:   u.PointSourceID.Int64,
		PointSourceName: u.PointSourceName.String,
		CreatedAt:       u.CreatedAt.Time,
		UpdatedAt:       u.UpdatedAt.Time,
		CreatedBy:       u.CreatedBy.Int64,
		UpdatedBy:       u.UpdatedBy.Int64,
	})
}

type IPointSourceUseCase interface {
	Get(ctx context.Context, page int, limit int, search string) (*PointSourceResponse, error)
	CreatePointSource(ctx context.Context, req PointSourceRequest) (*PointSourceResponse, error)
	UpdatePointSource(ctx context.Context, req PointSourceRequest, id int) (*PointSourceResponse, error)
	DeletePointSource(ctx context.Context, id int) (*PointSourceResponse, error)
}

type IPointSourceRepository interface {
	Get(ctx context.Context, page int, limit int, search string) (user []*PointSource, err error)
	GetTotalPointSource(ctx context.Context, search string) (total int, err error)
	CreatePointSource(ctx context.Context, payload *PointSource) (int64, error)
	UpdatePointSource(ctx context.Context, payload *PointSource, id int) error
	DeletePointSource(ctx context.Context, id int) error
}
