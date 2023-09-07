package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/naonweh-studio/bubbme-backend/bootstrap"
	"gitlab.com/naonweh-studio/bubbme-backend/domain"
	"net/http"
)

type DiaryController struct {
	DiaryUsecase domain.DiaryUsecase
	Env          *bootstrap.Env
}

func (d *DiaryController) FetchDiary(ctx *gin.Context) {
	var request domain.DiaryRequest
	_ = ctx.Request.URL.Query()

	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	listData, err := d.DiaryUsecase.FetchDiary(ctx, &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	diaryResponse := domain.DiaryResponse{
		Code:    200000,
		Message: "success",
		Data:    listData,
	}
	ctx.JSON(http.StatusOK, diaryResponse)

}

func (d *DiaryController) CreateDiary(ctx *gin.Context) {
	var request domain.DiaryReq
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := d.DiaryUsecase.FindUserByEmail(ctx, request.UserEmail)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	request.UserID = user.UserID.Int64
	_, err = d.DiaryUsecase.CreateDiary(ctx, &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	diaryResponse := domain.DiaryResponse{
		Code:    200000,
		Message: "created",
	}
	ctx.JSON(http.StatusOK, diaryResponse)
}

func (d *DiaryController) UpdateDiary(ctx *gin.Context) {
	var request domain.DiaryReq
	userID := ctx.Param("id")

	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := d.DiaryUsecase.FindUserByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	request.UserID = user.UserID.Int64
	err = d.DiaryUsecase.UpdateDiary(ctx, &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	diaryResponse := domain.DiaryResponse{
		Code:    200000,
		Message: "updated",
	}
	ctx.JSON(http.StatusOK, diaryResponse)
}

func (d *DiaryController) DeleteDiary(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid Argument"})
		return
	}
	err := d.DiaryUsecase.DeleteDiary(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "invalid Argument"})
		return
	}
	diaryResponse := domain.DiaryResponse{
		Code:    200000,
		Message: "deleted",
	}
	ctx.JSON(http.StatusOK, diaryResponse)
}
