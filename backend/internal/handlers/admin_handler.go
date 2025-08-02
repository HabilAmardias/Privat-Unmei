package handlers

import (
	"net/http"
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

func (ah *AdminHandlerImpl) GetStudentList(ctx *gin.Context) {
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
	if req.Limit < 0 {
		param.Limit = constants.DefaultLimit
	}
	if req.Page < 0 {
		param.Page = constants.DefaultPage
	}
	students, totalRow, err := ah.as.GetStudentList(ctx, param)
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
