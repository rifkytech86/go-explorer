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

type pointSourceUseCase struct {
	pointSourceRepo cms_domain.IPointSourceRepository
	env             *bootstrap.Env
	contextTimeout  time.Duration
}

func NewPointSourceUseCase(coinSource cms_domain.IPointSourceRepository, env *bootstrap.Env, timeout time.Duration) cms_domain.IPointSourceUseCase {
	return &pointSourceUseCase{
		pointSourceRepo: coinSource,
		contextTimeout:  timeout,
		env:             env,
	}
}

func (u *pointSourceUseCase) Get(ctx context.Context, page int, limit int, search string) (*cms_domain.PointSourceResponse, error) {
	listPetStatus, err := u.pointSourceRepo.Get(ctx, page, limit, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	totalData, err := u.pointSourceRepo.GetTotalPointSource(ctx, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.PointSourceResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Total:   totalData,
		Data:    listPetStatus,
	}
	return res, nil
}

func (u *pointSourceUseCase) CreatePointSource(ctx context.Context, req cms_domain.PointSourceRequest) (*cms_domain.PointSourceResponse, error) {
	idFromCtx := ctx.Value("x-user-id").(string)
	createdBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.PointSource{
		PointSourceName: sql.NullString{String: req.PointSourceName, Valid: true},
		CreatedBy:       sql.NullInt64{Int64: int64(createdBy), Valid: true},
	}

	_, err := u.pointSourceRepo.CreatePointSource(ctx, payload)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.PointSourceResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.PointSource{},
	}
	return res, nil
}

func (u *pointSourceUseCase) UpdatePointSource(ctx context.Context, req cms_domain.PointSourceRequest, id int) (*cms_domain.PointSourceResponse, error) {
	currentTime := time.Now()
	idFromCtx := ctx.Value("x-user-id").(string)
	updateBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.PointSource{
		PointSourceName: sql.NullString{String: req.PointSourceName, Valid: true},
		UpdatedAt:       sql.NullTime{Time: currentTime, Valid: true},
		UpdatedBy:       sql.NullInt64{Int64: int64(updateBy), Valid: true},
	}

	err := u.pointSourceRepo.UpdatePointSource(ctx, payload, id)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.PointSourceResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.PointSource{},
	}
	return res, nil
}

func (u *pointSourceUseCase) DeletePointSource(ctx context.Context, id int) (*cms_domain.PointSourceResponse, error) {
	err := u.pointSourceRepo.DeletePointSource(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.PointSourceResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.PointSource{},
	}
	return res, nil
}
