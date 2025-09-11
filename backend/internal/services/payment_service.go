package services

import (
	"context"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
)

type PaymentServiceImpl struct {
	pr  *repositories.PaymentRepositoryImpl
	ar  *repositories.AdminRepositoryImpl
	ur  *repositories.UserRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func CreatePaymentService(
	pr *repositories.PaymentRepositoryImpl,
	ar *repositories.AdminRepositoryImpl,
	ur *repositories.UserRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
) *PaymentServiceImpl {
	return &PaymentServiceImpl{pr, ar, ur, tmr}
}

func (ps *PaymentServiceImpl) GetAllPaymentMethod(ctx context.Context, param entity.GetAllPaymentMethodParam) (*[]entity.GetPaymentMethodQuery, *int64, error) {
	methods := new([]entity.GetPaymentMethodQuery)
	totalRow := new(int64)
	if err := ps.pr.GetAllPaymentMethod(ctx, param.Search, param.Limit, param.LastID, totalRow, methods); err != nil {
		return nil, nil, err
	}
	return methods, totalRow, nil
}

func (ps *PaymentServiceImpl) UpdatePaymentMethod(ctx context.Context, param entity.UpdatePaymentMethodParam) error {
	if param.MethodNewName == nil {
		return nil
	}
	admin := new(entity.Admin)
	count := new(int64)
	if err := ps.ar.FindByID(ctx, param.AdminID, admin); err != nil {
		return err
	}
	if err := ps.pr.FindPaymentMethodByName(ctx, *param.MethodNewName, count); err != nil {
		return err
	}
	if *count > 0 {
		return customerrors.NewError(
			"payment method with same name already exist",
			errors.New("payment method with same name already exist"),
			customerrors.InvalidAction,
		)
	}
	if err := ps.pr.UpdatePaymentMethod(ctx, param.MethodNewName, param.MethodID); err != nil {
		return err
	}
	return nil
}

func (ps *PaymentServiceImpl) DeletePaymentMethod(ctx context.Context, param entity.DeletePaymentMethodParam) error {
	admin := new(entity.Admin)
	method := new(entity.PaymentMethod)
	return ps.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ps.ar.FindByID(ctx, param.AdminID, admin); err != nil {
			return err
		}
		if err := ps.pr.FindPaymentMethodByID(ctx, param.MethodID, method); err != nil {
			return err
		}
		if err := ps.pr.UnassignPaymentMethodFromAllMentor(ctx, param.MethodID); err != nil {
			return err
		}
		if err := ps.pr.DeletePaymentMethod(ctx, param.MethodID); err != nil {
			return err
		}
		return nil
	})
}

func (ps *PaymentServiceImpl) CreatePaymentMethod(ctx context.Context, param entity.CreatePaymentMethodParam) (*int, error) {
	admin := new(entity.Admin)
	paymentMethod := new(entity.PaymentMethod)
	count := new(int64)
	if err := ps.ar.FindByID(ctx, param.AdminID, admin); err != nil {
		return nil, err
	}
	if err := ps.pr.FindPaymentMethodByName(ctx, param.MethodName, count); err != nil {
		return nil, err
	}
	if *count > 0 {
		return nil, customerrors.NewError(
			"payment method already exist",
			errors.New("payment method already exist"),
			customerrors.ItemAlreadyExist,
		)
	}
	if err := ps.pr.CreatePaymentMethod(ctx, param.MethodName, paymentMethod); err != nil {
		return nil, err
	}

	return &paymentMethod.ID, nil
}
