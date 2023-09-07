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

type itemTypeUseCase struct {
	petStatusRepository cms_domain.IItemTypeRepository
	env                 *bootstrap.Env
	contextTimeout      time.Duration
}

func NewItemTypeUseCase(petStatusRepository cms_domain.IItemTypeRepository, env *bootstrap.Env, timeout time.Duration) cms_domain.IItemTypeUseCase {
	return &itemTypeUseCase{
		petStatusRepository: petStatusRepository,
		contextTimeout:      timeout,
		env:                 env,
	}
}

func (u *itemTypeUseCase) Get(ctx context.Context, page int, limit int, search string) (*cms_domain.ItemTypeResponse, error) {
	listPetStatus, err := u.petStatusRepository.Get(ctx, page, limit, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	totalData, err := u.petStatusRepository.GetTotalItemType(ctx, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.ItemTypeResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Total:   totalData,
		Data:    listPetStatus,
	}
	return res, nil
}

func (u *itemTypeUseCase) CreateItemType(ctx context.Context, req cms_domain.ItemTypeRequest) (*cms_domain.ItemTypeResponse, error) {
	idFromCtx := ctx.Value("x-user-id").(string)
	createdBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.ItemType{
		ItemTypeName: sql.NullString{String: req.ItemTypeName, Valid: true},
		CreatedBy:    sql.NullInt64{Int64: int64(createdBy), Valid: true},
	}

	_, err := u.petStatusRepository.CreateItemType(ctx, payload)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.ItemTypeResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.ItemType{},
	}
	return res, nil
}

func (u *itemTypeUseCase) UpdateItemType(ctx context.Context, req cms_domain.ItemTypeRequest, id int) (*cms_domain.ItemTypeResponse, error) {
	currentTime := time.Now()
	idFromCtx := ctx.Value("x-user-id").(string)
	updateBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.ItemType{
		ItemTypeName: sql.NullString{String: req.ItemTypeName, Valid: true},
		UpdatedAt:    sql.NullTime{Time: currentTime, Valid: true},
		UpdatedBy:    sql.NullInt64{Int64: int64(updateBy), Valid: true},
	}

	err := u.petStatusRepository.UpdateItemType(ctx, payload, id)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.ItemTypeResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.ItemType{},
	}
	return res, nil
}

func (u *itemTypeUseCase) DeleteItemType(ctx context.Context, id int) (*cms_domain.ItemTypeResponse, error) {
	err := u.petStatusRepository.DeleteItemType(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.ItemTypeResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.ItemType{},
	}
	return res, nil
}
