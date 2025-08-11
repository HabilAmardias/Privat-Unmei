package handlers

import (
	"net/http"
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
