package handlers

import (
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseCategoryHandlerImpl struct {
	ccs *services.CourseCategoryServiceImpl
}

func CreateCourseCategoryHandler(ccs *services.CourseCategoryServiceImpl) *CourseCategoryHandlerImpl {
	return &CourseCategoryHandlerImpl{ccs}
}

func (cch *CourseCategoryHandlerImpl) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid category",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.DeleteCategoryParam{
		AdminID: claim.Subject,
		ID:      id,
	}
	if err := cch.ccs.DeleteCategory(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.CategoryIDRes{
			ID: id,
		},
	})
}

func (cch *CourseCategoryHandlerImpl) UpdateCategory(ctx *gin.Context) {
	var req dtos.UpdateCategoryReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid category",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.UpdateCategoryParam{
		ID:      id,
		Name:    req.Name,
		AdminID: claim.Subject,
	}
	if err := cch.ccs.UpdateCategory(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.CategoryIDRes{
			ID: id,
		},
	})
}

func (cch *CourseCategoryHandlerImpl) CreateCategory(ctx *gin.Context) {
	var req dtos.CreateCategoryReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.CreateCategoryParam{
		Name:    req.Name,
		AdminID: claim.Subject,
	}
	newCategory, err := cch.ccs.CreateCategory(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.CategoryIDRes{
			ID: newCategory.ID,
		},
	})
}

func (cch *CourseCategoryHandlerImpl) GetCategoriesList(ctx *gin.Context) {
	var req dtos.ListCourseCategoryReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.ListCourseCategoryParam{
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
		Search: req.Search,
	}
	if param.Limit <= 0 || param.Limit > constants.MaxLimit {
		param.Limit = constants.DefaultLimit
	}
	if param.Page <= 0 {
		param.Page = constants.DefaultPage
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
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.PaginatedResponse[dtos.ListCourseCategoryRes]{
			Entries: entries,
			PageInfo: dtos.PaginatedInfo{
				Page:     param.Page,
				Limit:    param.Limit,
				TotalRow: *totalRow,
				FilterBy: filters,
			},
		},
	})
}
