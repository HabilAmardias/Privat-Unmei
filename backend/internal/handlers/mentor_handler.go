package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
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
}

func CreateMentorHandler(ms *services.MentorServiceImpl) *MentorHandlerImpl {
	return &MentorHandlerImpl{ms}
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
	claim, err := getAuthenticationPayload(ctx)
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

func (mh *MentorHandlerImpl) GetMentorProfileForStudent(ctx *gin.Context) {
	mentorID := ctx.Param("id")
	param := entity.GetMentorProfileForStudentParam{
		MentorID: mentorID,
	}
	detail, err := mh.ms.GetMentorProfileForStudent(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	profile := dtos.GetMentorProfileForStudentRes{
		MentorID:                detail.MentorID,
		MentorName:              detail.MentorName,
		MentorEmail:             detail.MentorEmail,
		MentorBio:               detail.MentorBio,
		MentorProfileImage:      detail.MentorProfileImage,
		MentorResume:            detail.MentorResume,
		MentorAverageRating:     detail.MentorAverageRating,
		MentorYearsOfExperience: detail.MentorYearsOfExperience,
		MentorDegree:            detail.MentorDegree,
		MentorMajor:             detail.MentorMajor,
		MentorCampus:            detail.MentorCampus,
		MentorAvailabilities:    []dtos.MentorAvailabilityRes{},
	}
	for _, sc := range detail.MentorAvailabilities {
		profile.MentorAvailabilities = append(profile.MentorAvailabilities, dtos.MentorAvailabilityRes{
			DayOfWeek: sc.DayOfWeek,
			StartTime: dtos.TimeOnly(sc.StartTime),
			EndTime:   dtos.TimeOnly(sc.EndTime),
		})
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    profile,
	})
}

func (mh *MentorHandlerImpl) GetProfileForMentor(ctx *gin.Context) {
	claims, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetProfileMentorParam{
		MentorID: claims.Subject,
	}
	res, err := mh.ms.GetProfileForMentor(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}

	profile := dtos.GetProfileMentorRes{
		ResumeFile:           res.ResumeFile,
		ProfileImage:         res.ProfileImage,
		Name:                 res.Name,
		Bio:                  res.Bio,
		YearsOfExperience:    res.YearsOfExperience,
		GopayNumber:          res.GopayNumber,
		Degree:               res.Degree,
		Major:                res.Major,
		Campus:               res.Campus,
		MentorAvailabilities: []dtos.MentorAvailabilityRes{},
	}
	for _, sched := range res.MentorAvailabilities {
		profile.MentorAvailabilities = append(profile.MentorAvailabilities, dtos.MentorAvailabilityRes{
			DayOfWeek: sched.DayOfWeek,
			StartTime: dtos.TimeOnly(sched.StartTime),
			EndTime:   dtos.TimeOnly(sched.EndTime),
		})
	}
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
	claim, err := getAuthenticationPayload(ctx)
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
	var req dtos.LoginMentorReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.LoginMentorParam{
		Email:    req.Email,
		Password: req.Password,
	}
	token, err := mh.ms.Login(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.LoginMentorRes{
			Token: token,
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
	param := entity.DeleteMentorParam{
		ID: id,
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
	if req.GopayNumber != nil {
		if err := ValidatePhoneNumber(*req.GopayNumber); err != nil {
			ctx.Error(err)
			return
		}
	}
	param := entity.UpdateMentorParam{
		ID:                id,
		GopayNumber:       req.GopayNumber,
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
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	var req dtos.UpdateMentorReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	resumeHeader, _ := ctx.FormFile("resume_file")
	profileHeader, _ := ctx.FormFile("profile_image")

	resumeFile, err := ValidateFile(resumeHeader, constants.FileSizeThreshold, constants.PDFType)
	if err != nil {
		ctx.Error(err)
		return
	}
	profileImageFile, err := ValidateFile(profileHeader, constants.FileSizeThreshold, constants.PNGType)
	if err != nil {
		ctx.Error(err)
		return
	}
	if req.GopayNumber != nil {
		if err := ValidatePhoneNumber(*req.GopayNumber); err != nil {
			ctx.Error(err)
			return
		}
	}
	if req.Degree != nil {
		if err := ValidateDegree(*req.Degree); err != nil {
			ctx.Error(err)
			return
		}
	}

	param := entity.UpdateMentorParam{
		ID:                claim.Subject,
		Resume:            resumeFile,
		ProfileImage:      profileImageFile,
		Name:              req.Name,
		Bio:               req.Bio,
		YearsOfExperience: req.YearsOfExperience,
		GopayNumber:       req.GopayNumber,
		Degree:            req.Degree,
		Major:             req.Major,
		Campus:            req.Campus,
		MentorSchedules:   []entity.MentorSchedule{},
	}
	if len(req.MentorSchedules) > 0 {
		dowsMap := make(map[int]bool)
		for _, sched := range req.MentorSchedules {
			var mentorSchedule dtos.MentorAvailabilityReq
			if err := json.Unmarshal([]byte(sched), &mentorSchedule); err != nil {
				ctx.Error(customerrors.NewError(
					"failed to parse input",
					err,
					customerrors.CommonErr,
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
	if err := ValidateDegree(req.Degree); err != nil {
		ctx.Error(err)
		return
	}
	if err := ValidatePhoneNumber(req.GopayNumber); err != nil {
		ctx.Error(err)
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
	headerFile, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(
			customerrors.NewError(
				"Failed to upload resume",
				err,
				customerrors.InvalidAction,
			),
		)
		return
	}
	file, err := ValidateFile(headerFile, constants.FileSizeThreshold, constants.PDFType)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.AddNewMentorParam{
		Name:              req.Name,
		Email:             req.Email,
		Password:          req.Password,
		ResumeFile:        file,
		Bio:               req.Bio,
		YearsOfExperience: req.YearsOfExperience,
		GopayNumber:       req.GopayNumber,
		Degree:            req.Degree,
		Major:             req.Major,
		Campus:            req.Campus,
		MentorSchedules:   []entity.MentorSchedule{},
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
