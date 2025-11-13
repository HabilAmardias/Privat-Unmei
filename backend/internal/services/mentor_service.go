package services

import (
	"context"
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/logger"
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
	crr *repositories.CourseRequestRepositoryImpl
	cr  *repositories.CourseRepositoryImpl
	pr  *repositories.PaymentRepositoryImpl
	ar  *repositories.AdminRepositoryImpl
	bu  *utils.BcryptUtil
	ju  *utils.JWTUtil
	cu  *utils.CloudinaryUtil
	gu  *utils.GomailUtil
	lg  logger.CustomLogger
}

func CreateMentorService(
	tmr *repositories.TransactionManagerRepositories,
	ur *repositories.UserRepositoryImpl,
	mr *repositories.MentorRepositoryImpl,
	tr *repositories.TopicRepositoryImpl,
	ccr *repositories.CourseCategoryRepositoryImpl,
	car *repositories.MentorAvailabilityRepositoryImpl,
	crr *repositories.CourseRequestRepositoryImpl,
	cr *repositories.CourseRepositoryImpl,
	pr *repositories.PaymentRepositoryImpl,
	ar *repositories.AdminRepositoryImpl,
	bu *utils.BcryptUtil,
	ju *utils.JWTUtil,
	cu *utils.CloudinaryUtil,
	gu *utils.GomailUtil,
	lg logger.CustomLogger,
) *MentorServiceImpl {
	return &MentorServiceImpl{tmr, ur, mr, tr, ccr, car, crr, cr, pr, ar, bu, ju, cu, gu, lg}
}

func (ms *MentorServiceImpl) GetMyPaymentMethod(ctx context.Context, param entity.GetMyPaymentMethodParam) (*[]entity.GetMentorPaymentMethodQuery, error) {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	methods := new([]entity.GetMentorPaymentMethodQuery)
	if err := ms.ur.FindByID(ctx, param.MentorID, user); err != nil {
		return nil, err
	}
	if err := ms.mr.FindByID(ctx, param.MentorID, mentor, false); err != nil {
		return nil, err
	}
	if err := ms.pr.GetMentorPaymentMethod(ctx, param.MentorID, methods); err != nil {
		return nil, err
	}
	if len(*methods) == 0 {
		return nil, customerrors.NewError(
			"mentor payment method not found",
			errors.New("mentor payment method not found"),
			customerrors.ItemNotExist,
		)
	}
	return methods, nil
}

func (ms *MentorServiceImpl) GetMentorAvailability(ctx context.Context, param entity.GetMentorAvailabilityParam) (*[]entity.MentorAvailability, error) {
	scheds := new([]entity.MentorAvailability)
	user := new(entity.User)
	mentor := new(entity.Mentor)
	if err := ms.ur.FindByID(ctx, param.MentorID, user); err != nil {
		return nil, err
	}
	if err := ms.mr.FindByID(ctx, param.MentorID, mentor, false); err != nil {
		return nil, err
	}
	if err := ms.car.GetAvailabilityByMentorID(ctx, param.MentorID, scheds); err != nil {
		return nil, err
	}
	return scheds, nil
}

func (ms *MentorServiceImpl) GetDOWAvailability(ctx context.Context, param entity.GetDOWAvailabilityParam) (*[]int, error) {
	dows := new([]int)
	user := new(entity.User)
	course := new(entity.Course)
	if err := ms.ur.FindByID(ctx, param.UserID, user); err != nil {
		return nil, err
	}
	if param.Role == constants.StudentRole {
		if user.Status != constants.VerifiedStatus {
			return nil, customerrors.NewError(
				"unverified",
				errors.New("user unverified"),
				customerrors.Unauthenticate,
			)
		}
	}
	if err := ms.cr.FindByID(ctx, param.CourseID, course, false); err != nil {
		return nil, err
	}
	if err := ms.car.GetDOWAvailability(ctx, course.MentorID, dows); err != nil {
		return nil, err
	}
	if len(*dows) == 0 {
		return nil, customerrors.NewError(
			"mentor availability not found",
			errors.New("mentor availability not found"),
			customerrors.ItemNotExist,
		)
	}
	return dows, nil
}

