package services

import (
	"context"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
)

type CourseServiceImpl struct {
	car *repositories.CourseAvailabilityRepositoryImpl
	cr  *repositories.CourseRepositoryImpl
	ccr *repositories.CourseCategoryRepositoryImpl
	tr  *repositories.TopicRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
	cor *repositories.CourseRequestRepositoryImpl
}

func CreateCourseService(
	car *repositories.CourseAvailabilityRepositoryImpl,
	cr *repositories.CourseRepositoryImpl,
	ccr *repositories.CourseCategoryRepositoryImpl,
	tr *repositories.TopicRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
	cor *repositories.CourseRequestRepositoryImpl,
) *CourseServiceImpl {
	return &CourseServiceImpl{car, cr, ccr, tr, tmr, cor}
}

func (cs *CourseServiceImpl) UpdateCourse(ctx context.Context, param entity.UpdateCourseParam) error {
	updateCourseQuery := new(entity.UpdateCourseQuery)
	course := new(entity.Course)
	scheds := new([]entity.CourseAvailability)
	orders := new([]entity.CourseOrder)
	categories := new([]entity.CourseCategory)
	return cs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := cs.cr.FindByID(ctx, param.CourseID, course, true); err != nil {
			return err
		}
		if course.MentorID != param.MentorID {
			return customerrors.NewError(
				"unauthorized access",
				errors.New("different mentor"),
				customerrors.InvalidAction,
			)
		}
		if err := cs.cor.FindOngoingByCourseID(ctx, param.CourseID, orders); err != nil {
			return err
		}
		if len(*orders) > 0 {
			return customerrors.NewError(
				"there is an ongoing order for this course",
				errors.New("there is an ongoing order for this course"),
				customerrors.InvalidAction,
			)
		}
		if len(param.CourseTopic) > 0 {
			if err := cs.tr.DeleteTopics(ctx, param.CourseID); err != nil {
				return err
			}
			newTopics := new([]entity.CourseTopic)
			for _, topic := range param.CourseTopic {
				*newTopics = append(*newTopics, entity.CourseTopic{
					CourseID:    course.ID,
					Title:       topic.Title,
					Description: topic.Description,
				})
			}
			if err := cs.tr.CreateTopics(ctx, newTopics); err != nil {
				return err
			}
		}
		if len(param.CourseCategories) > 0 {
			if err := cs.ccr.UnassignCategories(ctx, param.CourseID); err != nil {
				return err
			}
			if err := cs.ccr.FindByMultipleIDs(ctx, param.CourseCategories, categories); err != nil {
				return err
			}
			if len(*categories) != len(param.CourseCategories) {
				return customerrors.NewError(
					"invalid course categories",
					errors.New("number of categories and number of ids does not match"),
					customerrors.InvalidAction,
				)
			}
			if err := cs.ccr.AssignCategories(ctx, course.ID, param.CourseCategories); err != nil {
				return err
			}
		}
		if len(param.CourseSchedule) > 0 {
			if err := cs.car.DeleteAvailability(ctx, param.CourseID); err != nil {
				return err
			}
			for _, schedule := range param.CourseSchedule {
				if schedule.EndTime.Hour < schedule.StartTime.Hour {
					return customerrors.NewError(
						"invalid schedule",
						errors.New("invalid schedule"),
						customerrors.InvalidAction,
					)
				}
				*scheds = append(*scheds, entity.CourseAvailability{
					CourseID:  course.ID,
					DayOfWeek: schedule.DayOfWeek,
					StartTime: schedule.StartTime,
					EndTime:   schedule.EndTime,
				})
			}
			if err := cs.car.CreateAvailability(ctx, scheds); err != nil {
				return err
			}
		}

		updateCourseQuery.Title = param.Title
		updateCourseQuery.Description = param.Description
		updateCourseQuery.Domicile = param.Description
		updateCourseQuery.MaxDurationDays = param.MaxDurationDays
		updateCourseQuery.MaxPrice = param.MaxPrice
		updateCourseQuery.Method = param.Method
		updateCourseQuery.MinDurationDays = param.MinDurationDays
		updateCourseQuery.MinPrice = param.MinPrice

		if err := cs.cr.UpdateCourse(ctx, param.CourseID, updateCourseQuery); err != nil {
			return err
		}
		return nil
	})
}

