package cms_usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain/cms_domain"
	"gitlab.com/naonweh-studio/bubbme-backend/internal"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"strconv"
	"time"
)

type usersAppUseCase struct {
	userAppRepository cms_domain.IUserAppRepository
	env               *bootstrap.Env
	contextTimeout    time.Duration
}

func NewUsersAppUseCase(userAdminRepository cms_domain.IUserAppRepository, env *bootstrap.Env, timeout time.Duration) cms_domain.IUsersAppUseCase {
	return &usersAppUseCase{
		userAppRepository: userAdminRepository,
		contextTimeout:    timeout,
		env:               env,
	}
}

func (u *usersAppUseCase) GetUserApp(ctx context.Context, page int, limit int, search string) (*cms_domain.UserAppResponse, error) {
	users, err := u.userAppRepository.GetUserApp(ctx, page, limit, search)
	if err != nil {
		logger.Info(err.Error(), nil)
		return nil, err
	}
	totalData, err := u.userAppRepository.GetTotalUserApp(ctx)
	if err != nil {
		logger.Info(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.UserAppResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Total:   totalData,
		Data:    users,
	}
	return res, nil
}

func (u *usersAppUseCase) CreateUserApp(ctx context.Context, req cms_domain.UserAppRequest) (*cms_domain.UserAppResponse, error) {
	idFromCtx := ctx.Value("x-user-id").(string)
	createdBy, _ := strconv.Atoi(idFromCtx)

	checkPhone, err := u.userAppRepository.GetUserAppByPhone(ctx, req.UserPhone)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	if checkPhone != nil && checkPhone.UserPhone.String != "" {
		return nil, errors.New("phone number already exist")
	}

	fmt.Println("=====", int64(createdBy))
	payload := &cms_domain.UsersApp{
		UserName:      sql.NullString{String: req.Username},
		UserPhone:     sql.NullString{String: req.UserPhone},
		UserISVerify:  sql.NullInt64{Int64: req.UserIsVerify},
		UserCreatedBy: sql.NullInt64{Int64: int64(createdBy)},
	}

	_, err = u.userAppRepository.CreateUserApp(ctx, payload)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.UserAppResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.UsersApp{},
	}
	return res, nil
}

func (u *usersAppUseCase) UpdateUserApp(ctx context.Context, req cms_domain.UserAppRequest, id int) (*cms_domain.UserAppResponse, error) {
	idFromCtx := ctx.Value("x-user-id").(string)
	updateBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.UsersApp{
		UserName:      sql.NullString{String: req.Username, Valid: true},
		UserPhone:     sql.NullString{String: req.UserPhone, Valid: true},
		UserISVerify:  sql.NullInt64{Int64: req.UserIsVerify, Valid: true},
		UserUpdatedBy: sql.NullInt64{Int64: int64(updateBy), Valid: true},
	}

	err := u.userAppRepository.UpdateUserApp(ctx, payload, id)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	res := &cms_domain.UserAppResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.UsersApp{},
	}
	return res, nil
}

func (u *usersAppUseCase) DeleteUserApp(ctx context.Context, id int) (*cms_domain.UserAppResponse, error) {
	err := u.userAppRepository.DeleteUserApp(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.UserAppResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.UsersApp{},
	}
	return res, nil
}
