package services

import (
	"context"
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
)

type AdditionalCostServiceImpl struct {
	acr *repositories.AdditionalCostRepositoryImpl
	ur  *repositories.UserRepositoryImpl
	ar  *repositories.AdminRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func CreateAdditionalCostService(
	acr *repositories.AdditionalCostRepositoryImpl,
	ur *repositories.UserRepositoryImpl,
	ar *repositories.AdminRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
) *AdditionalCostServiceImpl {
	return &AdditionalCostServiceImpl{acr, ur, ar, tmr}
}

func (acs *AdditionalCostServiceImpl) GetOperationalCost(ctx context.Context, param entity.GetOperationalCostParam) (*float64, error) {
	operationalCost := new(float64)
	user := new(entity.User)
	if err := acs.ur.FindByID(ctx, param.UserID, user); err != nil {
		return nil, err
	}
	if user.Status != "verified" {
		return nil, customerrors.NewError(
			"unverified",
			errors.New("user is not verified"),
			customerrors.Unauthenticate,
		)
	}
	if err := acs.acr.GetOperationalCost(ctx, operationalCost); err != nil {
		return nil, err
	}
	return operationalCost, nil
}

func (acs *AdditionalCostServiceImpl) GetAllAdditionalCost(ctx context.Context, param entity.GetAllAdditionalCostParam) (*[]entity.GetAdditionalCostQuery, *int64, error) {
	admin := new(entity.Admin)
	totalRow := new(int64)
	costs := new([]entity.GetAdditionalCostQuery)
	user := new(entity.User)
	if err := acs.ur.FindByID(ctx, param.AdminID, user); err != nil {
		return nil, nil, err
	}
	if user.Status == constants.UnverifiedStatus {
		return nil, nil, customerrors.NewError(
			"unauthorized",
			errors.New("admin is not verified"),
			customerrors.Unauthenticate,
		)
	}
	if err := acs.ar.FindByID(ctx, param.AdminID, admin); err != nil {
		return nil, nil, err
	}
	if err := acs.acr.GetAllAdditionalCost(ctx, param.Limit, param.Page, totalRow, costs); err != nil {
		return nil, nil, err
	}
	return costs, totalRow, nil
}

func (acs *AdditionalCostServiceImpl) DeleteCost(ctx context.Context, param entity.DeleteAdditionalCostParam) error {
	admin := new(entity.Admin)
	cost := new(entity.AdditionalCost)
	user := new(entity.User)
	if err := acs.ur.FindByID(ctx, param.AdminID, user); err != nil {
		return err
	}
	if user.Status == constants.UnverifiedStatus {
		return customerrors.NewError(
			"unauthorized",
			errors.New("admin is not verified"),
			customerrors.Unauthenticate,
		)
	}
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
	user := new(entity.User)

	if err := acs.ur.FindByID(ctx, param.AdminID, user); err != nil {
		return 0, err
	}
	if user.Status == constants.UnverifiedStatus {
		return 0, customerrors.NewError(
			"unauthorized",
			errors.New("admin is not verified"),
			customerrors.Unauthenticate,
		)
	}
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
	user := new(entity.User)

	if err := acs.ur.FindByID(ctx, param.AdminID, user); err != nil {
		return err
	}
	if user.Status == constants.UnverifiedStatus {
		return customerrors.NewError(
			"unauthorized",
			errors.New("admin is not verified"),
			customerrors.Unauthenticate,
		)
	}
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
