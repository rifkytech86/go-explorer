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

type coinSourceUseCase struct {
	coinSourceRepo cms_domain.ICoinSourceRepository
	env            *bootstrap.Env
	contextTimeout time.Duration
}

func NewCoinSourceUseCase(coinSource cms_domain.ICoinSourceRepository, env *bootstrap.Env, timeout time.Duration) cms_domain.ICoinSourceUseCase {
	return &coinSourceUseCase{
		coinSourceRepo: coinSource,
		contextTimeout: timeout,
		env:            env,
	}
}

func (u *coinSourceUseCase) Get(ctx context.Context, page int, limit int, search string) (*cms_domain.CoinSourceResponse, error) {
	listPetStatus, err := u.coinSourceRepo.Get(ctx, page, limit, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	totalData, err := u.coinSourceRepo.GetTotalCoinSource(ctx, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.CoinSourceResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Total:   totalData,
		Data:    listPetStatus,
	}
	return res, nil
}

func (u *coinSourceUseCase) CreateCoinSource(ctx context.Context, req cms_domain.CoinSourceRequest) (*cms_domain.CoinSourceResponse, error) {
	idFromCtx := ctx.Value("x-user-id").(string)
	createdBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.CoinSource{
		CoinSourceName: sql.NullString{String: req.CoinSourceName, Valid: true},
		CreatedBy:      sql.NullInt64{Int64: int64(createdBy), Valid: true},
	}

	_, err := u.coinSourceRepo.CreateCoinSource(ctx, payload)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.CoinSourceResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.CoinSource{},
	}
	return res, nil
}

func (u *coinSourceUseCase) UpdateCoinSource(ctx context.Context, req cms_domain.CoinSourceRequest, id int) (*cms_domain.CoinSourceResponse, error) {
	currentTime := time.Now()
	idFromCtx := ctx.Value("x-user-id").(string)
	updateBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.CoinSource{
		CoinSourceName: sql.NullString{String: req.CoinSourceName, Valid: true},
		UpdatedAt:      sql.NullTime{Time: currentTime, Valid: true},
		UpdatedBy:      sql.NullInt64{Int64: int64(updateBy), Valid: true},
	}

	err := u.coinSourceRepo.UpdateCoinSource(ctx, payload, id)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.CoinSourceResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.CoinSource{},
	}
	return res, nil
}

func (u *coinSourceUseCase) DeleteCoinSource(ctx context.Context, id int) (*cms_domain.CoinSourceResponse, error) {
	err := u.coinSourceRepo.DeleteCoinSource(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.CoinSourceResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.CoinSource{},
	}
	return res, nil
}
