package handlers

import (
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"

	"github.com/gin-gonic/gin"
)

type CourseCategoryHandlerImpl struct {
	ccs *services.CourseCategoryServiceImpl
}

func CreateCourseCategoryHandler(ccs *services.CourseCategoryServiceImpl) *CourseCategoryHandlerImpl {
	return &CourseCategoryHandlerImpl{ccs}
}

func (cch *CourseCategoryHandlerImpl) GetCategoriesList(ctx *gin.Context) {
	var req dtos.ListCourseCategoryReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.ListCourseCategoryParam{
		SeekPaginatedParam: entity.SeekPaginatedParam{
			Limit:  req.Limit,
			LastID: req.LastID,
		},
		Search: req.Search,
	}
	if req.Limit <= 0 {
		param.Limit = constants.DefaultLimit
	}
	if req.LastID <= 0 {
		param.LastID = constants.DefaultLastID
	}
	res, totalRow, err := cch.ccs.GetCategoriesList(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.ListCourseCategoryRes{}
	for _, cat := range *res {
		entries = append(entries, dtos.ListCourseCategoryRes(cat))
	}
	var filters []dtos.FilterInfo
	if req.Search != nil {
		filter := dtos.FilterInfo{
			Name:  "Search",
			Value: *req.Search,
		}
		filters = append(filters, filter)
	}
	var lastID int
	if len(entries) > 0 {
		lastID = entries[len(entries)-1].ID
	} else {
		lastID = 0
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.SeekPaginatedResponse[dtos.ListCourseCategoryRes]{
			Entries: entries,
			PageInfo: dtos.SeekPaginatedInfo{
				LastID:   lastID,
				Limit:    param.Limit,
				TotalRow: *totalRow,
				FilterBy: filters,
			},
		},
	})
}
