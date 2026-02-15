package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type MentorHandlerImpl struct {
	ms *services.MentorServiceImpl
	cs *services.CourseServiceImpl
}

func CreateMentorHandler(ms *services.MentorServiceImpl, cs *services.CourseServiceImpl) *MentorHandlerImpl {
	return &MentorHandlerImpl{ms, cs}
}

func (mh *MentorHandlerImpl) GetMyCourses(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	var req dtos.MentorListCourseReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.MentorListCourseParam{
		MentorID: claim.Subject,
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
		Search:         req.Search,
		CourseCategory: req.CourseCategory,
		IsProtected:    true,
	}
	if param.Limit <= 0 || param.Limit > constants.MaxLimit {
		param.Limit = constants.DefaultLimit
	}
	if param.Page <= 0 {
		param.Page = constants.DefaultPage
	}
	res, totalRow, err := mh.cs.MentorListCourse(ctx, param)
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

func (mh *MentorHandlerImpl) GetMyPaymentMethod(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetMyPaymentMethodParam{
		MentorID: claim.Subject,
	}
	methods, err := mh.ms.GetMyPaymentMethod(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.GetMentorPaymentMethodRes{}
	for _, method := range *methods {
		entries = append(entries, dtos.GetMentorPaymentMethodRes(method))
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    entries,
	})
}

func (mh *MentorHandlerImpl) GetMyAvailability(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetMentorAvailabilityParam{
		MentorID: claim.Subject,
	}
	scheds, err := mh.ms.GetMentorAvailability(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := []dtos.GetMentorAvailabilityRes{}
	for _, sch := range *scheds {
		res = append(res, dtos.GetMentorAvailabilityRes{
			DayOfWeek: sch.DayOfWeek,
			StartTime: sch.StartTime.ToString(),
			EndTime:   sch.EndTime.ToString(),
		})
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    res,
	})
}

func (mh *MentorHandlerImpl) GetMentorAvailability(ctx *gin.Context) {
	id := ctx.Param("id")
	param := entity.GetMentorAvailabilityParam{
		MentorID: id,
	}
	scheds, err := mh.ms.GetMentorAvailability(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := []dtos.GetMentorAvailabilityRes{}
	for _, sch := range *scheds {
		res = append(res, dtos.GetMentorAvailabilityRes{
			DayOfWeek: sch.DayOfWeek,
			StartTime: sch.StartTime.ToString(),
			EndTime:   sch.EndTime.ToString(),
		})
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    res,
	})
}

func (mh *MentorHandlerImpl) GetDOWAvailability(ctx *gin.Context) {
	courseID, err := strconv.Atoi(ctx.Param("id"))
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
	param := entity.GetDOWAvailabilityParam{
		Role:     claim.Role,
		CourseID: courseID,
		UserID:   claim.Subject,
	}
	dows, err := mh.ms.GetDOWAvailability(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.GetDOWAvailabilityRes{
			DayOfWeeks: *dows,
		},
	})
}

func (mh *MentorHandlerImpl) GetMyProfile(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.MentorProfileParam{
		ID: claim.Subject,
	}
	detail, err := mh.ms.GetMentorProfile(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	profile := dtos.MentorProfileRes(*detail)
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    profile,
	})
}

func (mh *MentorHandlerImpl) GetMentorProfile(ctx *gin.Context) {
	mentorID := ctx.Param("id")
	param := entity.MentorProfileParam{
		ID: mentorID,
	}
	detail, err := mh.ms.GetMentorProfile(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	profile := dtos.MentorProfileRes(*detail)
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    profile,
	})
}

func (mh *MentorHandlerImpl) ChangePassword(ctx *gin.Context) {
	var req dtos.MentorChangePasswordReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.MentorChangePasswordParam{
		ID:          claim.Subject,
		NewPassword: req.NewPassword,
	}
	if err := mh.ms.ChangePassword(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Succesfully change password",
		},
	})
}

func (mh *MentorHandlerImpl) Login(ctx *gin.Context) {
	domain := os.Getenv("COOKIE_DOMAIN")
	var req dtos.LoginMentorReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.LoginMentorParam{
		Email:    req.Email,
		Password: req.Password,
	}
	authToken, refreshToken, userStatus, err := mh.ms.Login(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	secure := os.Getenv("ENVIRONMENT_OPTION") == constants.Production
	ctx.SetCookie(constants.AUTH_COOKIE_KEY, authToken, int(constants.AUTH_AGE), "/", domain, secure, true)
	ctx.SetCookie(constants.REFRESH_COOKIE_KEY, refreshToken, int(constants.REFRESH_AGE), "/", domain, secure, true)
	ctx.SetCookie("role", strconv.Itoa(constants.MentorRole), int(constants.REFRESH_AGE), "/", domain, secure, true)
	ctx.SetCookie("status", userStatus, int(constants.REFRESH_AGE), "/", domain, secure, true)
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Successfully logged in",
		},
	})
}