func (cs *CourseServiceImpl) CourseDetail(ctx context.Context, param entity.CourseDetailParam) (*entity.CourseDetailQuery, error) {
	query := new(entity.CourseDetailQuery)
	topics := new([]entity.CourseTopic)
	scheds := new([]entity.CourseAvailability)
	if err := cs.cr.CourseDetail(ctx, query, param.ID); err != nil {
		return nil, err
	}
	if err := cs.tr.GetTopicsByCourseID(ctx, param.ID, topics); err != nil {
		return nil, err
	}
	if len(*topics) == 0 {
		return nil, customerrors.NewError(
			"no course topic found",
			errors.New("the course does not have a topic"),
			customerrors.ItemNotExist,
		)
	}
	query.Topics = topics
	if err := cs.car.GetAvailabilityByCourseID(ctx, param.ID, scheds); err != nil {
		return nil, err
	}
	if len(*scheds) == 0 {
		return nil, customerrors.NewError(
			"schedule availability for this course does not exist",
			errors.New("schedule availability for this course does not exist"),
			customerrors.ItemNotExist,
		)
	}
	query.Schedules = scheds
	return query, nil
}

func (cs *CourseServiceImpl) ListCourse(ctx context.Context, param entity.ListCourseParam) (*[]entity.CourseListQuery, *int64, error) {
	query := new([]entity.CourseListQuery)
	totalRow := new(int64)
	if err := cs.cr.ListCourse(ctx, query, totalRow, param); err != nil {
		return nil, nil, err
	}
	return query, totalRow, nil
}

func (cs *CourseServiceImpl) MostBoughtCourses(ctx context.Context) (*[]entity.CourseListQuery, error) {
	query := new([]entity.CourseListQuery)
	if err := cs.cr.GetMostBoughtCourses(ctx, query); err != nil {
		return nil, err
	}
	return query, nil
}

func (cs *CourseServiceImpl) MentorListCourse(ctx context.Context, param entity.MentorListCourseParam) (*[]entity.MentorListCourseQuery, *int64, error) {
	query := new([]entity.MentorListCourseQuery)
	totalRow := new(int64)
	if err := cs.cr.MentorListCourse(ctx, query, totalRow, param); err != nil {
		return nil, nil, err
	}
	return query, totalRow, nil
}

func (cs *CourseServiceImpl) DeleteCourse(ctx context.Context, param entity.DeleteCourseParam) error {
	course := new(entity.Course)
	return cs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := cs.cr.FindByID(ctx, param.CourseID, course, false); err != nil {
			return err
		}
		if param.MentorID != course.MentorID {
			return customerrors.NewError(
				"not authorized to delete",
				errors.New("course mentor and mentor input is different"),
				customerrors.Unauthenticate,
			)
		}
		if course.TransactionCount > 0 {
			return customerrors.NewError(
				"course can only be deleted if no transaction relate to the course",
				errors.New("transaction count is more than zero"),
				customerrors.InvalidAction,
			)
		}
		if err := cs.ccr.UnassignCategories(ctx, course.ID); err != nil {
			return err
		}
		if err := cs.car.DeleteAvailability(ctx, param.CourseID); err != nil {
			return err
		}
		if err := cs.tr.DeleteTopics(ctx, param.CourseID); err != nil {
			return err
		}
		if err := cs.cr.DeleteCourse(ctx, param.CourseID); err != nil {
			return err
		}
		return nil
	})
}

func (cs *CourseServiceImpl) CreateCourse(ctx context.Context, param entity.CreateCourseParam) (int, error) {
	course := new(entity.Course)
	topics := new([]entity.CourseTopic)
	schedules := new([]entity.CourseAvailability)
	categories := new([]entity.CourseCategory)

	err := cs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if param.MaxDuration < param.MinDuration {
			return customerrors.NewError(
				"invalid course duration",
				errors.New("invalid course duration"),
				customerrors.InvalidAction,
			)
		}
		if param.MaxPrice < param.MinPrice {
			return customerrors.NewError(
				"invalid course price",
				errors.New("invalid course price"),
				customerrors.InvalidAction,
			)
		}
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
		if err := cs.ccr.FindByMultipleIDs(ctx, param.Categories, categories); err != nil {
			return err
		}
		if len(*categories) != len(param.Categories) {
			return customerrors.NewError(
				"invalid course categories",
				errors.New("number of categories and number of ids does not match"),
				customerrors.InvalidAction,
			)
		}
		if err := cs.ccr.AssignCategories(ctx, course.ID, param.Categories); err != nil {
			return err
		}
		if err := cs.tr.CreateTopics(ctx, topics); err != nil {
			return err
		}
		for _, schedule := range param.CourseAvailability {
			if schedule.EndTime.Hour < schedule.StartTime.Hour {
				return customerrors.NewError(
					"invalid schedule",
					errors.New("invalid schedule"),
					customerrors.InvalidAction,
				)
			}
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
