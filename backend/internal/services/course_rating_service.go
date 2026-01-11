package services

import (
	"context"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"

	"golang.org/x/sync/errgroup"
)

type CourseRatingServiceImpl struct {
	cr  *repositories.CourseRepositoryImpl
	crr *repositories.CourseRatingRepositoryImpl
	cor *repositories.CourseRequestRepositoryImpl
	mr  *repositories.MentorRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func CreateCourseRatingService(
	cr *repositories.CourseRepositoryImpl,
	crr *repositories.CourseRatingRepositoryImpl,
	cor *repositories.CourseRequestRepositoryImpl,
	mr *repositories.MentorRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
) *CourseRatingServiceImpl {
	return &CourseRatingServiceImpl{cr, crr, cor, mr, tmr}
}

func (crs *CourseRatingServiceImpl) GetCourseReview(ctx context.Context, param entity.GetCourseRatingParam) (*[]entity.CourseRatingQuery, *int64, error) {
	g, ctx := errgroup.WithContext(ctx)
	reviews := new([]entity.CourseRatingQuery)
	totalRow := new(int64)
	course := new(entity.Course)
	g.Go(func() error {
		return crs.cr.FindByID(ctx, param.CourseID, course, false)
	})
	g.Go(func() error {
		return crs.crr.GetCourseReviews(ctx, param.CourseID, param.Page, param.Limit, totalRow, reviews)
	})
	if err := g.Wait(); err != nil {
		return nil, nil, err
	}
	return reviews, totalRow, nil
}

func (crs *CourseRatingServiceImpl) AddReview(ctx context.Context, param entity.CreateRatingParam) (int, error) {
	g, ctx := errgroup.WithContext(ctx)
	course := new(entity.Course)
	rating := new(entity.CourseRating)
	count := new(int)
	mentor := new(entity.Mentor)
	updateMentor := new(entity.UpdateMentorQuery)

	g.Go(func() error {
		return crs.cr.FindByID(ctx, param.CourseID, course, false)
	})

	g.Go(func() error {
		if err := crs.cor.FindCompletedByStudentIDAndCourseID(ctx, param.StudentID, param.CourseID, count); err != nil {
			return err
		}
		if *count == 0 {
			return customerrors.NewError(
				"user have not bought or completed the course",
				errors.New("user have not bought or completed the course"),
				customerrors.InvalidAction,
			)
		}
		return nil
	})
	g.Go(func() error {
		if err := crs.crr.FindByCourseIDAndStudentID(ctx, param.CourseID, param.StudentID, rating); err != nil {
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
				"you already reviewed this course",
				errors.New("user already reviewed the course"),
				customerrors.ItemAlreadyExist,
			)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return 0, err
	}
	if err := crs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := crs.mr.FindByID(ctx, course.MentorID, mentor, true); err != nil {
			return err
		}

		newRatingCount := mentor.RatingCount + 1
		updateMentor.RatingCount = &newRatingCount

		newTotalRating := float64(mentor.TotalRating) + float64(param.Rating)
		updateMentor.TotalRating = &newTotalRating

		if err := crs.mr.UpdateMentor(ctx, course.MentorID, updateMentor); err != nil {
			return err
		}
		if err := crs.crr.CreateReview(ctx, param.CourseID, param.StudentID, param.Rating, param.Feedback, rating); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}
	return rating.ID, nil
}

func (crs *CourseRatingServiceImpl) IsAlreadyReviewed(ctx context.Context, param entity.IsReviewedParam) (bool, error) {
	rating := new(entity.CourseRating)
	if err := crs.crr.FindByCourseIDAndStudentID(ctx, param.CourseID, param.StudentID, rating); err != nil {
		var parsedErr *customerrors.CustomError
		if !errors.As(err, &parsedErr) {
			return false, customerrors.NewError(
				"something went wrong",
				errors.New("cannot parse error"),
				customerrors.CommonErr,
			)
		}
		if parsedErr.ErrCode != customerrors.ItemNotExist {
			return false, err
		}
		return false, nil
	}
	return true, nil
}
