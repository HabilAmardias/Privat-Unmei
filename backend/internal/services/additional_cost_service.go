package services

import (
	"context"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
)

type AdditionalCostServiceImpl struct {
	acr *repositories.AdditionalCostRepositoryImpl
	ar  *repositories.AdminRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func CreateAdditionalCostService(
	acr *repositories.AdditionalCostRepositoryImpl,
	ar *repositories.AdminRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
) *AdditionalCostServiceImpl {
	return &AdditionalCostServiceImpl{acr, ar, tmr}
}

func (acs *AdditionalCostServiceImpl) DeleteCost(ctx context.Context, param entity.DeleteAdditionalCostParam) error {
	admin := new(entity.Admin)
	cost := new(entity.AdditionalCost)

	if err := acs.ar.FindByID(ctx, param.AdminID, admin); err != nil {
		return err
	}
	if err := acs.acr.FindByID(ctx, param.CostID, cost); err != nil {
		return err
	}
	if err := acs.acr.DeleteCost(ctx, param.CostID); err != nil {
		return err
	}
	return nil
}

func (acs *AdditionalCostServiceImpl) CreateNewAdditionalCost(ctx context.Context, param entity.CreateAdditionalCostParam) (int, error) {
	admin := new(entity.Admin)
	id := new(int)
	count := new(int64)

	if err := acs.ar.FindByID(ctx, param.AdminID, admin); err != nil {
		return 0, err
	}
	if err := acs.acr.FindByName(ctx, param.Name, count); err != nil {
		return 0, err
	}
	if *count > 0 {
		return 0, customerrors.NewError(
			"operational cost already exist",
			errors.New("operational cost already exist"),
			customerrors.ItemAlreadyExist,
		)
	}
	if err := acs.acr.CreateOperationalCost(ctx, param.Name, param.Amount, id); err != nil {
		return 0, err
	}
	return *id, nil
}
func (acs *AdditionalCostServiceImpl) UpdateCostAmount(ctx context.Context, param entity.UpdateAdditonalCostParam) error {
	admin := new(entity.Admin)
	cost := new(entity.AdditionalCost)

	if err := acs.ar.FindByID(ctx, param.AdminID, admin); err != nil {
		return err
	}
	if err := acs.acr.FindByID(ctx, param.CostID, cost); err != nil {
		return err
	}
	if err := acs.acr.UpdateCostAmount(ctx, param.CostID, param.Amount); err != nil {
		return err
	}
	return nil
}
