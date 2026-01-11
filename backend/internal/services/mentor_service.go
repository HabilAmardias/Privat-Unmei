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

	"golang.org/x/sync/errgroup"
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
	chr *repositories.ChatRepositoryImpl
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
	chr *repositories.ChatRepositoryImpl,
	bu *utils.BcryptUtil,
	ju *utils.JWTUtil,
	cu *utils.CloudinaryUtil,
	gu *utils.GomailUtil,
	lg logger.CustomLogger,
) *MentorServiceImpl {
	return &MentorServiceImpl{tmr, ur, mr, tr, ccr, car, crr, cr, pr, ar, chr, bu, ju, cu, gu, lg}
}

func (ms *MentorServiceImpl) GetMyPaymentMethod(ctx context.Context, param entity.GetMyPaymentMethodParam) (*[]entity.GetMentorPaymentMethodQuery, error) {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	methods := new([]entity.GetMentorPaymentMethodQuery)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return ms.ur.FindByID(ctx, param.MentorID, user)
	})
	g.Go(func() error {
		return ms.mr.FindByID(ctx, param.MentorID, mentor, false)
	})
	g.Go(func() error {
		return ms.pr.GetMentorPaymentMethod(ctx, param.MentorID, methods)
	})
	if err := g.Wait(); err != nil {
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

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return ms.ur.FindByID(ctx, param.MentorID, user)
	})
	g.Go(func() error {
		return ms.mr.FindByID(ctx, param.MentorID, mentor, false)
	})
	g.Go(func() error {
		return ms.car.GetAvailabilityByMentorID(ctx, param.MentorID, scheds)
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return scheds, nil
}

func (ms *MentorServiceImpl) GetDOWAvailability(ctx context.Context, param entity.GetDOWAvailabilityParam) (*[]int, error) {
	dows := new([]int)
	user := new(entity.User)
	course := new(entity.Course)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := ms.ur.FindByID(ctx, param.UserID, user); err != nil {
			return err
		}
		if user.Status != constants.VerifiedStatus {
			return customerrors.NewError(
				"unverified",
				errors.New("user unverified"),
				customerrors.Unauthenticate,
			)
		}
		return nil
	})
	g.Go(func() error {
		if err := ms.cr.FindByID(ctx, param.CourseID, course, false); err != nil {
			return err
		}
		if err := ms.car.GetDOWAvailability(ctx, course.MentorID, dows); err != nil {
			return err
		}
		return nil
	})
	if err := g.Wait(); err != nil {
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

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return ms.ur.FindByID(ctx, param.ID, user)
	})
	g.Go(func() error {
		return ms.mr.FindByID(ctx, param.ID, mentor, false)
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	res.ID = user.ID
	res.Name = user.Name
	res.PublicID = user.PublicID
	res.Bio = user.Bio
	res.Campus = mentor.Campus
	res.Degree = mentor.Degree
	res.Major = mentor.Major
	res.ProfileImage = user.ProfileImage
	res.YearsOfExperience = mentor.YearsOfExperience
	if mentor.RatingCount > 0 {
		res.AverageRating = float64(mentor.TotalRating) / float64(mentor.RatingCount)
	}

	return res, nil
}

func (ms *MentorServiceImpl) ChangePassword(ctx context.Context, param entity.MentorChangePasswordParam) error {
	user := new(entity.User)
	mentor := new(entity.Mentor)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return ms.ur.FindByID(ctx, param.ID, user)
	})
	g.Go(func() error {
		return ms.mr.FindByID(ctx, param.ID, mentor, false)
	})
	if err := g.Wait(); err != nil {
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
	return ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
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

func (ms *MentorServiceImpl) Login(ctx context.Context, param entity.LoginMentorParam) (string, string, string, error) {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	authToken := new(string)
	refreshToken := new(string)

	if err := ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ms.ur.FindByEmail(ctx, param.Email, user); err != nil {
			var parsedErr *customerrors.CustomError
			if !errors.As(err, &parsedErr) {
				return customerrors.NewError(
					"something went wrong",
					errors.New("cannot parse error"),
					customerrors.CommonErr,
				)
			}
			if parsedErr.ErrCode == customerrors.ItemNotExist {
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
		return "", "", "", err
	}
	return *authToken, *refreshToken, user.Status, nil
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

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
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
		return nil
	})
	g.Go(func() error {
		return ms.ar.FindByID(ctx, param.AdminID, admin)
	})
	g.Go(func() error {
		return ms.ur.FindByID(ctx, param.ID, user)
	})
	g.Go(func() error {
		return ms.mr.FindByID(ctx, param.ID, mentor, false)
	})
	g.Go(func() error {
		return ms.cr.GetMaximumTransactionCount(ctx, maxTransactionCount, param.ID)
	})
	g.Go(func() error {
		return ms.cr.FindByMentorID(ctx, param.ID, courseIDs)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if *maxTransactionCount > 0 {
			if err := ms.crr.DeleteAllMentorOrders(ctx, mentor.ID); err != nil {
				return err
			}
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
		if err := ms.chr.DeleteUserMessages(ctx, param.ID); err != nil {
			return err
		}
		if err := ms.chr.DeleteMentorChatrooms(ctx, param.ID); err != nil {
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

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return ms.ur.FindByID(ctx, param.ID, user)
	})
	g.Go(func() error {
		return ms.mr.FindByID(ctx, param.ID, mentor, false)
	})
	userQuery.Name = param.Name
	userQuery.Bio = param.Bio

	if err := g.Wait(); err != nil {
		return err
	}

	if param.ProfileImage != nil {
		filename := param.ID
		res, err := ms.cu.UploadFile(context.Background(), param.ProfileImage, filename, constants.AvatarFolder)
		if err != nil {
			return err
		}
		userQuery.ProfileImage = &res.SecureURL
	}

	return ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
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

		if err := ms.ur.UpdateUserProfile(ctx, userQuery, param.ID); err != nil {
			return err
		}
		if len(param.MentorPayments) > 0 {
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

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return ms.ur.FindByID(ctx, param.ID, user)
	})
	g.Go(func() error {
		return ms.mr.FindByID(ctx, param.ID, mentor, false)
	})
	if err := g.Wait(); err != nil {
		return err
	}
	query.YearsOfExperience = param.YearsOfExperience
	if err := ms.mr.UpdateMentor(ctx, mentor.ID, query); err != nil {
		return err
	}
	return nil
}

func (ms *MentorServiceImpl) AddNewMentor(ctx context.Context, param entity.AddNewMentorParam) error {
	user := new(entity.User)
	userAdmin := new(entity.User)
	mentor := new(entity.Mentor)
	admin := new(entity.Admin)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
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
		return nil
	})
	g.Go(func() error {
		return ms.ar.FindByID(ctx, param.AdminID, admin)
	})
	g.Go(func() error {
		if err := ms.ur.FindByEmail(ctx, param.Email, user); err != nil {
			var parsedErr *customerrors.CustomError
			if !errors.As(err, &parsedErr) {
				return customerrors.NewError(
					"something went wrong",
					errors.New("cannot parse error"),
					customerrors.CommonErr,
				)
			}
			if parsedErr.ErrCode != customerrors.ItemNotExist {
				return err
			}
		} else {
			return customerrors.NewError(
				"user already exist",
				errors.New("user already exist"),
				customerrors.ItemAlreadyExist,
			)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
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

	mentor.Campus = param.Campus
	mentor.Degree = param.Degree
	mentor.Major = param.Major
	mentor.YearsOfExperience = param.YearsOfExperience

	if err := ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ms.ur.AddNewUser(ctx, user); err != nil {
			return err
		}
		mentor.ID = user.ID
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

		return nil
	}); err != nil {
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
			if tErr := ms.tmr.WithTransaction(context.Background(), func(ctx context.Context) error {
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
}
