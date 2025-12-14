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

type CourseCategoryServiceImpl struct {
	ur  *repositories.UserRepositoryImpl
	ar  *repositories.AdminRepositoryImpl
	ccr *repositories.CourseCategoryRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func CreateCourseCategoryService(ur *repositories.UserRepositoryImpl, ar *repositories.AdminRepositoryImpl, ccr *repositories.CourseCategoryRepositoryImpl, tmr *repositories.TransactionManagerRepositories) *CourseCategoryServiceImpl {
	return &CourseCategoryServiceImpl{ur, ar, ccr, tmr}
}

func (ccs *CourseCategoryServiceImpl) DeleteCategory(ctx context.Context, param entity.DeleteCategoryParam) error {
	g, ctx := errgroup.WithContext(ctx)
	admin := new(entity.Admin)
	user := new(entity.User)
	category := new(entity.CourseCategory)
	count := new(int)
	g.Go(func() error {
		if err := ccs.ur.FindByID(ctx, param.AdminID, user); err != nil {
			return err
		}
		if user.Status == constants.UnverifiedStatus {
			return customerrors.NewError(
				"unauthorized",
				errors.New("admin is not verified"),
				customerrors.Unauthenticate,
			)
		}
		return nil
	})
	g.Go(func() error {
		return ccs.ar.FindByID(ctx, param.AdminID, admin)
	})
	g.Go(func() error {
		return ccs.ccr.FindByID(ctx, param.ID, category)
	})
	g.Go(func() error {
		if err := ccs.ccr.GetLeastCategoryCount(ctx, param.ID, count); err != nil {
			return err
		}
		if *count == 1 {
			return customerrors.NewError(
				"there is a course with only one category",
				errors.New("there is a course with only one category"),
				customerrors.InvalidAction,
			)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return ccs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ccs.ccr.UnassignCategoriesByCategoryID(ctx, param.ID); err != nil {
			return err
		}
		if err := ccs.ccr.DeleteCategory(ctx, param.ID); err != nil {
			return err
		}
		return nil
	})
}

func (ccs *CourseCategoryServiceImpl) GetCategoriesList(ctx context.Context, param entity.ListCourseCategoryParam) (*[]entity.ListCourseCategoryQuery, *int64, error) {
	categories := new([]entity.ListCourseCategoryQuery)
	totalRow := new(int64)
	if err := ccs.ccr.GetCourseCategoryList(ctx, categories, totalRow, param.Page, param.Limit, param.Search); err != nil {
		return nil, nil, err
	}
	return categories, totalRow, nil
}

func (ccs *CourseCategoryServiceImpl) CreateCategory(ctx context.Context, param entity.CreateCategoryParam) (*entity.CreateCategoryQuery, error) {
	g, ctx := errgroup.WithContext(ctx)
	category := new(entity.CourseCategory)
	newCategory := new(entity.CreateCategoryQuery)
	admin := new(entity.Admin)
	user := new(entity.User)

	g.Go(func() error {
		if err := ccs.ur.FindByID(ctx, param.AdminID, user); err != nil {
			return err
		}
		if user.Status == constants.UnverifiedStatus {
			return customerrors.NewError(
				"unauthenticate",
				errors.New("admin is not verified"),
				customerrors.Unauthenticate,
			)
		}
		return nil
	})
	g.Go(func() error {
		return ccs.ar.FindByID(ctx, param.AdminID, admin)
	})
	g.Go(func() error {
		if err := ccs.ccr.FindByName(ctx, param.Name, category); err != nil {
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
				"category already exist",
				errors.New("category already exist"),
				customerrors.ItemAlreadyExist,
			)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}
	newCategory.Name = param.Name
	if err := ccs.ccr.CreateCategory(ctx, newCategory); err != nil {
		return nil, err
	}
	return newCategory, nil
}

func (ccs *CourseCategoryServiceImpl) UpdateCategory(ctx context.Context, param entity.UpdateCategoryParam) error {
	if param.Name == nil {
		return nil
	}
	g, ctx := errgroup.WithContext(ctx)
	category := new(entity.CourseCategory)
	admin := new(entity.Admin)
	user := new(entity.User)

	g.Go(func() error {
		if err := ccs.ur.FindByID(ctx, param.AdminID, user); err != nil {
			return err
		}
		if user.Status == constants.UnverifiedStatus {
			return customerrors.NewError(
				"unauthenticate",
				errors.New("admin is not verified"),
				customerrors.Unauthenticate,
			)
		}
		return nil
	})
	g.Go(func() error {
		return ccs.ar.FindByID(ctx, param.AdminID, admin)
	})
	g.Go(func() error {
		return ccs.ccr.FindByID(ctx, param.ID, category)
	})
	g.Go(func() error {
		if err := ccs.ccr.FindByName(ctx, *param.Name, category); err != nil {
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
				"category already exist",
				errors.New("category already exist"),
				customerrors.ItemAlreadyExist,
			)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	if err := ccs.ccr.UpdateCategory(ctx, param); err != nil {
		return err
	}
	return nil
}
