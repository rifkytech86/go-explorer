package cms_usecase

import (
	"context"
	"database/sql"
	"errors"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain/cms_domain"
	"gitlab.com/naonweh-studio/bubbme-backend/internal"
	"strconv"
	"time"
)

type usersAdminUseCase struct {
	userAdminRepository cms_domain.IUserAdminRepository
	env                 *bootstrap.Env
	contextTimeout      time.Duration
}

func NewUsersAdminUseCase(userAdminRepository cms_domain.IUserAdminRepository, env *bootstrap.Env, timeout time.Duration) cms_domain.IUsersAdminUseCase {
	return &usersAdminUseCase{
		userAdminRepository: userAdminRepository,
		contextTimeout:      timeout,
		env:                 env,
	}
}

func (u *usersAdminUseCase) GetUserAdmin(ctx context.Context, page int, limit int, search string) (*cms_domain.UserAdminResponse, error) {

	users, err := u.userAdminRepository.GetUserAdmin(ctx, page, limit, search)
	if err != nil {
		return nil, err
	}

	totalData, err := u.userAdminRepository.GetTotalUserAdmin(ctx)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.UserAdminResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Total:   totalData,
		Data:    users,
	}
	return res, nil
}

func (u *usersAdminUseCase) CreateUserAdmin(ctx context.Context, req cms_domain.UserRequest) (*cms_domain.UserAdminResponse, error) {
	hasPassword, err := internal.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	idFromCtx := ctx.Value("x-user-id").(string)
	idCreatedBy, _ := strconv.Atoi(idFromCtx)
	payload := &cms_domain.Users{
		Email:     req.Email,
		Name:      req.Name,
		IsLogin:   0,
		Username:  req.Username,
		Password:  hasPassword,
		CreatedBy: idCreatedBy,
	}

	checkEmail, err := u.userAdminRepository.GetUserAdminByEmail(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if checkEmail != nil && checkEmail.Email != "" {
		return nil, errors.New("email already registered")
	}

	_, err = u.userAdminRepository.CreateUserAdmin(ctx, payload)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.UserAdminResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.Users{},
	}
	return res, nil
}

func (u *usersAdminUseCase) UpdateUserAdmin(ctx context.Context, req cms_domain.UserRequest, id int) (*cms_domain.UserAdminResponse, error) {
	payload := &cms_domain.Users{
		Email:    req.Email,
		Name:     req.Name,
		IsLogin:  0,
		Username: req.Username,
	}

	err := u.userAdminRepository.UpdateUserAdmin(ctx, payload, id)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.UserAdminResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.Users{},
	}
	return res, nil
}

func (u *usersAdminUseCase) DeleteUserAdmin(ctx context.Context, id int) (*cms_domain.UserAdminResponse, error) {
	err := u.userAdminRepository.DeleteUserAdmin(ctx, id)
	if err != nil {
		return nil, err
	}

	res := &cms_domain.UserAdminResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*cms_domain.Users{},
	}
	return res, nil
}
