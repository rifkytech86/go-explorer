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

type moodUseCase struct {
	moodRepository cms_domain.IMoodRepository
	env            *bootstrap.Env
	contextTimeout time.Duration
}

func NewMoodUseCase(userAdminRepository cms_domain.IMoodRepository, env *bootstrap.Env, timeout time.Duration) cms_domain.IMoodUseCase {
	return &moodUseCase{
		moodRepository: userAdminRepository,
		contextTimeout: timeout,
		env:            env,
	}
}

func (u *moodUseCase) Get(ctx context.Context, page int, limit int, search string) (*cms_domain.MoodResponse, error) {
	users, err := u.moodRepository.Get(ctx, page, limit, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	totalData, err := u.moodRepository.GetTotalUserMood(ctx)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.MoodResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Total:   totalData,
		Data:    users,
	}
	return res, nil
}

func (u *moodUseCase) CreateUserApp(ctx context.Context, req cms_domain.MoodRequest) (*cms_domain.MoodResponse, error) {
	idFromCtx := ctx.Value("x-user-id").(string)
	createdBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.Mood{
		Name:      sql.NullString{String: req.Name, Valid: true},
		CreatedBy: sql.NullInt64{Int64: int64(createdBy), Valid: true},
	}

	_, err := u.moodRepository.CreateUserApp(ctx, payload)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.MoodResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.Mood{},
	}
	return res, nil
}

func (u *moodUseCase) UpdateUserApp(ctx context.Context, req cms_domain.MoodRequest, id int) (*cms_domain.MoodResponse, error) {

	currentTime := time.Now()
	idFromCtx := ctx.Value("x-user-id").(string)
	updateBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.Mood{
		Name:      sql.NullString{String: req.Name, Valid: true},
		UpdatedAT: &currentTime,
		UpdatedBy: sql.NullInt64{Int64: int64(updateBy), Valid: true},
	}

	err := u.moodRepository.UpdateUserApp(ctx, payload, id)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.MoodResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.Mood{},
	}
	return res, nil
}

func (u *moodUseCase) DeleteUserApp(ctx context.Context, id int) (*cms_domain.MoodResponse, error) {
	err := u.moodRepository.DeleteUserApp(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.MoodResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.Mood{},
	}
	return res, nil
}
