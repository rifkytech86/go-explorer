package cms_domain

import (
	"context"
	"database/sql"
	"encoding/json"
)

type ItemRequest struct {
	ItemName   string `json:"item_name"`
	ItemDesc   string `json:"item_desc"`
	ItemValue  int64  `json:"item_value"`
	ItemTypeID int64  `json:"item_type_id"`
}

type ItemResponse struct {
	Message string  `json:"msg"`
	Code    int64   `json:"code"`
	Total   int     `json:"total"`
	Data    []*Item `json:"data"`
}

type Item struct {
	ItemID     sql.NullInt64  `json:"item_id"`
	ItemName   sql.NullString `json:"item_name"`
	ItemDesc   sql.NullString `json:"item_desc"`
	ItemValue  sql.NullInt64  `json:"item_value"`
	ItemTypeID sql.NullInt64  `json:"item_type_id"`
	CreatedAt  sql.NullTime   `json:"created_at"`
	UpdatedAt  sql.NullTime   `json:"updated_at"`
	CreatedBy  sql.NullInt64  `json:"created_by"`
	UpdatedBy  sql.NullInt64  `json:"updated_by"`
}

func (u Item) MarshalJSON() ([]byte, error) {
	type Alias Item
	return json.Marshal(&struct {
		Alias
		ItemID     interface{} `json:"item_id,omitempty"`
		ItemName   interface{} `json:"item_name,omitempty"`
		ItemDesc   interface{} `json:"item_desc,omitempty"`
		ItemValue  interface{} `json:"item_value,omitempty"`
		ItemTypeID interface{} `json:"item_type_id,omitempty"`
		CreatedAt  interface{} `json:"created_at,omitempty"`
		UpdatedAt  interface{} `json:"updated_at,omitempty"`
		CreatedBy  interface{} `json:"created_by"`
		UpdatedBy  interface{} `json:"updated_by"`
	}{
		Alias:      (Alias)(u),
		ItemID:     u.ItemID.Int64,
		ItemName:   u.ItemName.String,
		ItemDesc:   u.ItemDesc.String,
		ItemValue:  u.ItemValue.Int64,
		ItemTypeID: u.ItemTypeID.Int64,
		CreatedAt:  u.CreatedAt.Time,
		UpdatedAt:  u.UpdatedAt.Time,
		CreatedBy:  u.CreatedBy.Int64,
		UpdatedBy:  u.UpdatedBy.Int64,
	})
}

type IItemUseCase interface {
	Get(ctx context.Context, page int, limit int, search string) (*ItemResponse, error)
	CreateItemType(ctx context.Context, req ItemRequest) (*ItemResponse, error)
	UpdateItemType(ctx context.Context, req ItemRequest, id int) (*ItemResponse, error)
	DeleteItemType(ctx context.Context, id int) (*ItemResponse, error)
}

type IItemRepository interface {
	Get(ctx context.Context, page int, limit int, search string) (user []*Item, err error)
	GetTotalItem(ctx context.Context, search string) (total int, err error)
	CreateItem(ctx context.Context, payload *Item) (int64, error)
	UpdateItem(ctx context.Context, payload *Item, id int) error
	DeleteItem(ctx context.Context, id int) error
}
