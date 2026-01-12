package services

import (
	"context"
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"

	"golang.org/x/sync/errgroup"
)

type CourseServiceImpl struct {
	cr  *repositories.CourseRepositoryImpl
	ccr *repositories.CourseCategoryRepositoryImpl
	tr  *repositories.TopicRepositoryImpl
	ur  *repositories.UserRepositoryImpl
	mr  *repositories.MentorRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
	cor *repositories.CourseRequestRepositoryImpl
}

func CreateCourseService(
	cr *repositories.CourseRepositoryImpl,
	ccr *repositories.CourseCategoryRepositoryImpl,
	tr *repositories.TopicRepositoryImpl,
	ur *repositories.UserRepositoryImpl,
	mr *repositories.MentorRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
	cor *repositories.CourseRequestRepositoryImpl,
) *CourseServiceImpl {
	return &CourseServiceImpl{cr, ccr, tr, ur, mr, tmr, cor}
}

func (cs *CourseServiceImpl) UpdateCourse(ctx context.Context, param entity.UpdateCourseParam) error {
	updateCourseQuery := new(entity.UpdateCourseQuery)
	course := new(entity.Course)
	count := new(int)
	categories := new([]entity.CourseCategory)
	mentor := new(entity.Mentor)
	user := new(entity.User)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := cs.ur.FindByID(ctx, param.MentorID, user); err != nil {
			return err
		}
		if user.Status != constants.VerifiedStatus {
			return customerrors.NewError(
				"please verify your account",
				errors.New("unverified account"),
				customerrors.Unauthenticate,
			)
		}
		return nil
	})
	g.Go(func() error {
		return cs.mr.FindByID(ctx, param.MentorID, mentor, false)
	})
	g.Go(func() error {
		if err := cs.cor.FindOngoingByCourseID(ctx, param.CourseID, count); err != nil {
			return err
		}
		if *count > 0 {
			return customerrors.NewError(
				"there is an ongoing order for this course",
				errors.New("there is an ongoing order for this course"),
				customerrors.InvalidAction,
			)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
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

		updateCourseQuery.Title = param.Title
		updateCourseQuery.Description = param.Description
		updateCourseQuery.Domicile = param.Domicile
		updateCourseQuery.SessionDuration = param.SessionDuration
		updateCourseQuery.Price = param.Price
		updateCourseQuery.Method = param.Method
		updateCourseQuery.MaxSession = param.MaxSession

		if err := cs.cr.UpdateCourse(ctx, param.CourseID, updateCourseQuery); err != nil {
			return err
		}
		return nil
	})
}

func (cs *CourseServiceImpl) GetCourseTopic(ctx context.Context, param entity.CourseDetailParam) (*[]entity.CourseTopic, error) {
	query := new([]entity.CourseTopic)
	if err := cs.tr.GetTopicsByCourseID(ctx, param.ID, query); err != nil {
		return nil, err
	}
	if len(*query) == 0 {
		return nil, customerrors.NewError(
			"no topic found",
			errors.New("no topic found"),
			customerrors.ItemNotExist,
		)
	}
	return query, nil
}

func (cs *CourseServiceImpl) GetCourseCategory(ctx context.Context, param entity.CourseDetailParam) (*[]entity.GetCategoriesQuery, error) {
	query := new([]entity.GetCategoriesQuery)
	if err := cs.ccr.GetCategoriesByCourseID(ctx, param.ID, query); err != nil {
		return nil, err
	}
	if len(*query) == 0 {
		return nil, customerrors.NewError(
			"no categories found",
			errors.New("no categories found"),
			customerrors.ItemNotExist,
		)
	}
	return query, nil
}

func (cs *CourseServiceImpl) CourseDetail(ctx context.Context, param entity.CourseDetailParam) (*entity.CourseDetailQuery, error) {
	query := new(entity.CourseDetailQuery)
	if err := cs.cr.CourseDetail(ctx, query, param.ID); err != nil {
		return nil, err
	}

	return query, nil
}

func (cs *CourseServiceImpl) ListCourse(ctx context.Context, param entity.ListCourseParam) (*[]entity.CourseListQuery, *int64, error) {
	query := new([]entity.CourseListQuery)
	totalRow := new(int64)
	if err := cs.cr.ListCourse(ctx, query, totalRow, param.Limit, param.Page, param.Search, param.Method, param.CourseCategory); err != nil {
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
	mentor := new(entity.Mentor)
	user := new(entity.User)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := cs.ur.FindByID(ctx, param.MentorID, user); err != nil {
			return err
		}
		if param.IsProtected && user.Status != constants.VerifiedStatus {
			return customerrors.NewError(
				"please verify your account",
				errors.New("unverified account"),
				customerrors.Unauthenticate,
			)
		}
		return nil
	})
	g.Go(func() error {
		return cs.mr.FindByID(ctx, param.MentorID, mentor, false)
	})
	g.Go(func() error {
		return cs.cr.MentorListCourse(ctx, query, totalRow, param.Limit, param.Page, param.MentorID, param.Search, param.CourseCategory)
	})
	if err := g.Wait(); err != nil {
		return nil, nil, err
	}
	return query, totalRow, nil
}

func (cs *CourseServiceImpl) DeleteCourse(ctx context.Context, param entity.DeleteCourseParam) error {
	course := new(entity.Course)
	mentor := new(entity.Mentor)
	user := new(entity.User)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := cs.ur.FindByID(ctx, param.MentorID, user); err != nil {
			return err
		}
		if user.Status != constants.VerifiedStatus {
			return customerrors.NewError(
				"please verify your account",
				errors.New("unverified account"),
				customerrors.Unauthenticate,
			)
		}
		return nil
	})
	g.Go(func() error {
		return cs.mr.FindByID(ctx, param.MentorID, mentor, false)
	})
	g.Go(func() error {
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
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return cs.tmr.WithTransaction(ctx, func(ctx context.Context) error {

		if err := cs.ccr.UnassignCategories(ctx, course.ID); err != nil {
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
	categories := new([]entity.CourseCategory)
	mentor := new(entity.Mentor)

	err := cs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := cs.mr.FindByID(ctx, param.MentorID, mentor, false); err != nil {
			return err
		}
		if err := cs.cr.CreateCourse(
			ctx,
			param.MentorID,
			param.Title,
			param.Description,
			param.Domicile,
			param.Price,
			param.SessionDuration,
			param.MaxSession,
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
		return nil
	})
	if err != nil {
		return 0, err
	}
	return course.ID, nil
}
