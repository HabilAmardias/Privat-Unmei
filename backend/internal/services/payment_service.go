package services

import (
	"context"
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
)

type PaymentServiceImpl struct {
	pr  *repositories.PaymentRepositoryImpl
	ar  *repositories.AdminRepositoryImpl
	ur  *repositories.UserRepositoryImpl
	mr  *repositories.MentorRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func CreatePaymentService(
	pr *repositories.PaymentRepositoryImpl,
	ar *repositories.AdminRepositoryImpl,
	ur *repositories.UserRepositoryImpl,
	mr *repositories.MentorRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
) *PaymentServiceImpl {
	return &PaymentServiceImpl{pr, ar, ur, mr, tmr}
}

func (ps *PaymentServiceImpl) GetMentorPaymentMethod(ctx context.Context, param entity.GetMentorPaymentMethodParam) (*[]entity.GetPaymentMethodQuery, error) {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	methods := new([]entity.GetPaymentMethodQuery)
	if err := ps.ur.FindByID(ctx, param.UserID, user); err != nil {
		return nil, err
	}
	if err := ps.mr.FindByID(ctx, param.MentorID, mentor, false); err != nil {
		return nil, err
	}
	if err := ps.pr.GetMentorPaymentMethod(ctx, param.MentorID, methods); err != nil {
		return nil, err
	}
	if len(*methods) == 0 {
		return nil, customerrors.NewError(
			"mentor payment method not found",
			errors.New("mentor payment method not found"),
			customerrors.ItemNotExist,
		)
	}
	return methods, nil
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
	user := new(entity.User)
	if err := ps.ur.FindByID(ctx, param.AdminID, user); err != nil {
		return err
	}
	if user.Status == constants.UnverifiedStatus {
		return customerrors.NewError(
			"unauthenticate",
			errors.New("admin is not verified"),
			customerrors.Unauthenticate,
		)
	}
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
	user := new(entity.User)
	return ps.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ps.ur.FindByID(ctx, param.AdminID, user); err != nil {
			return err
		}
		if user.Status == constants.UnverifiedStatus {
			return customerrors.NewError(
				"unauthenticate",
				errors.New("admin is not verified"),
				customerrors.Unauthenticate,
			)
		}
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
	user := new(entity.User)
	if err := ps.ur.FindByID(ctx, param.AdminID, user); err != nil {
		return nil, err
	}
	if user.Status == constants.UnverifiedStatus {
		return nil, customerrors.NewError(
			"unauthenticate",
			errors.New("admin is not verified"),
			customerrors.Unauthenticate,
		)
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
