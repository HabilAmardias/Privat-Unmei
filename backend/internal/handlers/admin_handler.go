package handlers

import (
	"net/http"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminHandlerImpl struct {
	as *services.AdminServiceImpl
}

func CreateAdminHandler(as *services.AdminServiceImpl) *AdminHandlerImpl {
	return &AdminHandlerImpl{as}
}

func (ah *AdminHandlerImpl) GenerateRandomPassword(ctx *gin.Context) {
	pass, err := ah.as.GenerateRandomPassword()
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.GeneratePasswordRes{
			Password: pass,
		},
	})
}

func (ah *AdminHandlerImpl) Login(ctx *gin.Context) {
	var req dtos.AdminLoginReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.AdminLoginParam(req)
	token, err := ah.as.Login(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.AdminLoginRes{
			Token: token,
		},
	})
}
