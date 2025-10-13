package handlers

import (
	"net/http"
	"os"
	"privat-unmei/internal/constants"
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

func (ah *AdminHandlerImpl) ChangePassword(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	var req dtos.AdminUpdatePasswordReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.AdminUpdatePasswordParam{
		AdminID:  claim.Subject,
		Password: req.Password,
	}
	if err := ah.as.UpdatePassword(ctx, param); err != nil {
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
	domain := os.Getenv("COOKIE_DOMAIN")
	var req dtos.AdminLoginReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.AdminLoginParam(req)
	authToken, refreshToken, status, err := ah.as.Login(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.SetCookie(constants.AUTH_COOKIE_KEY, *authToken, int(constants.AUTH_AGE), "/", domain, false, true)
	ctx.SetCookie(constants.REFRESH_COOKIE_KEY, *refreshToken, int(constants.REFRESH_AGE), "/", domain, false, true)
	ctx.SetCookie("status", *status, int(constants.REFRESH_AGE), "/", domain, false, true)
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.AdminLoginRes{
			Status: *status,
		},
	})
}
