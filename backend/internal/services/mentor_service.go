package services

import (
	"context"
	"errors"
	"fmt"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/utils"
)

type MentorServiceImpl struct {
	tmr *repositories.TransactionManagerRepositories
	ur  *repositories.UserRepositoryImpl
	mr  *repositories.MentorRepositoryImpl
	tr  *repositories.TopicRepositoryImpl
	ccr *repositories.CourseCategoryRepositoryImpl
	car *repositories.MentorAvailabilityRepositoryImpl
	cr  *repositories.CourseRepositoryImpl
	bu  *utils.BcryptUtil
	ju  *utils.JWTUtil
	cu  *utils.CloudinaryUtil
	gu  *utils.GomailUtil
}

func CreateMentorService(
	tmr *repositories.TransactionManagerRepositories,
	ur *repositories.UserRepositoryImpl,
	mr *repositories.MentorRepositoryImpl,
	tr *repositories.TopicRepositoryImpl,
	ccr *repositories.CourseCategoryRepositoryImpl,
	car *repositories.MentorAvailabilityRepositoryImpl,
	cr *repositories.CourseRepositoryImpl,
	bu *utils.BcryptUtil,
	ju *utils.JWTUtil,
	cu *utils.CloudinaryUtil,
	gu *utils.GomailUtil,
) *MentorServiceImpl {
	return &MentorServiceImpl{tmr, ur, mr, tr, ccr, car, cr, bu, ju, cu, gu}
}

func (ms *MentorServiceImpl) GetMentorProfileForStudent(ctx context.Context, param entity.GetMentorProfileForStudentParam) (*entity.GetMentorProfileForStudentQuery, error) {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	mentorAvailability := new([]entity.MentorAvailability)
	res := new(entity.GetMentorProfileForStudentQuery)
	if err := ms.ur.FindByID(ctx, param.MentorID, user); err != nil {
		return nil, err
	}
	if err := ms.mr.FindByID(ctx, param.MentorID, mentor, false); err != nil {
		return nil, err
	}
	if err := ms.car.GetAvailabilityByMentorID(ctx, param.MentorID, mentorAvailability); err != nil {
		return nil, err
	}

	res.MentorAvailabilities = []entity.MentorSchedule{}
	res.MentorAverageRating = constants.NoRating
	if mentor.RatingCount > constants.NoRating {
		res.MentorAverageRating = mentor.TotalRating / float64(mentor.RatingCount)
	}
	res.MentorBio = user.Bio
	res.MentorCampus = mentor.Campus
	res.MentorDegree = mentor.Degree
	res.MentorEmail = user.Email
	res.MentorID = mentor.ID
	res.MentorMajor = mentor.Major
	res.MentorName = user.Name
	res.MentorProfileImage = user.ProfileImage
	res.MentorResume = mentor.Resume
	res.MentorYearsOfExperience = mentor.YearsOfExperience

	for _, sc := range *mentorAvailability {
		res.MentorAvailabilities = append(res.MentorAvailabilities, entity.MentorSchedule{
			DayOfWeek: sc.DayOfWeek,
			StartTime: sc.StartTime,
			EndTime:   sc.EndTime,
		})
	}

	return res, nil
}

func (ms *MentorServiceImpl) GetProfileForMentor(ctx context.Context, param entity.GetProfileMentorParam) (*entity.GetProfileMentorQuery, error) {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	mentorAvailability := new([]entity.MentorAvailability)
	res := new(entity.GetProfileMentorQuery)

	if err := ms.ur.FindByID(ctx, param.MentorID, user); err != nil {
		return nil, err
	}
	if err := ms.mr.FindByID(ctx, user.ID, mentor, false); err != nil {
		return nil, err
	}
	if err := ms.car.GetAvailabilityByMentorID(ctx, mentor.ID, mentorAvailability); err != nil {
		return nil, err
	}
	if len(*mentorAvailability) <= 0 {
		return nil, customerrors.NewError(
			"mentor data does not exist",
			errors.New("mentor availability does not exist"),
			customerrors.ItemNotExist,
		)
	}
	res.ResumeFile = mentor.Resume
	res.ProfileImage = user.ProfileImage
	res.Bio = user.Bio
	res.Campus = mentor.Campus
	res.Degree = mentor.Degree
	res.GopayNumber = mentor.GopayNumber
	res.Major = mentor.Major
	res.Name = user.Name
	res.YearsOfExperience = mentor.YearsOfExperience
	res.MentorAvailabilities = []entity.MentorSchedule{}
	for _, sched := range *mentorAvailability {
		res.MentorAvailabilities = append(res.MentorAvailabilities, entity.MentorSchedule{
			DayOfWeek: sched.DayOfWeek,
			StartTime: sched.StartTime,
			EndTime:   sched.EndTime,
		})
	}
	return res, nil
}

