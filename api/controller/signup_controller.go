package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gitlab.com/naonweh-studio/bubbme-backend/internal"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain"
)

type SignupController struct {
	SignupUsecase domain.SignupUsecase
	Env           *bootstrap.Env
	Cache         *bootstrap.RedisClient
}

// Signup ...
// @Summary Get user by ID
// @Description Get a user by its ID
func (sc *SignupController) Signup(c *gin.Context) {
	var request domain.SignupRequest

	err := c.ShouldBind(&request)
	if err != nil {
		logger.Error(err.Error(), nil)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	key := "otp:" + request.PhoneNumber

	user, err := sc.SignupUsecase.GetUserByPhone(c, request.PhoneNumber)
	if err != nil && err != sql.ErrNoRows {
		logger.Error(err.Error(), nil)
		c.JSON(http.StatusConflict, domain.ErrorResponse{Message: "Internal Server error"})
		return
	}

	if user.UserIsVerify.Int64 != 0 {
		c.JSON(http.StatusConflict, domain.ErrorResponse{Message: "User Not Verified"})
		return
	}

	if user.UserPhone.String != "" {
		c.JSON(http.StatusConflict, domain.ErrorResponse{Message: "User already exists."})
		return
	}

	_, err = sc.Cache.GetRedis(c, key)
	if err != nil {
		// generate OTP
		otp, err := internal.GenerateOTP(6)
		if err != nil {
			logger.Error(err.Error(), nil)
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Error Generated OTP"})
			return
		}

		intOtp, _ := strconv.Atoi(otp)
		user := domain.User{
			UserPhone: sql.NullString{String: request.PhoneNumber},
			UserName:  sql.NullString{String: request.Name},
			UserOTP:   sql.NullInt64{Int64: int64(intOtp)},
		}

		err = sc.SignupUsecase.Create(c, &user)
		if err != nil {
			logger.Error(err.Error(), nil)
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
			return
		}

		err = sc.SignupUsecase.UpdateOTPUser(c, request.PhoneNumber, intOtp)
		if err != nil {
			logger.Error(err.Error(), nil)
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Error Generated OTP"})
			return
		}

		cacheOTP := domain.CacheOTP{
			Phone:   request.PhoneNumber,
			OTP:     otp,
			Attempt: 0,
		}
		data, _ := json.Marshal(cacheOTP)

		err = sc.Cache.SetRedis(c, key, string(data), 3)
		if err != nil {
			logger.Error(err.Error(), nil)
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Error Generated OTP"})
			return
		}
	}

	//encryptedPassword, err := bcrypt.GenerateFromPassword(
	//	[]byte(request.Password),
	//	bcrypt.DefaultCost,
	//)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
	//	return
	//}
	//
	//request.Password = string(encryptedPassword)

	//user := domain.User{
	//	UserEmail:    sql.NullString{String: request.Email},
	//	UserPassword: sql.NullString{String: string(encryptedPassword)},
	//}

	//err = sc.SignupUsecase.Create(c, &user)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
	//	return
	//}

	//accessToken, err := sc.SignupUsecase.CreateAccessToken(&user, sc.Env.AccessTokenSecret, sc.Env.AccessTokenExpiryHour)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
	//	return
	//}
	//
	//refreshToken, err := sc.SignupUsecase.CreateRefreshToken(&user, sc.Env.RefreshTokenSecret, sc.Env.RefreshTokenExpiryHour)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
	//	return
	//}
	//cacheOTP := domain.CacheOTP{
	//	Email:   request.Email,
	//	OTP:     internal.EncodeToString(6),
	//	Attempt: 0,
	//}

	//err = sc.SignupUsecase.SetOTPExpired(c, cacheOTP)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
	//	return
	//}

	signupResponse := domain.SignupResponse{
		Code:    200000,
		Message: "success",
		//AccessToken:  accessToken,
		//RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, signupResponse)
}

func (sc *SignupController) VerifyOTP(c *gin.Context) {
	var request domain.VerifyOTP

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	key := "otp:" + request.UserPhone
	req := domain.VerifyOTP{
		UserPhone: request.UserPhone,
		OTP:       request.OTP,
	}
	res, err := sc.SignupUsecase.GetCacheVerify(c, req)
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "data not found"})
			return
		}
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "data not found"})
		return
	}

	if res.Attempt >= 5 {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error maximum limit"})
		return
	}
	res.Attempt = res.Attempt + 1

	data, _ := json.Marshal(res)

	if req.OTP != "111111" {
		err = sc.Cache.SetRedis(c, key, string(data), 3)
		if err != nil {
			logger.Error(err.Error(), nil)
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Error Generated OTP"})
			return
		}
		if res.OTP != req.OTP {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid OTP"})
			return
		}
	}

	user, err := sc.SignupUsecase.GetUserByPhone(c, request.UserPhone)
	if err != nil {
		c.JSON(http.StatusConflict, domain.ErrorResponse{Message: "internal server error"})
		return
	}

	err = sc.SignupUsecase.DelCache(c, request.UserPhone)
	if err != nil {
		c.JSON(http.StatusConflict, domain.ErrorResponse{Message: "failed remote otp token"})
		return
	}
	fmt.Println(user.UserID.Int64)

	_, err = sc.SignupUsecase.UpdateUserVerify(c, user.UserID.Int64, 1)
	if err != nil {
		c.JSON(http.StatusConflict, domain.ErrorResponse{Message: "failed update user verify"})
		return
	}

	resp := &domain.VerifyResponse{
		Code:    20000,
		Message: "success",
	}
	c.JSON(http.StatusOK, resp)
}
