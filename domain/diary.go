package domain

import (
	"context"
	"database/sql"
)

type DiaryRequest struct {
	DiaryID        int64  `json:"diary_id"`
	UserID         int64  `json:"user_id"`
	DiaryDesc      string `json:"diary_desc"`
	DiaryLinkVideo string `json:"diary_link_video"`
	DiaryIsActive  bool   `json:"diary_is_active"`
	Limit          int64  `json:"limit"`
}

type DiaryReq struct {
	UserEmail      string `json:"user_email"`
	DiaryDesc      string `json:"diary_desc"`
	UserID         int64  `json:"user_id"`
	DiaryLinkVideo string `json:"diary_link_video"`
}

type DiaryResponse struct {
	Code    int64    `json:"code"`
	Message string   `json:"msg"`
	Data    []*Diary `json:"data"`
}

type Diary struct {
	DiaryID        int64        `json:"diary_id"`
	UserID         int64        `json:"user_id"`
	DiaryDesc      string       `json:"diary_desc"`
	DiaryLinkVideo string       `json:"diary_link_video"`
	DiaryIsActive  bool         `json:"diary_is_active"`
	CreatedAt      sql.NullTime `json:"created_at"`
	UpdatedAt      sql.NullTime `json:"updated_at"`
}

type DiaryUsecase interface {
	FetchDiary(c context.Context, diary *DiaryRequest) ([]*Diary, error)
	CreateDiary(c context.Context, diary *DiaryReq) (lastID int64, err error)
	FindUserByEmail(c context.Context, userEmail string) (User, error)
	FindUserByID(c context.Context, userEmail string) (User, error)
	UpdateDiary(c context.Context, diary *DiaryReq) error
	DeleteDiary(c context.Context, diaryID string) error
}

type DiaryRepository interface {
	FetchDiaryRepository(c context.Context, diary *DiaryRequest) ([]*Diary, error)
	CreateDiaryRepository(c context.Context, diary *DiaryReq) (lastID int64, err error)
	FindUserRepository(c context.Context, userEmail string) (User, error)
	FindUserByIDRepository(c context.Context, userID string) (User, error)
	UpdateDiaryRepository(c context.Context, diary *DiaryReq) error
	DeleteDiaryRepository(c context.Context, diaryID string) error
}
