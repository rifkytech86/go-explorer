package cms_controllers

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/naonweh-studio/bubbme-backend/domain"
	"net/http"
)

type AuthorizationController struct {
	AuthorizationUseCase domain.IAuthorizationUsecase
}

func (a *AuthorizationController) Login(ctx *gin.Context) {

	var request domain.AuthRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	resp, err := a.AuthorizationUseCase.Login(ctx, request.Email, request.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)

}

func (a *AuthorizationController) Logout(ctx *gin.Context) {

	var request domain.AuthRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	resp, err := a.AuthorizationUseCase.LogOut(ctx, request.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)

}
