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

type userPointUseCase struct {
	userPointRepo  cms_domain.IUserPointRepository
	env            *bootstrap.Env
	contextTimeout time.Duration
}

func NewUserPointCase(userPointRepo cms_domain.IUserPointRepository, env *bootstrap.Env, timeout time.Duration) cms_domain.IUserPointUseCase {
	return &userPointUseCase{
		userPointRepo:  userPointRepo,
		contextTimeout: timeout,
		env:            env,
	}
}

func (u *userPointUseCase) Get(ctx context.Context, page int, limit int, search string) (*cms_domain.UserPointResponse, error) {
	listPetStatus, err := u.userPointRepo.Get(ctx, page, limit, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	totalData, err := u.userPointRepo.GetTotal(ctx, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.UserPointResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Total:   totalData,
		Data:    listPetStatus,
	}
	return res, nil
}

func (u *userPointUseCase) Create(ctx context.Context, req cms_domain.UserPointRequest) (*cms_domain.UserPointResponse, error) {
	idFromCtx := ctx.Value("x-user-id").(string)
	createdBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.UserPoint{
		UserID:         sql.NullInt64{Int64: req.UserID, Valid: true},
		UserPointValue: sql.NullInt64{Int64: req.UserPointValue, Valid: true},
		CreatedBy:      sql.NullInt64{Int64: int64(createdBy), Valid: true},
	}

	_, err := u.userPointRepo.Create(ctx, payload)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.UserPointResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.UserPoint{},
	}
	return res, nil
}

func (u *userPointUseCase) Update(ctx context.Context, req cms_domain.UserPointRequest, id int) (*cms_domain.UserPointResponse, error) {
	currentTime := time.Now()
	idFromCtx := ctx.Value("x-user-id").(string)
	updateBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.UserPoint{
		UserID:         sql.NullInt64{Int64: req.UserID, Valid: true},
		UserPointValue: sql.NullInt64{Int64: req.UserPointValue, Valid: true},
		UpdatedAt:      sql.NullTime{Time: currentTime, Valid: true},
		UpdatedBy:      sql.NullInt64{Int64: int64(updateBy), Valid: true},
	}

	err := u.userPointRepo.Update(ctx, payload, id)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.UserPointResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.UserPoint{},
	}
	return res, nil
}

func (u *userPointUseCase) Delete(ctx context.Context, id int) (*cms_domain.UserPointResponse, error) {
	err := u.userPointRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.UserPointResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.UserPoint{},
	}
	return res, nil
}