func (ms *MentorServiceImpl) GetMentorProfile(ctx context.Context, param entity.MentorProfileParam) (*entity.MentorProfileQuery, error) {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	res := new(entity.MentorProfileQuery)
	if err := ms.ur.FindByID(ctx, param.ID, user); err != nil {
		return nil, err
	}
	if err := ms.mr.FindByID(ctx, param.ID, mentor, false); err != nil {
		return nil, err
	}
	res.ID = user.ID
	res.Name = user.Name
	res.Email = user.Email
	res.Bio = user.Bio
	res.Campus = mentor.Campus
	res.Degree = mentor.Degree
	res.Major = mentor.Major
	res.ProfileImage = user.ProfileImage
	res.Resume = mentor.Resume
	res.YearsOfExperience = mentor.YearsOfExperience
	res.Status = user.Status

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

func (ms *MentorServiceImpl) Login(ctx context.Context, param entity.LoginMentorParam) (string, string, error) {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	authToken := new(string)
	refreshToken := new(string)

	if err := ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ms.ur.FindByEmail(ctx, param.Email, user); err != nil {
			var parsedErr *customerrors.CustomError
			if errors.As(err, &parsedErr) {
				if parsedErr.ErrCode == customerrors.ItemNotExist {
					return customerrors.NewError(
						"invalid email or password",
						errors.New("invalid email or password"),
						customerrors.InvalidAction,
					)
				}
				return err
			}
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
		loginToken, err := ms.ju.GenerateJWT(mentor.ID, constants.MentorRole, constants.ForLogin, user.Status, constants.AUTH_AGE)
		if err != nil {
			return err
		}
		rtoken, err := ms.ju.GenerateJWT(mentor.ID, constants.MentorRole, constants.ForRefresh, user.Status, constants.REFRESH_AGE)
		if err != nil {
			return err
		}
		*authToken = loginToken
		*refreshToken = rtoken
		return nil
	}); err != nil {
		return "", "", err
	}
	return *authToken, *refreshToken, nil
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
	userAdmin := new(entity.User)
	admin := new(entity.Admin)
	maxTransactionCount := new(int64)
	courseIDs := new([]int)

	return ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ms.ur.FindByID(ctx, param.AdminID, userAdmin); err != nil {
			return err
		}
		if userAdmin.Status == constants.UnverifiedStatus {
			return customerrors.NewError(
				"unauthorized",
				errors.New("admin is not verified"),
				customerrors.Unauthenticate,
			)
		}
		if err := ms.ar.FindByID(ctx, param.AdminID, admin); err != nil {
			return err
		}
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
			if err := ms.crr.DeleteAllMentorOrders(ctx, mentor.ID); err != nil {
				return err
			}
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
		if err := ms.pr.UnassignPaymentMethodFromMentor(ctx, param.ID); err != nil {
			return err
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
		if len(param.MentorPayments) > 0 {
			orderCount := new(int64)
			if err := ms.crr.FindOngoingOrderByMentorID(ctx, param.ID, orderCount); err != nil {
				return err
			}
			if *orderCount > 0 {
				return customerrors.NewError(
					"there is an ongoing order, cannot update payment method",
					errors.New("there is an ongoing order, cannot update payment method"),
					customerrors.InvalidAction,
				)
			}
			if err := ms.pr.UnassignPaymentMethodFromMentor(ctx, param.ID); err != nil {
				return err
			}
			if err := ms.pr.AssignPaymentMethodToMentor(ctx, param.ID, param.MentorPayments); err != nil {
				return err
			}
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
		query.YearsOfExperience = param.YearsOfExperience
		if err := ms.mr.UpdateMentor(ctx, mentor.ID, query); err != nil {
			return err
		}
		return nil
	})
}

func (ms *MentorServiceImpl) AddNewMentor(ctx context.Context, param entity.AddNewMentorParam) error {
	user := new(entity.User)
	userAdmin := new(entity.User)
	mentor := new(entity.Mentor)
	admin := new(entity.Admin)

	return ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ms.ur.FindByID(ctx, param.AdminID, userAdmin); err != nil {
			return err
		}
		if userAdmin.Status == constants.UnverifiedStatus {
			return customerrors.NewError(
				"unauthorized",
				errors.New("admin is unverified"),
				customerrors.Unauthenticate,
			)
		}
		if err := ms.ar.FindByID(ctx, param.AdminID, admin); err != nil {
			return err
		}
		if err := ms.ur.FindByEmail(ctx, param.Email, user); err != nil {
			var parsedErr *customerrors.CustomError
			if errors.As(err, &parsedErr) {
				if parsedErr.ErrCode != customerrors.ItemNotExist {
					return err
				}
			}
		} else {
			return customerrors.NewError(
				"user already exist",
				errors.New("user already exist"),
				customerrors.ItemAlreadyExist,
			)
		}

		user.Email = param.Email
		user.Name = param.Name
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
		mentor.YearsOfExperience = param.YearsOfExperience
		uploadRes, err := ms.cu.UploadFile(ctx, param.ResumeFile, mentor.ID, constants.ResumeFolder)
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
		if err := ms.pr.AssignPaymentMethodToMentor(ctx, mentor.ID, param.MentorPayments); err != nil {
			return err
		}
		go func() {
			emailParam := entity.SendEmailParams{
				Receiver:  param.Email,
				Subject:   "Login Credentials - Privat Unmei",
				EmailBody: constants.SendMentorAccEmailBody(param.Email, param.Password),
			}
			if err := ms.gu.SendEmail(emailParam); err != nil {
				ms.lg.Errorln(err.Error())
				newCtx := context.Background()
				if tErr := ms.tmr.WithTransaction(newCtx, func(ctx context.Context) error {
					if err := ms.pr.HardDeleteMentorPayment(ctx, mentor.ID); err != nil {
						return err
					}
					if err := ms.car.HardDeleteAvailability(ctx, mentor.ID); err != nil {
						return err
					}
					if err := ms.mr.HardDeleteMentor(ctx, mentor.ID); err != nil {
						return err
					}
					if err := ms.ur.HardDeleteUser(ctx, user.ID); err != nil {
						return err
					}
					return nil
				}); tErr != nil {
					ms.lg.Errorln(tErr.Error())
				}
			}
		}()

		return nil
	})
}
