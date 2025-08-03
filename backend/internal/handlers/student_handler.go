package handlers

import (
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"

	"github.com/gin-gonic/gin"
)

type StudentHandlerImpl struct {
	ss *services.StudentServiceImpl
}

func CreateStudentHandler(ss *services.StudentServiceImpl) *StudentHandlerImpl {
	return &StudentHandlerImpl{ss}
}

func (sh *StudentHandlerImpl) GetStudentList(ctx *gin.Context) {
	var req dtos.GetStudentListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.ListStudentParam{
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
	}
	if req.Limit <= 0 {
		param.Limit = constants.DefaultLimit
	}
	if req.Page <= 0 {
		param.Page = constants.DefaultPage
	}
	students, totalRow, err := sh.ss.GetStudentList(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.ListStudentRes{}
	for _, student := range *students {
		entries = append(entries, dtos.ListStudentRes(student))
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.PaginatedResponse[dtos.ListStudentRes]{
			Entries: entries,
			PageInfo: dtos.PaginatedInfo{
				Page:     param.Page,
				Limit:    param.Limit,
				TotalRow: *totalRow,
			},
		},
	})
}

func (sh *StudentHandlerImpl) SendVerificationEmail(ctx *gin.Context) {
	claims, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	if err := sh.ss.SendVerificationEmail(ctx, claims.Subject); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Successfully sent verification email",
		},
	})
}

func (sh *StudentHandlerImpl) ResetPassword(ctx *gin.Context) {
	var req dtos.ResetPasswordReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.ResetPasswordParam{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	}
	if err := sh.ss.ResetPassword(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Sucessfully reset password",
		},
	})
}

func (sh *StudentHandlerImpl) SendResetTokenEmail(ctx *gin.Context) {
	var req dtos.SendResetTokenEmailReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	if err := sh.ss.SendResetTokenEmail(ctx, req.Email); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Succesfully send reset password email",
		},
	})
}

func (sh *StudentHandlerImpl) Verify(ctx *gin.Context) {
	var req dtos.VerifyStudentReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	if err := sh.ss.Verify(ctx, req.Token); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Successfully Verified",
		},
	})
}

func (sh *StudentHandlerImpl) Login(ctx *gin.Context) {
	var req dtos.LoginStudentReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.StudentLoginParam{
		Email:    req.Email,
		Password: req.Password,
	}
	token, err := sh.ss.Login(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.LoginStudentRes{
			Token: token,
		},
	})
}

func (sh *StudentHandlerImpl) Register(ctx *gin.Context) {
	var req dtos.RegisterStudentReq

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.StudentRegisterParam{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Bio:      req.Bio,
		Status:   constants.UnverifiedStatus,
	}
	if err := sh.ss.Register(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Successfully Registered",
		},
	})
}