func (mh *MentorHandlerImpl) GetMentorList(ctx *gin.Context) {
	var req dtos.ListMentorReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.ListMentorParam{
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
		Search:               req.Search,
		SortYearOfExperience: req.SortYearOfExperience,
	}
	if param.Limit <= 0 || param.Limit > constants.MaxLimit {
		param.Limit = constants.DefaultLimit

	}
	if param.Page <= 0 {
		param.Page = constants.DefaultPage
	}
	res, totalRow, err := mh.ms.GetMentorList(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.ListMentorRes{}
	for _, mentor := range *res {
		entries = append(entries, dtos.ListMentorRes(mentor))
	}
	var filters []dtos.FilterInfo
	if req.Search != nil {
		filter := dtos.FilterInfo{
			Name:  "Search",
			Value: *req.Search,
		}
		filters = append(filters, filter)
	}
	var sorts []dtos.SortInfo
	if req.SortYearOfExperience != nil {
		sortInfo := dtos.SortInfo{
			Name: "years_of_experience",
			ASC:  *req.SortYearOfExperience,
		}
		sorts = append(sorts, sortInfo)
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.PaginatedResponse[dtos.ListMentorRes]{
			Entries: entries,
			PageInfo: dtos.PaginatedInfo{
				Page:     param.Page,
				Limit:    param.Limit,
				TotalRow: *totalRow,
				FilterBy: filters,
				SortBy:   sorts,
			},
		},
	})
}

func (mh *MentorHandlerImpl) DeleteMentor(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(customerrors.NewError(
			"no mentor id given",
			errors.New("no mentor id given"),
			customerrors.InvalidAction,
		))
		return
	}
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.DeleteMentorParam{
		ID:      id,
		AdminID: claim.Subject,
	}
	if err := mh.ms.DeleteMentor(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.UpdateMentorForAdminRes{
			ID: id,
		},
	})
}

func (mh *MentorHandlerImpl) UpdateMentorForAdmin(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(customerrors.NewError(
			"no mentor id given",
			errors.New("no mentor id given"),
			customerrors.InvalidAction,
		))
		return
	}
	var req dtos.UpdateMentorForAdminReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.UpdateMentorParam{
		ID:                id,
		YearsOfExperience: req.YearsOfExperience,
	}
	if err := mh.ms.UpdateMentorForAdmin(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.UpdateMentorForAdminRes{
			ID: param.ID,
		},
	})
}

func (mh *MentorHandlerImpl) UpdateMentor(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	var req dtos.UpdateMentorReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	profileHeader, _ := ctx.FormFile("profile_image")
	profileImageFile, err := ValidateFile(profileHeader, constants.FileSizeThreshold, constants.PNGType)
	if err != nil {
		ctx.Error(err)
		return
	}
	if profileImageFile != nil {
		defer profileImageFile.Close()
	}
	if req.Degree != nil {
		if err := ValidateDegree(*req.Degree); err != nil {
			ctx.Error(err)
			return
		}
	}

	param := entity.UpdateMentorParam{
		ID:                claim.Subject,
		ProfileImage:      profileImageFile,
		Name:              req.Name,
		Bio:               req.Bio,
		YearsOfExperience: req.YearsOfExperience,
		MentorPayments:    []entity.AddMentorPaymentInfo{},
		Degree:            req.Degree,
		Major:             req.Major,
		Campus:            req.Campus,
		MentorSchedules:   []entity.MentorSchedule{},
	}
	if len(req.MentorPayments) > 0 {
		for _, info := range req.MentorPayments {
			var payment dtos.MentorPaymentInfoReq
			if err := json.Unmarshal([]byte(info), &payment); err != nil {
				ctx.Error(customerrors.NewError(
					"failed to parse input",
					err,
					customerrors.InvalidAction,
				))
				return
			}
			if err := binding.Validator.ValidateStruct(payment); err != nil {
				ctx.Error(err)
				return
			}
			param.MentorPayments = append(param.MentorPayments, entity.AddMentorPaymentInfo(payment))
		}
	}
	if len(req.MentorSchedules) > 0 {
		dowsMap := make(map[int]bool)
		for _, sched := range req.MentorSchedules {
			var mentorSchedule dtos.MentorAvailabilityReq
			if err := json.Unmarshal([]byte(sched), &mentorSchedule); err != nil {
				ctx.Error(customerrors.NewError(
					"failed to parse input",
					err,
					customerrors.InvalidAction,
				))
				return
			}
			if err := binding.Validator.ValidateStruct(mentorSchedule); err != nil {
				ctx.Error(err)
				return
			}
			if _, exist := dowsMap[mentorSchedule.DayOfWeek]; exist {
				ctx.Error(customerrors.NewError(
					"cannot have 2 same day of week",
					errors.New("cannot have 2 same day of week"),
					customerrors.InvalidAction,
				))
				return
			} else {
				dowsMap[mentorSchedule.DayOfWeek] = true
			}
			param.MentorSchedules = append(param.MentorSchedules, entity.MentorSchedule{
				DayOfWeek: mentorSchedule.DayOfWeek,
				StartTime: entity.TimeOnly(mentorSchedule.StartTime),
				EndTime:   entity.TimeOnly(mentorSchedule.EndTime),
			})
		}
	}
	if err := mh.ms.UpdateMentorProfile(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.UpdateMentorRes{
			ID: claim.Subject,
		},
	})
}

