package cms_domain

import (
	"context"
	"database/sql"
	"encoding/json"
)

type PetStatusRequest struct {
	PetStatusName string `json:"pet_status_name"`
}

type PetStatusResponse struct {
	Message string       `json:"msg"`
	Code    int64        `json:"code"`
	Total   int          `json:"total"`
	Data    []*PetStatus `json:"data"`
}

type PetStatus struct {
	PetStatusID   sql.NullInt64  `json:"pet_status_id"`
	PetStatusName sql.NullString `json:"pet_status_name"`
	CreatedAt     sql.NullTime   `json:"created_at"`
	UpdatedAt     sql.NullTime   `json:"updated_at"`
	CreatedBy     sql.NullInt64  `json:"created_by"`
	UpdatedBy     sql.NullInt64  `json:"updated_by"`
}

func (u PetStatus) MarshalJSON() ([]byte, error) {
	type Alias PetStatus
	return json.Marshal(&struct {
		Alias
		PetStatusID   interface{} `json:"pet_status_id,omitempty"`
		PetStatusName interface{} `json:"pet_status_name,omitempty"`
		CreatedAt     interface{} `json:"created_at,omitempty"`
		UpdatedAt     interface{} `json:"updated_at,omitempty"`
		CreatedBy     interface{} `json:"created_by"`
		UpdatedBy     interface{} `json:"updated_by"`
	}{
		Alias:         (Alias)(u),
		PetStatusID:   u.PetStatusID.Int64,
		PetStatusName: u.PetStatusName.String,
		CreatedAt:     u.CreatedAt.Time,
		UpdatedAt:     u.UpdatedAt.Time,
		CreatedBy:     u.CreatedBy.Int64,
		UpdatedBy:     u.UpdatedBy.Int64,
	})
}

type IPetStatusUseCase interface {
	Get(ctx context.Context, page int, limit int, search string) (*PetStatusResponse, error)
	CreatePetStatus(ctx context.Context, req PetStatusRequest) (*PetStatusResponse, error)
	UpdatePetStatus(ctx context.Context, req PetStatusRequest, id int) (*PetStatusResponse, error)
	DeletePetStatus(ctx context.Context, id int) (*PetStatusResponse, error)
}

type IPetStatusRepository interface {
	Get(ctx context.Context, page int, limit int, search string) (user []*PetStatus, err error)
	GetTotalPetStatus(ctx context.Context, search string) (total int, err error)
	CreatePetStatus(ctx context.Context, payload *PetStatus) (int64, error)
	UpdatePetStatus(ctx context.Context, payload *PetStatus, id int) error
	DeletePetStatus(ctx context.Context, id int) error
}
