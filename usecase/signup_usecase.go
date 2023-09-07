package usecase

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"time"

	"gitlab.com/naonweh-studio/bubbme-backend/domain"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/tokenutil"
)

type signupUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
	cache          *bootstrap.RedisClient
}

func NewSignupUsecase(userRepository domain.UserRepository, cache *bootstrap.RedisClient, timeout time.Duration) domain.SignupUsecase {
	return &signupUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
		cache:          cache,
	}
}

func (su *signupUsecase) Create(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	_, err := su.userRepository.Create(ctx, user)
	return err
}

func (su *signupUsecase) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	return su.userRepository.GetByEmail(ctx, email)
}

func (su *signupUsecase) GetUserByPhone(c context.Context, phoneNumber string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	return su.userRepository.GetUserByPhone(ctx, phoneNumber)
}

func (su *signupUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, secret, expiry)
}

func (su *signupUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}

func (su *signupUsecase) SetOTPExpired(ctx context.Context, params domain.CacheOTP) error {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	return su.userRepository.SetExpiredToken(ctx, params)
}

func (su *signupUsecase) GetCacheVerify(ctx context.Context, param domain.VerifyOTP) (*domain.CacheOTP, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()
	res, err := su.userRepository.GetCacheVerify(ctx, param)
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("data not found")
		}
		return nil, err
	}

	return res, nil
}

func (su *signupUsecase) SetCacheVerify(ctx context.Context, params domain.CacheOTP) error {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	return su.userRepository.SetExpiredToken(ctx, params)
}

func (su *signupUsecase) DelCache(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	return su.userRepository.DeleteCacheVerify(ctx, email)
}

func (su *signupUsecase) UpdateUserVerify(ctx context.Context, userID int64, isVerify int) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	return su.userRepository.UpdateUserVerify(ctx, userID, isVerify)
}

func (su *signupUsecase) UpdateOTPUser(ctx context.Context, phoneNumber string, otp int) error {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	return su.userRepository.UpdateOTPUserByPhone(ctx, phoneNumber, otp)
}
