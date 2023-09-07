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

type userCoinUseCase struct {
	userCoinRepo   cms_domain.IUserCoinRepository
	env            *bootstrap.Env
	contextTimeout time.Duration
}

func NewUserCoinUseCase(userCoinRepo cms_domain.IUserCoinRepository, env *bootstrap.Env, timeout time.Duration) cms_domain.IUserCoinUseCase {
	return &userCoinUseCase{
		userCoinRepo:   userCoinRepo,
		contextTimeout: timeout,
		env:            env,
	}
}

func (u *userCoinUseCase) Get(ctx context.Context, page int, limit int, search string) (*cms_domain.UserCoinResponse, error) {
	listPetStatus, err := u.userCoinRepo.Get(ctx, page, limit, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	totalData, err := u.userCoinRepo.GetTotal(ctx, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.UserCoinResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Total:   totalData,
		Data:    listPetStatus,
	}
	return res, nil
}

func (u *userCoinUseCase) Create(ctx context.Context, req cms_domain.UserCoinRequest) (*cms_domain.UserCoinResponse, error) {
	idFromCtx := ctx.Value("x-user-id").(string)
	createdBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.UserCoin{
		UserID:        sql.NullInt64{Int64: req.UserID, Valid: true},
		UserCoinValue: sql.NullInt64{Int64: req.UserCoinValue, Valid: true},
		CreatedBy:     sql.NullInt64{Int64: int64(createdBy), Valid: true},
	}

	_, err := u.userCoinRepo.Create(ctx, payload)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.UserCoinResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.UserCoin{},
	}
	return res, nil
}

func (u *userCoinUseCase) Update(ctx context.Context, req cms_domain.UserCoinRequest, id int) (*cms_domain.UserCoinResponse, error) {
	currentTime := time.Now()
	idFromCtx := ctx.Value("x-user-id").(string)
	updateBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.UserCoin{
		UserID:        sql.NullInt64{Int64: req.UserID, Valid: true},
		UserCoinValue: sql.NullInt64{Int64: req.UserCoinValue, Valid: true},
		UpdatedAt:     sql.NullTime{Time: currentTime, Valid: true},
		UpdatedBy:     sql.NullInt64{Int64: int64(updateBy), Valid: true},
	}

	err := u.userCoinRepo.Update(ctx, payload, id)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.UserCoinResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.UserCoin{},
	}
	return res, nil
}

func (u *userCoinUseCase) Delete(ctx context.Context, id int) (*cms_domain.UserCoinResponse, error) {
	err := u.userCoinRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.UserCoinResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.UserCoin{},
	}
	return res, nil
}
