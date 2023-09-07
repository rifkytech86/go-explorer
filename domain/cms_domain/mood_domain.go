package cms_domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

//type AuthAppRequest struct {
//	Password string `json:"password"`
//	Email    string `json:"email"`
//}

type MoodRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	DOB      string `json:"dob"`
	School   string `json:"school"`
}

type MoodResponse struct {
	Message string  `json:"msg"`
	Code    int64   `json:"code"`
	Total   int     `json:"total"`
	Data    []*Mood `json:"data"`
}

type Mood struct {
	ID        sql.NullInt64  `json:"id"`
	Name      sql.NullString `json:"name"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAT *time.Time     `json:"updated_at"`
	CreatedBy sql.NullInt64  `json:"created_by"`
	UpdatedBy sql.NullInt64  `json:"updated_by"`
}

func (u Mood) MarshalJSON() ([]byte, error) {
	type Alias Mood
	return json.Marshal(&struct {
		Alias
		ID        interface{} `json:"id,omitempty"`
		Name      interface{} `json:"name,omitempty"`
		CreatedBy interface{} `json:"created_by,omitempty"`
		UpdatedBy interface{} `json:"updated_by,omitempty"`
	}{
		Alias:     (Alias)(u),
		ID:        u.ID.Int64,
		Name:      u.Name.String,
		CreatedBy: u.CreatedBy.Int64,
		UpdatedBy: u.UpdatedBy.Int64,
	})
}

type IMoodUseCase interface {
	Get(ctx context.Context, page int, limit int, search string) (*MoodResponse, error)
	CreateUserApp(ctx context.Context, req MoodRequest) (*MoodResponse, error)
	UpdateUserApp(ctx context.Context, req MoodRequest, id int) (*MoodResponse, error)
	DeleteUserApp(ctx context.Context, id int) (*MoodResponse, error)
}

type IMoodRepository interface {
	Get(ctx context.Context, page int, limit int, search string) (user []*Mood, err error)
	GetTotalUserMood(ctx context.Context) (total int, err error)
	CreateUserApp(ctx context.Context, payload *Mood) (lastID int64, err error)
	UpdateUserApp(ctx context.Context, payload *Mood, id int) error
	DeleteUserApp(ctx context.Context, id int) error
}
