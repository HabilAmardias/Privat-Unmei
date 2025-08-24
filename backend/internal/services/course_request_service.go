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

func (crs *CourseRequestServiceImpl) HandleCourseRequest(ctx context.Context, param entity.HandleCourseRequestParam) error {
	course := new(entity.Course)
	courseRequest := new(entity.CourseRequest)
	user := new(entity.User)
	mentor := new(entity.Mentor)
	schedules := new([]entity.CourseRequestSchedule)
	now := time.Now()
	temp := now.Add(2 * 24 * time.Hour)
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

		totalPrice := course.Price * float64(len(param.PreferredSlots)) * constants.OperationalCostPercentage

		courseRequest.StudentID = param.StudentID
		courseRequest.CourseID = param.CourseID
		courseRequest.TotalPrice = totalPrice
		courseRequest.NumberOfSessions = len(param.PreferredSlots)
		eat := now.Add(2 * 24 * time.Hour)
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
