package handlers

import (
	"errors"
	"io"
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"

	"github.com/gin-gonic/gin"
)

type StudentHandlerImpl struct {
	ss *services.StudentServiceImpl
}

func (sh *StudentHandlerImpl) Register(ctx *gin.Context) {
	var req dtos.RegisterStudentReq

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(customerrors.NewError(
			errors.New("failed to register"),
			err,
			customerrors.InvalidAction,
		))
		return
	}

	fileReq, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(customerrors.NewError(
			errors.New("invalid image profile"),
			err,
			customerrors.InvalidAction,
		))
		return
	}

	file, err := fileReq.Open()
	if err != nil {
		ctx.Error(
			customerrors.NewError(
				errors.New("failed to upload file"),
				err,
				customerrors.InvalidAction,
			),
		)
		return
	}
	defer file.Close()

	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		ctx.Error(
			customerrors.NewError(
				errors.New("failed to upload file"),
				err,
				customerrors.InvalidAction,
			),
		)
		return
	}

	file.Seek(0, io.SeekStart)
	fileType := http.DetectContentType(buff)
	if fileType != constants.PNGType && fileType != constants.JPGType {
		ctx.Error(
			customerrors.NewError(
				errors.New("uploaded file must be .png or .jpg"),
				errors.New("uploaded file must be .png or .jpg"),
				customerrors.InvalidAction,
			),
		)
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
