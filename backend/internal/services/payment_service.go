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

func (ps *PaymentServiceImpl) CreatePaymentMethod(ctx context.Context, param entity.CreatePaymentMethodParam) (*int, error) {
	user := new(entity.User)
	admin := new(entity.Admin)
	paymentMethod := new(entity.PaymentMethod)
	count := new(int64)

	if err := ps.ur.FindByID(ctx, param.AdminID, user); err != nil {
		return nil, err
	}
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
