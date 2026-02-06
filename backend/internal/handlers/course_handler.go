package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseHandlerImpl struct {
	cs *services.CourseServiceImpl
}

func CreateCourseHandler(cs *services.CourseServiceImpl) *CourseHandlerImpl {
	return &CourseHandlerImpl{cs}
}

func (ch *CourseHandlerImpl) UpdateCourse(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid course",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	var req dtos.UpdateCourseReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	if req.Method != nil {
		if err := ValidateMethod(*req.Method); err != nil {
			ctx.Error(err)
			return
		}
	}
	if len(req.CourseCategories) > 0 {
		if err := ValidateCategories(req.CourseCategories); err != nil {
			ctx.Error(err)
			return
		}
	}
	if len(req.CourseCategories) > constants.MaxCourseCategories {
		ctx.Error(customerrors.NewError(
			fmt.Sprintf("max num of course categories is %d", constants.MaxCourseCategories),
			fmt.Errorf("max num of course categories is %d", constants.MaxCourseCategories),
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.UpdateCourseParam{
		MentorID: claim.Subject,
		CourseID: id,
		UpdateCourseQuery: entity.UpdateCourseQuery{
			Title:           req.Title,
			Description:     req.Description,
			Domicile:        req.Domicile,
			Price:           req.Price,
			Method:          req.Method,
			SessionDuration: req.SessionDuration,
			MaxSession:      req.MaxSession,
		},
		CourseTopic:      []entity.CreateTopic{},
		CourseCategories: req.CourseCategories,
	}
	if len(req.CourseTopic) > 0 {
		for _, t := range req.CourseTopic {
			param.CourseTopic = append(param.CourseTopic, entity.CreateTopic(t))
		}
	}
	if err := ch.cs.UpdateCourse(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.UpdateCourseRes{
			ID: id,
		},
	})
}

func (ch *CourseHandlerImpl) CourseCategories(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid course",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.CourseDetailParam{
		ID: id,
	}
	topics, err := ch.cs.GetCourseCategory(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := []dtos.GetCategoriesRes{}
	for _, c := range *topics {
		res = append(res, dtos.GetCategoriesRes{
			CategoryID:   c.CategoryID,
			CategoryName: c.CategoryName,
		})
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    res,
	})
}

func (ch *CourseHandlerImpl) CourseTopics(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid course",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.CourseDetailParam{
		ID: id,
	}
	topics, err := ch.cs.GetCourseTopic(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := []dtos.CourseTopicRes{}
	for _, c := range *topics {
		res = append(res, dtos.CourseTopicRes{
			Title:       c.Title,
			Description: c.Description,
		})
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    res,
	})
}

func (ch *CourseHandlerImpl) CourseDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid course",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.CourseDetailParam{
		ID: id,
	}
	res, err := ch.cs.CourseDetail(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entry := dtos.CourseDetailRes{
		CourseListRes: dtos.CourseListRes{
			MentorListCourseRes: dtos.MentorListCourseRes{
				ID:              res.ID,
				Title:           res.Title,
				Domicile:        res.Domicile,
				Method:          res.Method,
				Price:           res.Price,
				SessionDuration: res.SessionDuration,
				MaxSession:      res.MaxSession,
			},
			MentorID:           res.MentorID,
			MentorName:         res.MentorName,
			MentorPublicID:     res.MentorPublicID,
			MentorProfileImage: res.MentorProfileImage,
		},
		Description: res.Description,
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    entry,
	})
}

func (ch *CourseHandlerImpl) ListCourse(ctx *gin.Context) {
	var req dtos.ListCourseReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.ListCourseParam{
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
		Search:         req.Search,
		CourseCategory: req.CourseCategory,
		Method:         req.Method,
	}
	if param.Limit <= 0 || param.Limit > constants.MaxLimit {
		param.Limit = constants.DefaultLimit
	}
	if param.Page <= 0 {
		param.Page = constants.DefaultPage
	}
	res, totalRow, err := ch.cs.ListCourse(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.CourseListRes{}
	for _, course := range *res {
		item := dtos.CourseListRes{
			MentorListCourseRes: dtos.MentorListCourseRes{
				ID:              course.ID,
				Title:           course.Title,
				Domicile:        course.Domicile,
				Method:          course.Method,
				Price:           course.Price,
				SessionDuration: course.SessionDuration,
				MaxSession:      course.MaxSession,
			},
			MentorID:       course.MentorID,
			MentorName:     course.MentorName,
			MentorPublicID: course.MentorPublicID,
		}
		entries = append(entries, item)
	}
	var filters []dtos.FilterInfo
	if req.Search != nil {
		filter := dtos.FilterInfo{
			Name:  "Search",
			Value: *req.Search,
		}
		filters = append(filters, filter)
	}
	if req.CourseCategory != nil {
		filter := dtos.FilterInfo{
			Name:  "Course Category",
			Value: *req.CourseCategory,
		}
		filters = append(filters, filter)
	}
	if req.Method != nil {
		filter := dtos.FilterInfo{
			Name:  "Method",
			Value: *req.Method,
		}
		filters = append(filters, filter)
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.PaginatedResponse[dtos.CourseListRes]{
			Entries: entries,
			PageInfo: dtos.PaginatedInfo{
				Page:     param.Page,
				FilterBy: filters,
				Limit:    param.Limit,
				TotalRow: *totalRow,
			},
		},
	})
}

