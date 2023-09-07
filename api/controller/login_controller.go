package controller

import (
	"encoding/json"
	"gitlab.com/naonweh-studio/bubbme-backend/internal"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain"
)

type LoginController struct {
	LoginUsecase domain.LoginUsecase
	Env          *bootstrap.Env
	Cache        *bootstrap.RedisClient
}

// Login handles user login.
// @Summary Login user
// @Description Login user and get JWT token
// @Success 200 {object} domain.LoginResponse
// @Router /api/login [post]
func (lc *LoginController) Login(c *gin.Context) {
	var request domain.LoginRequest

	err := c.ShouldBind(&request)
	if err != nil {
		logger.Error(err.Error(), nil)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	key := "otp:" + request.PhoneNumber
	//user, err := lc.LoginUsecase.GetUserByEmail(c, request.Email)
	//if err != nil {
	//	c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "User not found with the given email"})
	//	return
	//}

	user, err := lc.LoginUsecase.GetUserByPhone(c, request.PhoneNumber)
	if err != nil {
		logger.Error(err.Error(), nil)
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "User not found."})
		return
	}

	if user.UserIsVerify.Int64 == 0 {
		c.JSON(http.StatusConflict, domain.ErrorResponse{Message: "User not verify yet"})
		return
	}

	//if bcrypt.CompareHashAndPassword([]byte(user.UserPassword.String), []byte(request.Password)) != nil {
	//	c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Invalid credentials"})
	//	return
	//}

	_, err = lc.Cache.GetRedis(c, key)
	if err != nil {
		// generate OTP
		otp, err := internal.GenerateOTP(6)
		if err != nil {
			logger.Error(err.Error(), nil)
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Error Generated OTP"})
			return
		}

		err = lc.LoginUsecase.UpdateOTPUser(c, request.PhoneNumber, otp)
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

		err = lc.Cache.SetRedis(c, key, string(data), 3)
		if err != nil {
			logger.Error(err.Error(), nil)
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Error Generated OTP"})
			return
		}
	}

	accessToken, err := lc.LoginUsecase.CreateAccessToken(&user, lc.Env.AccessTokenSecret, lc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	refreshToken, err := lc.LoginUsecase.CreateRefreshToken(&user, lc.Env.RefreshTokenSecret, lc.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	loginResponse := domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, loginResponse)
}
