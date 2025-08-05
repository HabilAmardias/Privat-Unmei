package services

import (
	"context"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
)

type CourseServiceImpl struct {
	car *repositories.CourseAvailabilityRepositoryImpl
	cr  *repositories.CourseRepositoryImpl
	tr  *repositories.TopicRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func CreateCourseService(
	car *repositories.CourseAvailabilityRepositoryImpl,
	cr *repositories.CourseRepositoryImpl,
	tr *repositories.TopicRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
) *CourseServiceImpl {
	return &CourseServiceImpl{car, cr, tr, tmr}
}

func (cs *CourseServiceImpl) CreateCourse(ctx context.Context, param entity.CreateCourseParam) (int, error) {
	course := new(entity.Course)
	topics := new([]entity.CourseTopic)
	schedules := new([]entity.CourseAvailability)

	err := cs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := cs.cr.CreateCourse(
			ctx,
			param.MentorID,
			param.Title,
			param.Description,
			param.Domicile,
			param.MinPrice,
			param.MaxPrice,
			param.MinDuration,
			param.MaxDuration,
			param.Method,
			course,
		); err != nil {
			return err
		}
		for _, topic := range param.Topics {
			*topics = append(*topics, entity.CourseTopic{
				CourseID:    course.ID,
				Title:       topic.Title,
				Description: topic.Description,
			})
		}
		if err := cs.tr.CreateTopics(ctx, topics); err != nil {
			return err
		}
		for _, schedule := range param.CourseAvailability {
			*schedules = append(*schedules, entity.CourseAvailability{
				CourseID:  course.ID,
				DayOfWeek: schedule.DayOfWeek,
				StartTime: schedule.StartTime,
				EndTime:   schedule.EndTime,
			})
		}
		if err := cs.car.CreateAvailability(ctx, schedules); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return course.ID, nil
}
