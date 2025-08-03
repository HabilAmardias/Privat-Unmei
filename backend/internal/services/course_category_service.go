package services

import (
	"context"
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
