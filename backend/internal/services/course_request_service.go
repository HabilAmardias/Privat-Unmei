package services

import (
	"context"
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
)

type CourseRequestServiceImpl struct {
	crr *repositories.CourseRequestRepositoryImpl
	cr  *repositories.CourseRepositoryImpl
	csr *repositories.CourseScheduleRepositoryImpl
	mar *repositories.MentorAvailabilityRepositoryImpl
	ur  *repositories.UserRepositoryImpl
	sr  *repositories.StudentRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func CreateCourseRequestService(
	crr *repositories.CourseRequestRepositoryImpl,
	cr *repositories.CourseRepositoryImpl,
	csr *repositories.CourseScheduleRepositoryImpl,
	mar *repositories.MentorAvailabilityRepositoryImpl,
	ur *repositories.UserRepositoryImpl,
	sr *repositories.StudentRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
) *CourseRequestServiceImpl {
	return &CourseRequestServiceImpl{crr, cr, csr, mar, ur, sr, tmr}
}

func (crs *CourseRequestServiceImpl) CreateReservation(ctx context.Context, param entity.CreateOrderParam) (*int64, error) {
	course := new(entity.Course)
	freeMentorSchedule := new(int64)
	existingSchedule := new(int64)
	newScheds := new([]entity.CreateRequestSchedule)
	courseRequestID := new(int64)
	updateCourse := new(entity.UpdateCourseQuery)
	user := new(entity.User)
	student := new(entity.Student)

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
		if err := crs.crr.CreateOrder(ctx, param.StudentID, param.CourseID, totalPrice, len(param.PreferredSlots), courseRequestID); err != nil {
			return err
		}
		if err := crs.csr.CreateSchedule(ctx, *courseRequestID, newScheds); err != nil {
			return err
		}
		updatedTransactionCount := course.TransactionCount + 1
		updateCourse.TransactionCount = &updatedTransactionCount
		if err := crs.cr.UpdateCourse(ctx, param.CourseID, updateCourse); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return courseRequestID, nil
}
