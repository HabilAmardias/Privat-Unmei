package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/utils"
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
	pr  *repositories.PaymentRepositoryImpl
	dr  *repositories.DiscountRepositoryImpl
	acr *repositories.AdditionalCostRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
	gu  *utils.GomailUtil
}

func CreateCourseRequestService(
	crr *repositories.CourseRequestRepositoryImpl,
	cr *repositories.CourseRepositoryImpl,
	csr *repositories.CourseScheduleRepositoryImpl,
	mar *repositories.MentorAvailabilityRepositoryImpl,
	ur *repositories.UserRepositoryImpl,
	sr *repositories.StudentRepositoryImpl,
	mr *repositories.MentorRepositoryImpl,
	pr *repositories.PaymentRepositoryImpl,
	dr *repositories.DiscountRepositoryImpl,
	acr *repositories.AdditionalCostRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
	gu *utils.GomailUtil,
) *CourseRequestServiceImpl {
	return &CourseRequestServiceImpl{crr, cr, csr, mar, ur, sr, mr, pr, dr, acr, tmr, gu}
}

func (crs *CourseRequestServiceImpl) StudentCourseRequestDetail(ctx context.Context, param entity.StudentCourseRequestDetailParam) (*entity.StudentCourseRequestDetailQuery, error) {
	courseRequest := new(entity.CourseRequest)
	course := new(entity.Course)
	userMentor := new(entity.User)
	mentor := new(entity.Mentor)
	userStudent := new(entity.User)
	student := new(entity.Student)
	schedules := new([]entity.CourseRequestSchedule)
	res := new(entity.StudentCourseRequestDetailQuery)
	payment := new(entity.Payment)

	if err := crs.crr.FindByID(ctx, param.CourseRequestID, courseRequest); err != nil {
		return nil, err
	}
	if err := crs.cr.FindByID(ctx, courseRequest.CourseID, course, false); err != nil {
		return nil, err
	}
	if err := crs.ur.FindByID(ctx, param.StudentID, userStudent); err != nil {
		return nil, err
	}
	if err := crs.sr.FindByID(ctx, userStudent.ID, student); err != nil {
		return nil, err
	}
	if courseRequest.StudentID != param.StudentID {
		return nil, customerrors.NewError(
			"the course does not belong to the student",
			errors.New("student id does not match"),
			customerrors.InvalidAction,
		)
	}
	if userStudent.Status != constants.VerifiedStatus {
		return nil, customerrors.NewError(
			"unauthorized user",
			errors.New("user is not verified"),
			customerrors.Unauthenticate,
		)
	}
	if err := crs.ur.FindByID(ctx, course.MentorID, userMentor); err != nil {
		return nil, err
	}
	if err := crs.mr.FindByID(ctx, userMentor.ID, mentor, false); err != nil {
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

	if err := crs.pr.FindPaymentByRequestID(ctx, courseRequest.ID, payment); err != nil {
		return nil, err
	}

	res.CourseRequestID = courseRequest.ID
	res.CourseName = course.Title
	res.MentorName = userMentor.Name
	res.MentorPublicID = userMentor.PublicID
	res.TotalPrice = payment.TotalPrice
	res.Subtotal = payment.SubTotal
	res.OperationalCost = payment.OperationalCost
	res.NumberOfSessions = courseRequest.NumberOfSessions
	res.NumberOfParticipant = courseRequest.NumberOfParticipant
	res.Status = courseRequest.Status
	res.ExpiredAt = courseRequest.ExpiredAt
	res.Schedules = *schedules
	res.MentorID = userMentor.ID
	res.AccountNumber = payment.AccountNumber
	res.PaymentMethodName = payment.PaymentMethodName

	return res, nil
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
	if err := crs.crr.StudentCourseRequestList(ctx, param.StudentID, param.Status, param.Search, param.Page, param.Limit, totalRow, requests); err != nil {
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
	payment := new(entity.Payment)

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

	if err := crs.pr.FindPaymentByRequestID(ctx, param.CourseRequestID, payment); err != nil {
		return nil, err
	}

	res.CourseRequestID = courseRequest.ID
	res.CourseName = course.Title
	res.StudentName = userStudent.Name
	res.StudentPublicID = userStudent.PublicID
	res.TotalPrice = payment.TotalPrice
	res.Subtotal = payment.SubTotal
	res.OperationalCost = payment.OperationalCost
	res.NumberOfSessions = courseRequest.NumberOfSessions
	res.PaymentMethod = payment.PaymentMethodName
	res.AccountNumber = payment.AccountNumber
	res.Status = courseRequest.Status
	res.NumberOfParticipant = courseRequest.NumberOfParticipant
	res.ExpiredAt = courseRequest.ExpiredAt
	res.Schedules = *schedules
	res.StudentID = courseRequest.StudentID

	return res, nil
}

func (crs *CourseRequestServiceImpl) MentorCourseRequestList(ctx context.Context, param entity.MentorCourseRequestListParam) (*[]entity.MentorCourseRequestQuery, *int64, error) {
	requests := new([]entity.MentorCourseRequestQuery)
	totalRow := new(int64)
	mentor := new(entity.Mentor)
	if err := crs.mr.FindByID(ctx, param.MentorID, mentor, false); err != nil {
		return nil, nil, err
	}
	if err := crs.crr.MentorCourseRequestList(ctx, param.MentorID, param.Status, param.Page, param.Limit, totalRow, requests); err != nil {
		return nil, nil, err
	}
	return requests, totalRow, nil
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
		if user.Status != constants.VerifiedStatus {
			return customerrors.NewError(
				"user is not verified",
				errors.New("user is not verified"),
				customerrors.Unauthenticate,
			)
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
	userStudent := new(entity.User)
	schedules := new([]entity.CourseRequestSchedule)
	updateCourse := new(entity.UpdateCourseQuery)
	now := time.Now()
	temp := now.Add(constants.ExpiredInterval)
	eat := &temp

	if err := crs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := crs.ur.FindByID(ctx, param.MentorID, user); err != nil {
			return err
		}
		if user.Status != constants.VerifiedStatus {
			return customerrors.NewError(
				"user is not verified",
				errors.New("user is not verified"),
				customerrors.Unauthenticate,
			)
		}
		if err := crs.mr.FindByID(ctx, param.MentorID, mentor, false); err != nil {
			return err
		}
		if err := crs.crr.FindByID(ctx, param.CourseRequestID, courseRequest); err != nil {
			return err
		}
		if err := crs.ur.FindByID(ctx, courseRequest.StudentID, userStudent); err != nil {
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
		} else {
			payment := new(entity.Payment)
			if err := crs.pr.FindPaymentByRequestID(ctx, param.CourseRequestID, payment); err != nil {
				return err
			}
			layout := "02-01-06 15:04:05"
			go func() {
				param := entity.SendEmailParams{
					Receiver:  userStudent.Email,
					Subject:   "Payment Detail",
					EmailBody: constants.PaymentInfoEmailBody(course.Title, user.Name, user.Email, payment.PaymentMethodName, payment.AccountNumber, payment.TotalPrice, courseRequest.ExpiredAt.Format(layout)),
				}
				if err := crs.gu.SendEmail(param); err != nil {
					log.Println(err.Error())
					return
				}
				log.Println("Send Email Success")
			}()
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (crs *CourseRequestServiceImpl) CreateReservation(ctx context.Context, param entity.CreateCourseRequestParam) (string, error) {
	ongoingOrderCount := new(int64)
	maxParticipant := new(int)
	operationalCost := new(float64)
	course := new(entity.Course)
	newScheds := new([]entity.CreateRequestSchedule)
	updateCourse := new(entity.UpdateCourseQuery)
	user := new(entity.User)
	student := new(entity.Student)
	userMentor := new(entity.User)
	courseRequest := new(entity.CourseRequest)
	availabilityRes := new(entity.AvailabilityResult)
	conflictingScheds := new([]entity.ConflictingSchedule)
	paymentInfo := new(entity.MentorPayment)
	discount := new(entity.Discount)
	dates := make([]time.Time, 0, len(param.PreferredSlots))
	startTimes := make([]string, 0, len(param.PreferredSlots))
	endTimes := make([]string, 0, len(param.PreferredSlots))
	now := time.Now()

	if err := crs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := crs.cr.FindByID(ctx, param.CourseID, course, true); err != nil {
			return err
		}
		if err := crs.ur.FindByID(ctx, course.MentorID, userMentor); err != nil {
			return err
		}

		if len(param.PreferredSlots) > course.MaxSession {
			return customerrors.NewError(
				fmt.Sprintf("cannot reserve more than %d slots", course.MaxSession),
				errors.New("reservation slots is more than max session"),
				customerrors.InvalidAction,
			)
		}
		if err := crs.ur.FindByID(ctx, param.StudentID, user); err != nil {
			return err
		}
		if err := crs.sr.FindByID(ctx, param.StudentID, student); err != nil {
			return err
		}
		if user.Status != constants.VerifiedStatus {
			return customerrors.NewError(
				"user are not verified yet",
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
		for _, slot := range param.PreferredSlots {
			endTime, err := CalculateEndTime(slot.StartTime.ToString(), course.SessionDuration)
			if err != nil {
				return err
			}
			dates = append(dates, slot.Date)
			startTimes = append(startTimes, slot.StartTime.ToString())
			endTimes = append(endTimes, endTime)
			*newScheds = append(*newScheds, entity.CreateRequestSchedule{
				Date:      slot.Date,
				StartTime: slot.StartTime.ToString(),
				EndTime:   endTime,
			})
		}
		if err := crs.mar.CheckMentorAvailability(ctx, course.MentorID, dates, startTimes, endTimes, availabilityRes); err != nil {
			return err
		}
		if availabilityRes.AvailableSlots < availabilityRes.TotalRequested {
			return customerrors.NewError(
				"mentor not available",
				errors.New("mentor not available"),
				customerrors.InvalidAction,
			)
		}
		if err := crs.csr.CheckScheduleConflicts(ctx, course.MentorID, dates, startTimes, endTimes, conflictingScheds); err != nil {
			return err
		}
		if len(*conflictingScheds) > 0 {
			return customerrors.NewError(
				"Some schedules are already reserved",
				errors.New("schedule already reserved"),
				customerrors.InvalidAction,
			)
		}

		if err := crs.pr.GetPaymentInfoByMentorAndMethodID(ctx, course.MentorID, param.PaymentMethodID, paymentInfo); err != nil {
			return err
		}
		if err := crs.dr.GetMaxParticipant(ctx, maxParticipant); err != nil {
			return err
		}
		participant := param.NumberOfParticipant
		if param.NumberOfParticipant > *maxParticipant {
			participant = *maxParticipant
		}
		if err := crs.dr.GetDiscountByNumberOfParticipant(ctx, participant, discount); err != nil {
			var parsedErr *customerrors.CustomError
			if errors.As(err, &parsedErr) {
				if parsedErr.ErrCode != customerrors.ItemNotExist {
					return err
				}
			}
		}
		if err := crs.acr.GetOperationalCost(ctx, operationalCost); err != nil {
			return err
		}

		pricePerSession := course.Price - discount.Amount
		if course.Price <= discount.Amount {
			pricePerSession = course.Price
		}

		subTotal := pricePerSession * float64(len(param.PreferredSlots))
		*operationalCost = *operationalCost * float64(len(param.PreferredSlots))
		courseRequest.StudentID = param.StudentID
		courseRequest.CourseID = param.CourseID
		totalPrice := subTotal + *operationalCost
		courseRequest.NumberOfParticipant = param.NumberOfParticipant
		courseRequest.NumberOfSessions = len(param.PreferredSlots)
		eat := now.Add(constants.ExpiredInterval)
		courseRequest.ExpiredAt = &eat

		if err := crs.crr.CreateOrder(ctx, courseRequest); err != nil {
			return err
		}

		if err := crs.csr.CreateSchedule(ctx, courseRequest.ID, newScheds); err != nil {
			return err
		}
		if err := crs.pr.CreatePaymentDetail(ctx, courseRequest.ID, subTotal, *operationalCost, totalPrice, paymentInfo.PaymentMethodName, paymentInfo.AccountNumber); err != nil {
			return err
		}

		updatedTransactionCount := course.TransactionCount + 1
		updateCourse.TransactionCount = &updatedTransactionCount
		if err := crs.cr.UpdateCourse(ctx, param.CourseID, updateCourse); err != nil {
			return err
		}

		go func() {
			param := entity.SendEmailParams{
				Receiver:  userMentor.Email,
				Subject:   "Request Detail",
				EmailBody: constants.RequestDetailEmailBody(course.Title, user.Name, user.Email, param.NumberOfParticipant, totalPrice),
			}
			if err := crs.gu.SendEmail(param); err != nil {
				log.Println(err.Error())
				return
			}
			log.Println("Send Email Success")
		}()

		return nil
	}); err != nil {
		return "", err
	}
	return courseRequest.ID, nil
}