func (ms *MentorServiceImpl) ChangePassword(ctx context.Context, param entity.MentorChangePasswordParam) error {
	user := new(entity.User)
	mentor := new(entity.Mentor)

	return ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ms.ur.FindByID(ctx, param.ID, user); err != nil {
			return err
		}
		if err := ms.mr.FindByID(ctx, param.ID, mentor, false); err != nil {
			return err
		}
		if match := ms.bu.ComparePassword(param.NewPassword, user.Password); match {
			return customerrors.NewError(
				"cannot change into same password",
				errors.New("new password same as previous password"),
				customerrors.InvalidAction,
			)
		}
		hashedPass, err := ms.bu.HashPassword(param.NewPassword)
		if err != nil {
			return err
		}
		if err := ms.ur.UpdateUserPassword(ctx, hashedPass, param.ID); err != nil {
			return err
		}
		if user.Status != constants.VerifiedStatus {
			if err := ms.ur.UpdateUserStatus(ctx, constants.VerifiedStatus, param.ID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (ms *MentorServiceImpl) Login(ctx context.Context, param entity.LoginMentorParam) (string, error) {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	token := new(string)
	if err := ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ms.ur.FindByEmail(ctx, param.Email, user); err != nil {
			parsedErr := err.(*customerrors.CustomError)
			if parsedErr.ErrUser == customerrors.UserNotFound {
				return customerrors.NewError(
					"invalid email or password",
					errors.New("invalid email or password"),
					customerrors.InvalidAction,
				)
			}
			return err
		}
		if err := ms.mr.FindByID(ctx, user.ID, mentor, false); err != nil {
			return err
		}
		if match := ms.bu.ComparePassword(param.Password, user.Password); !match {
			return customerrors.NewError(
				"invalid email or password",
				errors.New("invalid email or password"),
				customerrors.InvalidAction,
			)
		}
		loginToken, err := ms.ju.GenerateJWT(mentor.ID, constants.MentorRole, constants.ForLogin, user.Status)
		if err != nil {
			return err
		}
		*token = loginToken
		return nil
	}); err != nil {
		return "", err
	}
	return *token, nil
}

func (ms *MentorServiceImpl) GetMentorList(ctx context.Context, param entity.ListMentorParam) (*[]entity.ListMentorQuery, *int64, error) {
	mentors := new([]entity.ListMentorQuery)
	totalRow := new(int64)
	if err := ms.mr.GetMentorList(ctx, mentors, totalRow, param); err != nil {
		return nil, nil, err
	}
	return mentors, totalRow, nil
}

func (ms *MentorServiceImpl) DeleteMentor(ctx context.Context, param entity.DeleteMentorParam) error {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	maxTransactionCount := new(int64)
	courseIDs := new([]int)

	return ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ms.ur.FindByID(ctx, param.ID, user); err != nil {
			return err
		}
		if err := ms.mr.FindByID(ctx, user.ID, mentor, false); err != nil {
			return err
		}

		if err := ms.cr.GetMaximumTransactionCount(ctx, maxTransactionCount, param.ID); err != nil {
			return err
		}
		if *maxTransactionCount > 0 {
			return customerrors.NewError(
				"there are course that has been bought",
				errors.New("max transaction count is more than zero"),
				customerrors.InvalidAction,
			)
		}
		if err := ms.cr.FindByMentorID(ctx, param.ID, courseIDs); err != nil {
			return err
		}
		if len(*courseIDs) > 0 {
			if err := ms.tr.DeleteTopicsMultipleCourse(ctx, *courseIDs); err != nil {
				return err
			}
			if err := ms.ccr.UnassignCategoriesMultipleCourse(ctx, *courseIDs); err != nil {
				return err
			}
			if err := ms.cr.DeleteMentorCourse(ctx, param.ID); err != nil {
				return err
			}
		}
		if err := ms.car.DeleteAvailability(ctx, param.ID); err != nil {
			return err
		}
		if err := ms.mr.DeleteMentor(ctx, mentor.ID); err != nil {
			return err
		}
		if err := ms.ur.DeleteUser(ctx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (ms *MentorServiceImpl) UpdateMentorProfile(ctx context.Context, param entity.UpdateMentorParam) error {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	mentorQuery := new(entity.UpdateMentorQuery)
	userQuery := new(entity.UpdateUserQuery)

	return ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ms.ur.FindByID(ctx, param.ID, user); err != nil {
			return err
		}
		if err := ms.mr.FindByID(ctx, user.ID, mentor, false); err != nil {
			return err
		}
		userQuery.Name = param.Name
		userQuery.Bio = param.Bio

		if param.ProfileImage != nil {
			filename := mentor.ID
			res, err := ms.cu.UploadFile(ctx, param.ProfileImage, filename, constants.AvatarFolder)
			if err != nil {
				return err
			}
			userQuery.ProfileImage = &res.SecureURL
		}
		if len(param.MentorSchedules) > 0 {
			if err := ms.car.DeleteAvailability(ctx, param.ID); err != nil {
				return err
			}
			scheds := new([]entity.MentorAvailability)
			for _, sched := range param.MentorSchedules {
				*scheds = append(*scheds, entity.MentorAvailability{
					MentorID:  param.ID,
					DayOfWeek: sched.DayOfWeek,
					StartTime: sched.StartTime,
					EndTime:   sched.EndTime,
				})
			}
			if err := ms.car.CreateAvailability(ctx, scheds); err != nil {
				return err
			}
		}
		mentorQuery.Campus = param.Campus
		mentorQuery.Degree = param.Degree
		mentorQuery.Major = param.Major
		mentorQuery.GopayNumber = param.GopayNumber
		mentorQuery.YearsOfExperience = param.YearsOfExperience

		if param.Resume != nil {
			filename := mentor.ID
			res, err := ms.cu.UploadFile(ctx, param.Resume, filename, constants.ResumeFolder)
			if err != nil {
				return err
			}
			mentorQuery.Resume = &res.SecureURL
		}
		if err := ms.ur.UpdateUserProfile(ctx, userQuery, param.ID); err != nil {
			return err
		}
		if err := ms.mr.UpdateMentor(ctx, param.ID, mentorQuery); err != nil {
			return err
		}
		return nil
	})
}

func (ms *MentorServiceImpl) UpdateMentorForAdmin(ctx context.Context, param entity.UpdateMentorParam) error {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	query := new(entity.UpdateMentorQuery)

	return ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ms.ur.FindByID(ctx, param.ID, user); err != nil {
			return err
		}
		if err := ms.mr.FindByID(ctx, user.ID, mentor, false); err != nil {
			return err
		}
		query.GopayNumber = param.GopayNumber
		query.YearsOfExperience = param.YearsOfExperience
		if err := ms.mr.UpdateMentor(ctx, mentor.ID, query); err != nil {
			return err
		}
		return nil
	})
}

