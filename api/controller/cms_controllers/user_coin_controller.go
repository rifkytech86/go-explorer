package cms_controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/naonweh-studio/bubbme-backend/domain"
	"gitlab.com/naonweh-studio/bubbme-backend/domain/cms_domain"
	"gitlab.com/naonweh-studio/bubbme-backend/internal/logger"
	"net/http"
	"strconv"
)

type UserCoinController struct {
	UserCoinUseCase cms_domain.IUserCoinUseCase
}

func (u *UserCoinController) Get(ctx *gin.Context) {
	page := ctx.Query("page")
	if page == "" {
		page = "1"
	}
	intPage, err := strconv.Atoi(page)
	if err != nil {
		logger.Error(err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	limit := ctx.Query("limit")
	if limit == "" {
		limit = "10"
	}
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		logger.Error(err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	search := ctx.Query("q")
	if err != nil {
		logger.Error(err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	resp, err := u.UserCoinUseCase.Get(ctx, intPage, intLimit, search)
	if err != nil {
		logger.Error(err.Error(), nil)
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (u *UserCoinController) Create(ctx *gin.Context) {
	var request cms_domain.UserCoinRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	resp, err := u.UserCoinUseCase.Create(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (u *UserCoinController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		logger.Error("error update user admin", nil)
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error update user admin"})
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error update user admin"})
		return
	}
	var request cms_domain.UserCoinRequest
	err = ctx.ShouldBind(&request)
	if err != nil {
		logger.Error(err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	resp, err := u.UserCoinUseCase.Update(ctx, request, intID)
	if err != nil {
		logger.Error(err.Error(), nil)
		ctx.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (u *UserCoinController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		logger.Error("error delete user admin", nil)
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error delete user admin"})
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error delete user admin"})
		return
	}

	resp, err := u.UserCoinUseCase.Delete(ctx, intID)
	if err != nil {
		logger.Error(err.Error(), nil)
		ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
