package services

import (
	"context"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
)

type CourseCategoryServiceImpl struct {
	ccr *repositories.CourseCategoryRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func CreateCourseCategoryService(ccr *repositories.CourseCategoryRepositoryImpl, tmr *repositories.TransactionManagerRepositories) *CourseCategoryServiceImpl {
	return &CourseCategoryServiceImpl{ccr, tmr}
}

func (ccs *CourseCategoryServiceImpl) GetCategoriesList(ctx context.Context, param entity.ListCourseCategoryParam) (*[]entity.ListCourseCategoryQuery, *int64, error) {
	categories := new([]entity.ListCourseCategoryQuery)
	totalRow := new(int64)
	if err := ccs.ccr.GetCourseCategoryList(ctx, categories, totalRow, param); err != nil {
		return nil, nil, err
	}
	return categories, totalRow, nil
}

func (ccs *CourseCategoryServiceImpl) CreateCategory(ctx context.Context, param entity.CreateCategoryParam) (*entity.CreateCategoryQuery, error) {
	category := new(entity.CourseCategory)
	newCategory := new(entity.CreateCategoryQuery)

	if err := ccs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ccs.ccr.FindByName(ctx, param.Name, category); err != nil {
			if err.Error() != "category does not exist" {
				return err
			}
		} else {
			return customerrors.NewError(
				"category already exist",
				errors.New("category already exist"),
				customerrors.ItemAlreadyExist,
			)
		}
		newCategory.Name = param.Name
		if err := ccs.ccr.CreateCategory(ctx, newCategory); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return newCategory, nil
}