func (ch *CourseHandlerImpl) MostBoughtCourses(ctx *gin.Context) {
	res, err := ch.cs.MostBoughtCourses(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.CourseListRes{}
	for _, course := range *res {
		item := dtos.CourseListRes{
			MentorListCourseRes: dtos.MentorListCourseRes{
				ID:              course.ID,
				Title:           course.Title,
				Domicile:        course.Domicile,
				Method:          course.Method,
				Price:           course.Price,
				SessionDuration: course.SessionDuration,
				MaxSession:      course.MaxSession,
			},
			MentorID:       course.MentorID,
			MentorName:     course.MentorName,
			MentorPublicID: course.MentorPublicID,
		}
		entries = append(entries, item)
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    entries,
	})
}

func (ch *CourseHandlerImpl) MentorListCourse(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dtos.MentorListCourseReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.MentorListCourseParam{
		MentorID: id,
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
		Search:         req.Search,
		CourseCategory: req.CourseCategory,
	}
	if param.Limit <= 0 || param.Limit > constants.MaxLimit {
		param.Limit = constants.DefaultLimit
	}
	if param.Page <= 0 {
		param.Page = constants.DefaultPage
	}
	res, totalRow, err := ch.cs.MentorListCourse(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.MentorListCourseRes{}
	for _, item := range *res {
		entry := dtos.MentorListCourseRes{
			ID:              item.ID,
			Title:           item.Title,
			Domicile:        item.Domicile,
			Method:          item.Method,
			Price:           item.Price,
			SessionDuration: item.SessionDuration,
			MaxSession:      item.MaxSession,
		}
		entries = append(entries, entry)
	}
	var filters []dtos.FilterInfo
	if req.Search != nil {
		filter := dtos.FilterInfo{
			Name:  "Search",
			Value: *req.Search,
		}
		filters = append(filters, filter)
	}
	if req.CourseCategory != nil {
		filter := dtos.FilterInfo{
			Name:  "Course Category",
			Value: *req.CourseCategory,
		}
		filters = append(filters, filter)
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.PaginatedResponse[dtos.MentorListCourseRes]{
			Entries: entries,
			PageInfo: dtos.PaginatedInfo{
				Page:     param.Page,
				FilterBy: filters,
				Limit:    param.Limit,
				TotalRow: *totalRow,
			},
		},
	})
}

func (ch *CourseHandlerImpl) DeleteCourse(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	courseIDStr := ctx.Param("id")
	if courseIDStr == "" {
		ctx.Error(customerrors.NewError(
			"no course has been selected",
			errors.New("course id is empty"),
			customerrors.InvalidAction,
		))
		return
	}
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid course",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.DeleteCourseParam{
		MentorID: claim.Subject,
		CourseID: courseID,
	}

	if err := ch.cs.DeleteCourse(ctx, param); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.DeleteCourseRes{
			ID: courseID,
		},
	})
}

func (ch *CourseHandlerImpl) AddNewCourse(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	var req dtos.CreateCourseReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	if err := ValidateMethod(req.Method); err != nil {
		ctx.Error(err)
		return
	}
	if len(req.Topics) <= 0 {
		ctx.Error(customerrors.NewError(
			"mentor should enter course topic at least one",
			errors.New("mentor should enter course topic at least one"),
			customerrors.InvalidAction,
		))
		return
	}
	if len(req.Categories) <= 0 {
		ctx.Error(
			customerrors.NewError(
				"mentor should at least enter one category for the course",
				errors.New("no categories entered"),
				customerrors.InvalidAction,
			),
		)
		return
	}
	if len(req.Categories) > constants.MaxCourseCategories {
		ctx.Error(customerrors.NewError(
			fmt.Sprintf("max num of course categories is %d", constants.MaxCourseCategories),
			fmt.Errorf("max num of course categories is %d", constants.MaxCourseCategories),
			customerrors.InvalidAction,
		))
		return
	}
	if err := ValidateCategories(req.Categories); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.CreateCourseParam{
		MentorID:        claim.Subject,
		Title:           req.Title,
		Description:     req.Description,
		Domicile:        req.Domicile,
		Price:           req.Price,
		Method:          req.Method,
		SessionDuration: req.SessionDuration,
		MaxSession:      req.MaxSession,
		Topics:          []entity.CreateTopic{},
		Categories:      req.Categories,
	}
	for _, topic := range req.Topics {
		param.Topics = append(param.Topics, entity.CreateTopic(topic))
	}
	courseID, err := ch.cs.CreateCourse(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.CreateCourseRes{
			ID: courseID,
		},
	})
}
