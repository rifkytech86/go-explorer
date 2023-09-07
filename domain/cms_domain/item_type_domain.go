package cms_domain

import (
	"context"
	"database/sql"
	"encoding/json"
)

type ItemTypeRequest struct {
	ItemTypeName string `json:"item_type_name"`
}

type ItemTypeResponse struct {
	Message string      `json:"msg"`
	Code    int64       `json:"code"`
	Total   int         `json:"total"`
	Data    []*ItemType `json:"data"`
}

type ItemType struct {
	ItemTypeID   sql.NullInt64  `json:"item_type_id"`
	ItemTypeName sql.NullString `json:"item_type_name"`
	CreatedAt    sql.NullTime   `json:"created_at"`
	UpdatedAt    sql.NullTime   `json:"updated_at"`
	CreatedBy    sql.NullInt64  `json:"created_by"`
	UpdatedBy    sql.NullInt64  `json:"updated_by"`
}

func (u ItemType) MarshalJSON() ([]byte, error) {
	type Alias ItemType
	return json.Marshal(&struct {
		Alias
		ItemTypeID   interface{} `json:"item_type_id,omitempty"`
		ItemTypeName interface{} `json:"item_type_name,omitempty"`
		CreatedAt    interface{} `json:"created_at,omitempty"`
		UpdatedAt    interface{} `json:"updated_at,omitempty"`
		CreatedBy    interface{} `json:"created_by"`
		UpdatedBy    interface{} `json:"updated_by"`
	}{
		Alias:        (Alias)(u),
		ItemTypeID:   u.ItemTypeID.Int64,
		ItemTypeName: u.ItemTypeName.String,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
		CreatedBy:    u.CreatedBy.Int64,
		UpdatedBy:    u.UpdatedBy.Int64,
	})
}

type IItemTypeUseCase interface {
	Get(ctx context.Context, page int, limit int, search string) (*ItemTypeResponse, error)
	CreateItemType(ctx context.Context, req ItemTypeRequest) (*ItemTypeResponse, error)
	UpdateItemType(ctx context.Context, req ItemTypeRequest, id int) (*ItemTypeResponse, error)
	DeleteItemType(ctx context.Context, id int) (*ItemTypeResponse, error)
}

type IItemTypeRepository interface {
	Get(ctx context.Context, page int, limit int, search string) (user []*ItemType, err error)
	GetTotalItemType(ctx context.Context, search string) (total int, err error)
	CreateItemType(ctx context.Context, payload *ItemType) (int64, error)
	UpdateItemType(ctx context.Context, payload *ItemType, id int) error
	DeleteItemType(ctx context.Context, id int) error
}
