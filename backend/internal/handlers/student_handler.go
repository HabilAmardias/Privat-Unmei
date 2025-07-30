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
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Successfully Registered",
		},
	})
}
