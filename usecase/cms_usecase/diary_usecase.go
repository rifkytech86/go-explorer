package cms_usecase

import (
	"context"
	"database/sql"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain/cms_domain"
	"gitlab.com/naonweh-studio/bubbme-backend/internal"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"strconv"
	"time"
)

type diaryUseCase struct {
	diaryRepository cms_domain.IDiaryRepository
	env             *bootstrap.Env
	contextTimeout  time.Duration
}

func NewDiaryUseCase(userAdminRepository cms_domain.IDiaryRepository, env *bootstrap.Env, timeout time.Duration) cms_domain.IDiaryUseCase {
	return &diaryUseCase{
		diaryRepository: userAdminRepository,
		contextTimeout:  timeout,
		env:             env,
	}
}

func (u *diaryUseCase) Get(ctx context.Context, page int, limit int, search string) (*cms_domain.DiaryResponse, error) {
	users, err := u.diaryRepository.Get(ctx, page, limit, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	totalData, err := u.diaryRepository.GetTotalUserMood(ctx)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.DiaryResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Total:   totalData,
		Data:    users,
	}
	return res, nil
}

func (u *diaryUseCase) CreateUserApp(ctx context.Context, req cms_domain.DiaryRequest) (*cms_domain.DiaryResponse, error) {
	idFromCtx := ctx.Value("x-user-id").(string)
	createdBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.Mood{
		Name:      sql.NullString{String: req.Name, Valid: true},
		CreatedBy: sql.NullInt64{Int64: int64(createdBy), Valid: true},
	}

	_, err := u.diaryRepository.CreateUserApp(ctx, payload)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.DiaryResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.Diary{},
	}
	return res, nil
}

func (u *diaryUseCase) UpdateUserApp(ctx context.Context, req cms_domain.DiaryRequest, id int) (*cms_domain.DiaryResponse, error) {

	currentTime := time.Now()
	idFromCtx := ctx.Value("x-user-id").(string)
	updateBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.Mood{
		Name:      sql.NullString{String: req.Name, Valid: true},
		UpdatedAT: &currentTime,
		UpdatedBy: sql.NullInt64{Int64: int64(updateBy), Valid: true},
	}

	err := u.diaryRepository.UpdateUserApp(ctx, payload, id)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.DiaryResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.Diary{},
	}
	return res, nil
}

func (u *diaryUseCase) DeleteUserApp(ctx context.Context, id int) (*cms_domain.DiaryResponse, error) {
	err := u.diaryRepository.DeleteUserApp(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.DiaryResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.Diary{},
	}
	return res, nil
}
