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

func (ps *PaymentServiceImpl) GetMentorPaymentMethod(ctx context.Context, param entity.GetMentorPaymentMethodParam) (*[]entity.GetMentorPaymentMethodQuery, error) {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	methods := new([]entity.GetMentorPaymentMethodQuery)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return ps.ur.FindByID(ctx, param.UserID, user)
	})
	g.Go(func() error {
		return ps.mr.FindByID(ctx, param.MentorID, mentor, false)
	})
	g.Go(func() error {
		return ps.pr.GetMentorPaymentMethod(ctx, param.MentorID, methods)
	})
	if err := g.Wait(); err != nil {
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
	if err := ps.pr.GetAllPaymentMethod(ctx, param.Search, param.Limit, param.Page, totalRow, methods); err != nil {
		return nil, nil, err
	}
	return methods, totalRow, nil
}

func (ps *PaymentServiceImpl) UpdatePaymentMethod(ctx context.Context, param entity.UpdatePaymentMethodParam) error {
	if param.MethodNewName == nil {
		return nil
	}
	g, ctx := errgroup.WithContext(ctx)
	admin := new(entity.Admin)
	count := new(int64)
	user := new(entity.User)

	g.Go(func() error {
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
		return nil
	})
	g.Go(func() error {
		return ps.ar.FindByID(ctx, param.AdminID, admin)
	})
	g.Go(func() error {
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
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
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
	count := new(int)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
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
		return nil
	})
	g.Go(func() error {
		return ps.ar.FindByID(ctx, param.AdminID, admin)
	})
	g.Go(func() error {
		return ps.pr.FindPaymentMethodByID(ctx, param.MethodID, method)
	})
	g.Go(func() error {
		if err := ps.pr.GetLeastPaymentMethodCount(ctx, param.MethodID, count); err != nil {
			return err
		}
		if *count == 1 {
			return customerrors.NewError(
				"there is a mentor with only one method",
				errors.New("there is a mentor with only one method"),
				customerrors.InvalidAction,
			)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return ps.tmr.WithTransaction(ctx, func(ctx context.Context) error {
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
	g, ctx := errgroup.WithContext(ctx)

	admin := new(entity.Admin)
	paymentMethod := new(entity.PaymentMethod)
	count := new(int64)
	user := new(entity.User)

	g.Go(func() error {
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
		return nil
	})
	g.Go(func() error {
		return ps.ar.FindByID(ctx, param.AdminID, admin)
	})
	g.Go(func() error {
		if err := ps.pr.FindPaymentMethodByName(ctx, param.MethodName, count); err != nil {
			return err
		}
		if *count > 0 {
			return customerrors.NewError(
				"payment method already exist",
				errors.New("payment method already exist"),
				customerrors.ItemAlreadyExist,
			)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}
	if err := ps.pr.CreatePaymentMethod(ctx, param.MethodName, paymentMethod); err != nil {
		return nil, err
	}

	return &paymentMethod.ID, nil
}
