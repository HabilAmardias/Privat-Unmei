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

// TODO: Add change email and password admin for first login

func (ah *AdminHandlerImpl) VerifyAdmin(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	var req dtos.AdminVerificationReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.AdminVerificationParam{
		AdminID:  claim.Subject,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := ah.as.VerifyAdmin(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.AdminIDRes{
			ID: param.AdminID,
		},
	})
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
	token, status, err := ah.as.Login(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.AdminLoginRes{
			Token:  *token,
			Status: *status,
		},
	})
}
