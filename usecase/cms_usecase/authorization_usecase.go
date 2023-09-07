package cms_usecase

import (
	"context"
	"errors"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain"
	"gitlab.com/naonweh-studio/bubbme-backend/internal"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/tokenutil"
	"time"
)

type authorizationUsecase struct {
	authrozationRepository domain.IUserRepository
	env                    *bootstrap.Env
	contextTimeout         time.Duration
}

func NewAuthorization(authorizationRepository domain.IUserRepository, env *bootstrap.Env, timeout time.Duration) domain.IAuthorizationUsecase {
	return &authorizationUsecase{
		authrozationRepository: authorizationRepository,
		contextTimeout:         timeout,
		env:                    env,
	}
}

func (a *authorizationUsecase) Login(ctx context.Context, email string, password string) (*domain.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	//

	user, err := a.authrozationRepository.GetUser(ctx, email)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, errors.New(internal.ErrorUnauthorization)
	}

	setUser := &domain.Users{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}

	isMatch := internal.DoPasswordsMatch(user.Password, password)

	if !isMatch {
		logger.Error("loggin not match", nil)
		return nil, errors.New(internal.ErrorUnauthorization)
	}

	// generate JWT
	token, err := tokenutil.CreateAccessTokenAdmin(setUser, a.env.AccessTokenSecret, a.env.AccessTokenExpiryHour)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, err
	}

	err = a.authrozationRepository.UpdateIsLogin(ctx, email, 1)
	if err != nil {
		logger.Error(err.Error(), nil)
		return nil, errors.New(internal.ErrorUnauthorization)
	}

	res := &domain.AuthResponse{
		Message: "succes",
		Code:    internal.USERLOGIN,
		Token:   token,
		Data: []*domain.Users{
			setUser,
		},
	}

	return res, nil
}

func (a authorizationUsecase) LogOut(ctx context.Context, email string) (*domain.AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()
	//

	user, err := a.authrozationRepository.GetUser(ctx, email)
	if err != nil {
		return nil, errors.New(internal.ErrorUnauthorization)
	}

	err = a.authrozationRepository.UpdateIsLogin(ctx, user.Email, 0)
	if err != nil {
		return nil, errors.New(internal.ErrorUnauthorization)
	}
	res := &domain.AuthResponse{
		Message: "success",
		Code:    internal.USERLOGIN,
		Data:    []*domain.Users{},
	}

	return res, nil

}
