package cms_domain

import (
	"context"
	"database/sql"
	"encoding/json"
)

type UserPointRequest struct {
	UserID         int64 `json:"user_id"`
	UserPointValue int64 `json:"user_point_value"`
}

type UserPointResponse struct {
	Message string       `json:"msg"`
	Code    int64        `json:"code"`
	Total   int          `json:"total"`
	Data    []*UserPoint `json:"data"`
}

type UserPoint struct {
	UserPointID    sql.NullInt64 `json:"user_point_id"`
	UserID         sql.NullInt64 `json:"user_id"`
	UserPointValue sql.NullInt64 `json:"user_point_value"`
	CreatedAt      sql.NullTime  `json:"created_at"`
	UpdatedAt      sql.NullTime  `json:"updated_at"`
	CreatedBy      sql.NullInt64 `json:"created_by"`
	UpdatedBy      sql.NullInt64 `json:"updated_by"`
}

func (u UserPoint) MarshalJSON() ([]byte, error) {
	type Alias UserPoint
	return json.Marshal(&struct {
		Alias
		UserPoint      interface{} `json:"user_point_id,omitempty"`
		UserID         interface{} `json:"user_id,omitempty"`
		UserPointValue interface{} `json:"user_point_value,omitempty"`
		CreatedAt      interface{} `json:"created_at,omitempty"`
		UpdatedAt      interface{} `json:"updated_at,omitempty"`
		CreatedBy      interface{} `json:"created_by"`
		UpdatedBy      interface{} `json:"updated_by"`
	}{
		Alias:          (Alias)(u),
		UserPoint:      u.UserPointID.Int64,
		UserID:         u.UserID.Int64,
		UserPointValue: u.UserPointValue.Int64,
		CreatedAt:      u.CreatedAt.Time,
		UpdatedAt:      u.UpdatedAt.Time,
		CreatedBy:      u.CreatedBy.Int64,
		UpdatedBy:      u.UpdatedBy.Int64,
	})
}

type IUserPointUseCase interface {
	Get(ctx context.Context, page int, limit int, search string) (*UserPointResponse, error)
	Create(ctx context.Context, req UserPointRequest) (*UserPointResponse, error)
	Update(ctx context.Context, req UserPointRequest, id int) (*UserPointResponse, error)
	Delete(ctx context.Context, id int) (*UserPointResponse, error)
}

type IUserPointRepository interface {
	Get(ctx context.Context, page int, limit int, search string) (user []*UserPoint, err error)
	GetTotal(ctx context.Context, search string) (total int, err error)
	Create(ctx context.Context, payload *UserPoint) (int64, error)
	Update(ctx context.Context, payload *UserPoint, id int) error
	Delete(ctx context.Context, id int) error
}
