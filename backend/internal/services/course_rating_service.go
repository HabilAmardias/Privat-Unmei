package services

import (
	"context"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
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
	reviews := new([]entity.CourseRatingQuery)
	totalRow := new(int64)
	course := new(entity.Course)
	if err := crs.cr.FindByID(ctx, param.CourseID, course, false); err != nil {
		return nil, nil, err
	}
	if err := crs.crr.GetCourseReviews(ctx, param.CourseID, param.LastID, param.Limit, totalRow, reviews); err != nil {
		return nil, nil, err
	}
	return reviews, totalRow, nil
}

func (crs *CourseRatingServiceImpl) AddReview(ctx context.Context, param entity.CreateRatingParam) (int, error) {
	course := new(entity.Course)
	rating := new(entity.CourseRating)
	orders := new([]entity.CourseRequest)
	mentor := new(entity.Mentor)
	updateMentor := new(entity.UpdateMentorQuery)

	if err := crs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := crs.cr.FindByID(ctx, param.CourseID, course, false); err != nil {
			return err
		}
		if err := crs.cor.FindCompletedByStudentIDAndCourseID(ctx, param.StudentID, param.CourseID, orders); err != nil {
			return err
		}
		if len(*orders) == 0 {
			return customerrors.NewError(
				"user have not bought or completed the course",
				errors.New("user have not bought or completed the course"),
				customerrors.InvalidAction,
			)
		}
		if err := crs.crr.FindByCourseIDAndStudentID(ctx, param.CourseID, param.StudentID, rating); err != nil {
			parsedErr := err.(*customerrors.CustomError)
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
		if err := crs.mr.FindByID(ctx, course.MentorID, mentor, true); err != nil {
			return err
		}
		*updateMentor.RatingCount = mentor.RatingCount + 1
		*updateMentor.TotalRating = float64(mentor.TotalRating) + float64(param.Rating)
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
