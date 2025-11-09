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

type CourseRatingHandlerImpl struct {
	crs *services.CourseRatingServiceImpl
}

func CreateCourseRatingHandler(crs *services.CourseRatingServiceImpl) *CourseRatingHandlerImpl {
	return &CourseRatingHandlerImpl{crs}
}

func (crh *CourseRatingHandlerImpl) GetCourseReview(ctx *gin.Context) {
	courseIDStr := ctx.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid course credential",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	var req dtos.CourseRatingReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetCourseRatingParam{
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
		CourseID: courseID,
	}
	if param.Limit <= 0 || param.Limit > constants.MaxLimit {
		param.Limit = constants.DefaultLimit
	}
	if param.Page <= 0 {
		param.Page = constants.DefaultPage
	}
	res, totalRow, err := crh.crs.GetCourseReview(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.CourseRatingRes{}
	for _, course := range *res {
		entries = append(entries, dtos.CourseRatingRes(course))
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.PaginatedResponse[dtos.CourseRatingRes]{
			Entries: entries,
			PageInfo: dtos.PaginatedInfo{
				Page:     param.Page,
				Limit:    param.Limit,
				TotalRow: *totalRow,
			},
		},
	})
}

func (crh *CourseRatingHandlerImpl) AddReview(ctx *gin.Context) {
	courseIDStr := ctx.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid course credential",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	var req dtos.CreateRatingReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.CreateRatingParam{
		StudentID: claim.Subject,
		CourseID:  courseID,
		Rating:    req.Rating,
		Feedback:  req.Feedback,
	}
	reviewID, err := crh.crs.AddReview(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.CreateRatingRes{
			RatingID: reviewID,
		},
	})
}
