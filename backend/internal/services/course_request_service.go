package services

import (
	"context"
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
	"time"
)

type CourseRequestServiceImpl struct {
	crr *repositories.CourseRequestRepositoryImpl
	cr  *repositories.CourseRepositoryImpl
	csr *repositories.CourseScheduleRepositoryImpl
	mar *repositories.MentorAvailabilityRepositoryImpl
	ur  *repositories.UserRepositoryImpl
	sr  *repositories.StudentRepositoryImpl
	mr  *repositories.MentorRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func CreateCourseRequestService(
	crr *repositories.CourseRequestRepositoryImpl,
	cr *repositories.CourseRepositoryImpl,
	csr *repositories.CourseScheduleRepositoryImpl,
	mar *repositories.MentorAvailabilityRepositoryImpl,
	ur *repositories.UserRepositoryImpl,
	sr *repositories.StudentRepositoryImpl,
	mr *repositories.MentorRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
) *CourseRequestServiceImpl {
	return &CourseRequestServiceImpl{crr, cr, csr, mar, ur, sr, mr, tmr}
}

func (crs *CourseRequestServiceImpl) StudentCourseRequestList(ctx context.Context, param entity.StudentCourseRequestListParam) (*[]entity.StudentCourseRequestQuery, *int64, error) {
	requests := new([]entity.StudentCourseRequestQuery)
	totalRow := new(int64)
	user := new(entity.User)
	student := new(entity.Student)
	if err := crs.ur.FindByID(ctx, param.StudentID, user); err != nil {
		return nil, nil, err
	}
	if err := crs.sr.FindByID(ctx, user.ID, student); err != nil {
		return nil, nil, err
	}
	if user.Status != constants.VerifiedStatus {
		return nil, nil, customerrors.NewError(
			"unauthorized access",
			errors.New("user status is not verified"),
			customerrors.Unauthenticate,
		)
	}
	if err := crs.crr.StudentCourseRequestList(ctx, param.StudentID, param.Status, param.Search, param.LastID, param.Limit, totalRow, requests); err != nil {
		return nil, nil, err
	}
	return requests, totalRow, nil
}

func (crs *CourseRequestServiceImpl) MentorCourseRequestDetail(ctx context.Context, param entity.MentorCourseRequestDetailParam) (*entity.MentorCourseRequestDetailQuery, error) {
	courseRequest := new(entity.CourseRequest)
	course := new(entity.Course)
	userMentor := new(entity.User)
	mentor := new(entity.Mentor)
	userStudent := new(entity.User)
	student := new(entity.Student)
	schedules := new([]entity.CourseRequestSchedule)
	res := new(entity.MentorCourseRequestDetailQuery)

	if err := crs.crr.FindByID(ctx, param.CourseRequestID, courseRequest); err != nil {
		return nil, err
	}
	if err := crs.cr.FindByID(ctx, courseRequest.CourseID, course, false); err != nil {
		return nil, err
	}
	if err := crs.ur.FindByID(ctx, param.MentorID, userMentor); err != nil {
		return nil, err
	}
	if err := crs.mr.FindByID(ctx, userMentor.ID, mentor, false); err != nil {
		return nil, err
	}
	if course.MentorID != param.MentorID {
		return nil, customerrors.NewError(
			"the course does not belong to the mentor",
			errors.New("mentor id does not match"),
			customerrors.InvalidAction,
		)
	}
	if err := crs.ur.FindByID(ctx, courseRequest.StudentID, userStudent); err != nil {
		return nil, err
	}
	if err := crs.sr.FindByID(ctx, userStudent.ID, student); err != nil {
		return nil, err
	}
	if err := crs.csr.FindScheduleByCourseRequestID(ctx, param.CourseRequestID, schedules); err != nil {
		return nil, err
	}

	if len(*schedules) != courseRequest.NumberOfSessions {
		return nil, customerrors.NewError(
			"something went wrong",
			errors.New("course schedule and number of session does not match, integrity breached"),
			customerrors.CommonErr,
		)
	}

	res.CourseRequestID = courseRequest.ID
	res.CourseName = course.Title
	res.StudentName = userStudent.Name
	res.StudentEmail = userStudent.Email
	res.TotalPrice = courseRequest.TotalPrice
	res.Subtotal = courseRequest.SubTotal
	res.OperationalCost = courseRequest.OperationalCost
	res.NumberOfSessions = courseRequest.NumberOfSessions
	res.Status = courseRequest.Status
	res.ExpiredAt = courseRequest.ExpiredAt
	res.Schedules = *schedules

	return res, nil
}

func (crs *CourseRequestServiceImpl) MentorCourseRequestList(ctx context.Context, param entity.MentorCourseRequestListParam) (*[]entity.MentorCourseRequestQuery, *int64, error) {
	requests := new([]entity.MentorCourseRequestQuery)
	totalRow := new(int64)
	if err := crs.crr.MentorCourseRequestList(ctx, param.MentorID, param.Status, param.LastID, param.Limit, totalRow, requests); err != nil {
		return nil, nil, err
	}
	return requests, totalRow, nil
}

func (crs *CourseRequestServiceImpl) GetPaymentDetail(ctx context.Context, param entity.GetPaymentDetailParam) (*entity.PaymentDetailQuery, error) {
	courseRequest := new(entity.CourseRequest)
	course := new(entity.Course)
	userMentor := new(entity.User)
	mentor := new(entity.Mentor)
	userStudent := new(entity.User)
	student := new(entity.Student)
	query := new(entity.PaymentDetailQuery)
	now := time.Now()

	if err := crs.ur.FindByID(ctx, param.UserID, userStudent); err != nil {
		return nil, err
	}
	if err := crs.sr.FindByID(ctx, userStudent.ID, student); err != nil {
		return nil, err
	}
	if err := crs.crr.GetPaymentDetail(ctx, param.CourseRequestID, courseRequest); err != nil {
		return nil, err
	}
	if err := crs.cr.FindByID(ctx, courseRequest.CourseID, course, false); err != nil {
		return nil, err
	}
	if err := crs.ur.FindByID(ctx, course.MentorID, userMentor); err != nil {
		return nil, err
	}
	if err := crs.mr.FindByID(ctx, userMentor.ID, mentor, false); err != nil {
		return nil, err
	}
	if courseRequest.Status != constants.PendingPaymentStatus {
		return nil, customerrors.NewError(
			"invalid course request",
			errors.New("course request status is not pending payment"),
			customerrors.InvalidAction,
		)
	}
	if courseRequest.ExpiredAt == nil {
		return nil, customerrors.NewError(
			"invalid course request",
			errors.New("no expire date found"),
			customerrors.CommonErr,
		)
	}
	if now.After(*courseRequest.ExpiredAt) {
		return nil, customerrors.NewError(
			"course request has expired",
			errors.New("course request has expired"),
			customerrors.InvalidAction,
		)
	}
	if courseRequest.StudentID != student.ID {
		return nil, customerrors.NewError(
			"unauthorized",
			errors.New("course request student id is different"),
			customerrors.Unauthenticate,
		)
	}
	if userStudent.Status != constants.VerifiedStatus {
		return nil, customerrors.NewError(
			"need to verify account",
			errors.New("user is still unverified"),
			customerrors.Unauthenticate,
		)
	}

	query.CourseID = course.ID
	query.CourseRequestID = param.CourseRequestID
	query.CourseTitle = course.Title
	query.ExpiredAt = *courseRequest.ExpiredAt
	query.GopayNumber = mentor.GopayNumber
	query.MentorID = mentor.ID
	query.MentorName = userMentor.Name
	query.OperationalCost = courseRequest.OperationalCost
	query.Subtotal = courseRequest.SubTotal
	query.TotalCost = courseRequest.TotalPrice

	return query, nil
}

func (crs *CourseRequestServiceImpl) MentorConfirmPayment(ctx context.Context, param entity.ConfirmPaymentParam) error {
	course := new(entity.Course)
	courseRequest := new(entity.CourseRequest)
	user := new(entity.User)
	mentor := new(entity.Mentor)
	schedules := new([]entity.CourseRequestSchedule)
	now := time.Now()

	if err := crs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := crs.ur.FindByID(ctx, param.MentorID, user); err != nil {
			return err
		}
		if err := crs.mr.FindByID(ctx, param.MentorID, mentor, false); err != nil {
			return err
		}
		if err := crs.crr.FindByID(ctx, param.CourseRequestID, courseRequest); err != nil {
			return err
		}
		if err := crs.cr.FindByID(ctx, courseRequest.CourseID, course, false); err != nil {
			return err
		}
		if err := crs.csr.FindReservedScheduleByCourseRequestID(ctx, param.CourseRequestID, schedules); err != nil {
			return err
		}
		if course.MentorID != param.MentorID {
			return customerrors.NewError(
				"unauthorized access",
				errors.New("unauthorized access"),
				customerrors.Unauthenticate,
			)
		}
		if courseRequest.Status != constants.PendingPaymentStatus {
			return customerrors.NewError(
				"invalid course request",
				errors.New("course request is not on reserved status"),
				customerrors.InvalidAction,
			)
		}
		if courseRequest.ExpiredAt == nil {
			return customerrors.NewError(
				"invalid course request",
				errors.New("no expiration date"),
				customerrors.InvalidAction,
			)
		}
		if now.After(*courseRequest.ExpiredAt) {
			return customerrors.NewError(
				"course request already expired",
				errors.New("course request already expired"),
				customerrors.InvalidAction,
			)
		}
		if len(*schedules) != courseRequest.NumberOfSessions {
			return customerrors.NewError(
				"requested number of session and number of schedules does not match",
				errors.New("requested number of session and number of schedules does not match"),
				customerrors.InvalidAction,
			)
		}
		if err := crs.crr.ChangeRequestStatus(ctx, param.CourseRequestID, constants.ScheduledStatus, nil); err != nil {
			return err
		}
		if err := crs.csr.UpdateScheduleStatusByCourseRequestID(ctx, param.CourseRequestID, constants.ScheduledStatus); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (crs *CourseRequestServiceImpl) HandleCourseRequest(ctx context.Context, param entity.HandleCourseRequestParam) error {
	course := new(entity.Course)
	courseRequest := new(entity.CourseRequest)
	user := new(entity.User)
	mentor := new(entity.Mentor)
	schedules := new([]entity.CourseRequestSchedule)
	updateCourse := new(entity.UpdateCourseQuery)
	now := time.Now()
	temp := now.Add(constants.ExpiredInterval)
	eat := &temp

	if err := crs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := crs.ur.FindByID(ctx, param.MentorID, user); err != nil {
			return err
		}
		if err := crs.mr.FindByID(ctx, param.MentorID, mentor, false); err != nil {
			return err
		}
		if err := crs.crr.FindByID(ctx, param.CourseRequestID, courseRequest); err != nil {
			return err
		}
		if err := crs.cr.FindByID(ctx, courseRequest.CourseID, course, true); err != nil {
			return err
		}
		if err := crs.csr.FindReservedScheduleByCourseRequestID(ctx, param.CourseRequestID, schedules); err != nil {
			return err
		}
		if course.MentorID != param.MentorID {
			return customerrors.NewError(
				"unauthorized access",
				errors.New("unauthorized access"),
				customerrors.Unauthenticate,
			)
		}
		if courseRequest.Status != constants.ReservedStatus {
			return customerrors.NewError(
				"invalid course request",
				errors.New("course request is not on reserved status"),
				customerrors.InvalidAction,
			)
		}
		if courseRequest.ExpiredAt == nil {
			return customerrors.NewError(
				"invalid course request",
				errors.New("no expiration date"),
				customerrors.InvalidAction,
			)
		}
		if now.After(*courseRequest.ExpiredAt) {
			return customerrors.NewError(
				"course request already expired",
				errors.New("course request already expired"),
				customerrors.InvalidAction,
			)
		}
		if len(*schedules) != courseRequest.NumberOfSessions {
			return customerrors.NewError(
				"requested number of session and number of schedules does not match",
				errors.New("requested number of session and number of schedules does not match"),
				customerrors.InvalidAction,
			)
		}
		status := constants.PendingPaymentStatus
		if !param.Accept {
			status = constants.CancelledStatus
			eat = nil
		}
		if err := crs.crr.ChangeRequestStatus(ctx, param.CourseRequestID, status, eat); err != nil {
			return err
		}
		if !param.Accept {
			if err := crs.csr.UpdateScheduleStatusByCourseRequestID(ctx, param.CourseRequestID, constants.CancelledStatus); err != nil {
				return err
			}
			updatedTransactionCount := course.TransactionCount - 1
			updateCourse.TransactionCount = &updatedTransactionCount
			if err := crs.cr.UpdateCourse(ctx, course.ID, updateCourse); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (crs *CourseRequestServiceImpl) CreateReservation(ctx context.Context, param entity.CreateCourseRequestParam) (int, error) {
	freeMentorSchedule := new(int64)
	existingSchedule := new(int64)
	ongoingOrderCount := new(int64)
	course := new(entity.Course)
	newScheds := new([]entity.CreateRequestSchedule)
	updateCourse := new(entity.UpdateCourseQuery)
	user := new(entity.User)
	student := new(entity.Student)
	courseRequest := new(entity.CourseRequest)
	now := time.Now()

	if err := crs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := crs.cr.FindByID(ctx, param.CourseID, course, true); err != nil {
			return err
		}
		if err := crs.ur.FindByID(ctx, param.StudentID, user); err != nil {
			return err
		}
		if err := crs.sr.FindByID(ctx, param.StudentID, student); err != nil {
			return err
		}
		if user.Status != constants.VerifiedStatus {
			return customerrors.NewError(
				"you are not verified yet",
				errors.New("user are not verified yet"),
				customerrors.Unauthenticate,
			)
		}
		if err := crs.crr.FindOngoingByCourseIDAndStudentID(ctx, param.CourseID, param.StudentID, ongoingOrderCount); err != nil {
			return err
		}
		if *ongoingOrderCount > 0 {
			return customerrors.NewError(
				"There is an ongoing order for this course",
				errors.New("there is an ongoing order for this course"),
				customerrors.InvalidAction,
			)
		}
		// I hate this solution to the bone because of N+1 Query
		// but i can't find another solution for this case, my skill issue I guess
		for _, slot := range param.PreferredSlots {
			endTime, err := CalculateEndTime(slot.StartTime.ToString(), course.SessionDuration)
			if err != nil {
				return err
			}
			if err := crs.mar.IsMentorAvailable(ctx, course.MentorID, slot.StartTime.ToString(), endTime, freeMentorSchedule, slot.Date); err != nil {
				return err
			}
			if *freeMentorSchedule == 0 {
				return customerrors.NewError(
					"Mentor not available",
					errors.New("mentor not available"),
					customerrors.InvalidAction,
				)
			}
			if err := crs.csr.IsScheduleExist(ctx, course.MentorID, slot.Date, slot.StartTime.ToString(), endTime, existingSchedule); err != nil {
				return err
			}
			if *existingSchedule > 0 {
				return customerrors.NewError(
					"Schedule already reserved",
					errors.New("schedule already reserved"),
					customerrors.InvalidAction,
				)
			}
			*newScheds = append(*newScheds, entity.CreateRequestSchedule{
				Date:      slot.Date,
				StartTime: slot.StartTime.ToString(),
				EndTime:   endTime,
			})
		}
		courseRequest.SubTotal = course.Price * float64(len(param.PreferredSlots))
		courseRequest.OperationalCost = courseRequest.SubTotal * constants.OperationalCostPercentage
		totalPrice := courseRequest.SubTotal + courseRequest.OperationalCost

		courseRequest.StudentID = param.StudentID
		courseRequest.CourseID = param.CourseID
		courseRequest.TotalPrice = totalPrice
		courseRequest.NumberOfSessions = len(param.PreferredSlots)
		eat := now.Add(constants.ExpiredInterval)
		courseRequest.ExpiredAt = &eat

		if err := crs.crr.CreateOrder(ctx, courseRequest); err != nil {
			return err
		}
		if err := crs.csr.CreateSchedule(ctx, courseRequest.ID, newScheds); err != nil {
			return err
		}

		updatedTransactionCount := course.TransactionCount + 1
		updateCourse.TransactionCount = &updatedTransactionCount
		if err := crs.cr.UpdateCourse(ctx, param.CourseID, updateCourse); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return 0, err
	}
	return courseRequest.ID, nil
}
