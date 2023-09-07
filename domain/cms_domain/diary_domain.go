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

type DiaryRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	DOB      string `json:"dob"`
	School   string `json:"school"`
}

type DiaryResponse struct {
	Message string   `json:"msg"`
	Code    int64    `json:"code"`
	Total   int      `json:"total"`
	Data    []*Diary `json:"data"`
}

type Diary struct {
	ID             sql.NullInt64  `json:"diary_id"`
	UserID         sql.NullInt64  `json:"user_id"`
	DiaryDesc      sql.NullString `json:"diary_desc"`
	DiaryLinkVideo sql.NullString `json:"diary_link_video"`
	DiaryIsActive  sql.NullInt64  `json:"diary_is_active"`
	UpdatedAT      *time.Time     `json:"updated_at"`
	CreatedAt      *time.Time     `json:"created_at"`
}

func (u Diary) MarshalJSON() ([]byte, error) {
	type Alias Diary
	return json.Marshal(&struct {
		Alias
		ID             interface{} `json:"diary_id,omitempty"`
		UserID         interface{} `json:"user_id,omitempty"`
		DiaryDesc      interface{} `json:"diary_desc,omitempty"`
		DiaryLinkVideo interface{} `json:"diary_link_video,omitempty"`
		DiaryIsActive  interface{} `json:"diary_is_active,omitempty"`
		UpdatedAT      interface{} `json:"updated_at,omitempty"`
		CreatedAT      interface{} `json:"created_at,omitempty"`
	}{
		Alias:          (Alias)(u),
		ID:             u.ID.Int64,
		UserID:         u.UserID.Int64,
		DiaryDesc:      u.DiaryDesc.String,
		DiaryLinkVideo: u.DiaryLinkVideo.String,
	})
}

type IDiaryUseCase interface {
	Get(ctx context.Context, page int, limit int, search string) (*DiaryResponse, error)
	CreateUserApp(ctx context.Context, req DiaryRequest) (*DiaryResponse, error)
	UpdateUserApp(ctx context.Context, req DiaryRequest, id int) (*DiaryResponse, error)
	DeleteUserApp(ctx context.Context, id int) (*DiaryResponse, error)
}

type IDiaryRepository interface {
	Get(ctx context.Context, page int, limit int, search string) (user []*Diary, err error)
	GetTotalUserMood(ctx context.Context) (total int, err error)
	CreateUserApp(ctx context.Context, payload *Mood) (lastID int64, err error)
	UpdateUserApp(ctx context.Context, payload *Mood, id int) error
	DeleteUserApp(ctx context.Context, id int) error
}
