package usecase

import (
	"context"
	"gitlab.com/naonweh-studio/bubbme-backend/domain"
	"time"
)

type diaryUsecase struct {
	diaryRepository domain.DiaryRepository
	contextTimeout  time.Duration
}

func NewDiaryUsecase(diaryRepository domain.DiaryRepository, timeout time.Duration) domain.DiaryUsecase {
	return &diaryUsecase{
		diaryRepository: diaryRepository,
		contextTimeout:  timeout,
	}
}

func (ds *diaryUsecase) FetchDiary(ctx context.Context, diary *domain.DiaryRequest) ([]*domain.Diary, error) {
	ctx, cancel := context.WithTimeout(ctx, ds.contextTimeout)
	defer cancel()
	return ds.diaryRepository.FetchDiaryRepository(ctx, diary)
}

func (ds *diaryUsecase) CreateDiary(ctx context.Context, diary *domain.DiaryReq) (lastID int64, err error) {
	ctx, cancel := context.WithTimeout(ctx, ds.contextTimeout)
	defer cancel()
	return ds.diaryRepository.CreateDiaryRepository(ctx, diary)
}

func (ds *diaryUsecase) FindUserByEmail(ctx context.Context, userEmail string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, ds.contextTimeout)
	defer cancel()
	return ds.diaryRepository.FindUserRepository(ctx, userEmail)
}
func (ds *diaryUsecase) FindUserByID(ctx context.Context, userID string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, ds.contextTimeout)
	defer cancel()
	return ds.diaryRepository.FindUserByIDRepository(ctx, userID)
}

func (ds *diaryUsecase) UpdateDiary(ctx context.Context, diary *domain.DiaryReq) error {
	ctx, cancel := context.WithTimeout(ctx, ds.contextTimeout)
	defer cancel()
	return ds.diaryRepository.UpdateDiaryRepository(ctx, diary)
}

func (ds *diaryUsecase) DeleteDiary(ctx context.Context, diaryID string) error {
	ctx, cancel := context.WithTimeout(ctx, ds.contextTimeout)
	defer cancel()
	return ds.diaryRepository.DeleteDiaryRepository(ctx, diaryID)
}
