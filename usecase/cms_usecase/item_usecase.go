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

type itemUseCase struct {
	petStatusRepository cms_domain.IItemRepository
	env                 *bootstrap.Env
	contextTimeout      time.Duration
}

func NewItemUseCase(petStatusRepository cms_domain.IItemRepository, env *bootstrap.Env, timeout time.Duration) cms_domain.IItemUseCase {
	return &itemUseCase{
		petStatusRepository: petStatusRepository,
		contextTimeout:      timeout,
		env:                 env,
	}
}

func (u *itemUseCase) Get(ctx context.Context, page int, limit int, search string) (*cms_domain.ItemResponse, error) {
	listPetStatus, err := u.petStatusRepository.Get(ctx, page, limit, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	totalData, err := u.petStatusRepository.GetTotalItem(ctx, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.ItemResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Total:   totalData,
		Data:    listPetStatus,
	}
	return res, nil
}

func (u *itemUseCase) CreateItemType(ctx context.Context, req cms_domain.ItemRequest) (*cms_domain.ItemResponse, error) {
	idFromCtx := ctx.Value("x-user-id").(string)
	createdBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.Item{
		ItemName:   sql.NullString{String: req.ItemName, Valid: true},
		ItemDesc:   sql.NullString{String: req.ItemDesc, Valid: true},
		ItemValue:  sql.NullInt64{Int64: req.ItemValue, Valid: true},
		ItemTypeID: sql.NullInt64{Int64: req.ItemTypeID, Valid: true},
		CreatedBy:  sql.NullInt64{Int64: int64(createdBy), Valid: true},
	}

	_, err := u.petStatusRepository.CreateItem(ctx, payload)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.ItemResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.Item{},
	}
	return res, nil
}

func (u *itemUseCase) UpdateItemType(ctx context.Context, req cms_domain.ItemRequest, id int) (*cms_domain.ItemResponse, error) {
	currentTime := time.Now()
	idFromCtx := ctx.Value("x-user-id").(string)
	updateBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.Item{
		ItemName:   sql.NullString{String: req.ItemName, Valid: true},
		ItemDesc:   sql.NullString{String: req.ItemDesc, Valid: true},
		ItemValue:  sql.NullInt64{Int64: req.ItemValue, Valid: true},
		ItemTypeID: sql.NullInt64{Int64: req.ItemTypeID, Valid: true},
		UpdatedAt:  sql.NullTime{Time: currentTime, Valid: true},
		UpdatedBy:  sql.NullInt64{Int64: int64(updateBy), Valid: true},
	}

	err := u.petStatusRepository.UpdateItem(ctx, payload, id)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.ItemResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.Item{},
	}
	return res, nil
}

func (u *itemUseCase) DeleteItemType(ctx context.Context, id int) (*cms_domain.ItemResponse, error) {
	err := u.petStatusRepository.DeleteItem(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.ItemResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.Item{},
	}
	return res, nil
}