func (ms *MentorServiceImpl) AddNewMentor(ctx context.Context, param entity.AddNewMentorParam) error {
	user := new(entity.User)
	mentor := new(entity.Mentor)

	return ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		// to be honest idk how to make this clean enough but for now it should work
		if err := ms.mr.FindByGopay(ctx, param.GopayNumber, mentor); err != nil {
			if err.Error() != customerrors.UserNotFound {
				return err
			}
		} else {
			return customerrors.NewError(
				"user already exist",
				customerrors.ErrItemAlreadyExist,
				customerrors.ItemAlreadyExist,
			)
		}
		if err := ms.ur.FindByEmail(ctx, param.Email, user); err != nil {
			if err.Error() != customerrors.UserNotFound {
				return err
			}
		} else {
			return customerrors.NewError(
				"user already exist",
				customerrors.ErrItemAlreadyExist,
				customerrors.ItemAlreadyExist,
			)
		}

		user.Email = param.Email
		user.Name = param.Name
		user.Bio = param.Bio
		user.ProfileImage = constants.DefaultAvatar

		hashedPass, err := ms.bu.HashPassword(param.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPass

		// Will be changed to verified after mentor change their password for the first time
		user.Status = constants.UnverifiedStatus

		if err := ms.ur.AddNewUser(ctx, user); err != nil {
			return err
		}

		mentor.ID = user.ID
		mentor.Campus = param.Campus
		mentor.Degree = param.Degree
		mentor.Major = param.Major
		mentor.GopayNumber = param.GopayNumber
		mentor.YearsOfExperience = param.YearsOfExperience

		newFilename := fmt.Sprintf("%s.pdf", mentor.ID)
		uploadRes, err := ms.cu.UploadFile(ctx, param.ResumeFile, newFilename, constants.ResumeFolder)
		if err != nil {
			return err
		}

		mentor.Resume = uploadRes.SecureURL
		if err := ms.mr.AddNewMentor(ctx, mentor); err != nil {
			return err
		}

		schedules := new([]entity.MentorAvailability)
		for _, sched := range param.MentorSchedules {
			*schedules = append(*schedules, entity.MentorAvailability{
				MentorID:  mentor.ID,
				DayOfWeek: sched.DayOfWeek,
				StartTime: sched.StartTime,
				EndTime:   sched.EndTime,
			})
		}
		if err := ms.car.CreateAvailability(ctx, schedules); err != nil {
			return err
		}

		emailParam := entity.SendEmailParams{
			Receiver:  param.Email,
			Subject:   "Login Credentials - Privat Unmei",
			EmailBody: constants.SendMentorAccEmailBody(param.Email, param.Password),
		}
		if err := ms.gu.SendEmail(emailParam); err != nil {
			return err
		}

		return nil
	})
}