func (mh *MentorHandlerImpl) AddNewMentor(ctx *gin.Context) {
	var req dtos.AddNewMentorReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	if req.YearsOfExperience == nil {
		ctx.Error(customerrors.NewError(
			"years of experience should not empty",
			errors.New("years of experience should not empty"),
			customerrors.InvalidAction,
		))
		return
	}
	if err := ValidateDegree(req.Degree); err != nil {
		ctx.Error(err)
		return
	}
	if len(req.MentorPayments) <= 0 {
		ctx.Error(customerrors.NewError(
			"mentor need to submit their payment method",
			errors.New("mentor need to submit their payment method"),
			customerrors.InvalidAction,
		))
		return
	}
	if len(req.MentorSchedules) <= 0 {
		ctx.Error(customerrors.NewError(
			"mentor need to submit their schedule",
			errors.New("mentor need to submit their schedule"),
			customerrors.InvalidAction,
		))
		return
	}
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.AddNewMentorParam{
		AdminID:           claim.Subject,
		Name:              req.Name,
		Email:             req.Email,
		Password:          req.Password,
		YearsOfExperience: *req.YearsOfExperience,
		MentorPayments:    []entity.AddMentorPaymentInfo{},
		Degree:            req.Degree,
		Major:             req.Major,
		Campus:            req.Campus,
		MentorSchedules:   []entity.MentorSchedule{},
	}
	for _, info := range req.MentorPayments {
		var payment dtos.MentorPaymentInfoReq
		if err := json.Unmarshal([]byte(info), &payment); err != nil {
			ctx.Error(customerrors.NewError(
				"failed to parse input",
				err,
				customerrors.InvalidAction,
			))
			return
		}
		if err := binding.Validator.ValidateStruct(payment); err != nil {
			ctx.Error(err)
			return
		}
		param.MentorPayments = append(param.MentorPayments, entity.AddMentorPaymentInfo(payment))
	}
	dowsMap := make(map[int]bool)
	for _, sched := range req.MentorSchedules {
		var schedule dtos.MentorAvailabilityReq
		if err := json.Unmarshal([]byte(sched), &schedule); err != nil {
			ctx.Error(customerrors.NewError(
				"failed to parse input",
				err,
				customerrors.CommonErr,
			))
			return
		}
		if err := binding.Validator.ValidateStruct(schedule); err != nil {
			ctx.Error(err)
			return
		}
		if _, exist := dowsMap[schedule.DayOfWeek]; exist {
			ctx.Error(customerrors.NewError(
				"cannot have 2 same day of week",
				errors.New("cannot have 2 same day of week"),
				customerrors.InvalidAction,
			))
			return
		} else {
			dowsMap[schedule.DayOfWeek] = true
		}
		param.MentorSchedules = append(param.MentorSchedules, entity.MentorSchedule{
			DayOfWeek: schedule.DayOfWeek,
			StartTime: entity.TimeOnly(schedule.StartTime),
			EndTime:   entity.TimeOnly(schedule.EndTime),
		})
	}
	if err := mh.ms.AddNewMentor(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Succesfully create mentor account",
		},
	})
}
