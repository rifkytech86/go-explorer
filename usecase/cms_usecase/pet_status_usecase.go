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

type petStatusUseCase struct {
	petStatusRepository cms_domain.IPetStatusRepository
	env                 *bootstrap.Env
	contextTimeout      time.Duration
}

func NewPetStatusUseCase(petStatusRepository cms_domain.IPetStatusRepository, env *bootstrap.Env, timeout time.Duration) cms_domain.IPetStatusUseCase {
	return &petStatusUseCase{
		petStatusRepository: petStatusRepository,
		contextTimeout:      timeout,
		env:                 env,
	}
}

func (u *petStatusUseCase) Get(ctx context.Context, page int, limit int, search string) (*cms_domain.PetStatusResponse, error) {
	listPetStatus, err := u.petStatusRepository.Get(ctx, page, limit, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	totalData, err := u.petStatusRepository.GetTotalPetStatus(ctx, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.PetStatusResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Total:   totalData,
		Data:    listPetStatus,
	}
	return res, nil
}

func (u *petStatusUseCase) CreatePetStatus(ctx context.Context, req cms_domain.PetStatusRequest) (*cms_domain.PetStatusResponse, error) {
	idFromCtx := ctx.Value("x-user-id").(string)
	createdBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.PetStatus{
		PetStatusName: sql.NullString{String: req.PetStatusName, Valid: true},
		CreatedBy:     sql.NullInt64{Int64: int64(createdBy), Valid: true},
	}

	_, err := u.petStatusRepository.CreatePetStatus(ctx, payload)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.PetStatusResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.PetStatus{},
	}
	return res, nil
}

func (u *petStatusUseCase) UpdatePetStatus(ctx context.Context, req cms_domain.PetStatusRequest, id int) (*cms_domain.PetStatusResponse, error) {
	currentTime := time.Now()
	idFromCtx := ctx.Value("x-user-id").(string)
	updateBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.PetStatus{
		PetStatusName: sql.NullString{String: req.PetStatusName, Valid: true},
		UpdatedAt:     sql.NullTime{Time: currentTime, Valid: true},
		UpdatedBy:     sql.NullInt64{Int64: int64(updateBy), Valid: true},
	}

	err := u.petStatusRepository.UpdatePetStatus(ctx, payload, id)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.PetStatusResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.PetStatus{},
	}
	return res, nil
}

func (u *petStatusUseCase) DeletePetStatus(ctx context.Context, id int) (*cms_domain.PetStatusResponse, error) {
	err := u.petStatusRepository.DeletePetStatus(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.PetStatusResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.PetStatus{},
	}
	return res, nil
}
