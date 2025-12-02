package services

import (
	"context"
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
)

type DiscountServiceImpl struct {
	dr  *repositories.DiscountRepositoryImpl
	ur  *repositories.UserRepositoryImpl
	ar  *repositories.AdminRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func CreateDiscountService(
	dr *repositories.DiscountRepositoryImpl,
	ur *repositories.UserRepositoryImpl,
	ar *repositories.AdminRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
) *DiscountServiceImpl {
	return &DiscountServiceImpl{dr, ur, ar, tmr}
}

func (ds *DiscountServiceImpl) GetDiscount(ctx context.Context, param entity.GetDiscountParam) (*float64, error) {
	user := new(entity.User)
	discount := new(entity.Discount)
	maxParticipant := new(int)

	if err := ds.ur.FindByID(ctx, param.UserID, user); err != nil {
		return nil, err
	}
	if user.Status != "verified" {
		return nil, customerrors.NewError(
			"unauthorized",
			errors.New("user is not verified"),
			customerrors.Unauthenticate,
		)
	}
	if err := ds.dr.GetMaxParticipant(ctx, maxParticipant); err != nil {
		return nil, err
	}
	participant := param.Participant
	if participant > *maxParticipant {
		participant = *maxParticipant
	}
	if err := ds.dr.GetDiscountByNumberOfParticipant(ctx, participant, discount); err != nil {
		var parsedErr *customerrors.CustomError
		if errors.As(err, &parsedErr) {
			if parsedErr.ErrCode != customerrors.ItemNotExist {
				return nil, err
			}
		}
	}
	return &discount.Amount, nil
}

func (ds *DiscountServiceImpl) GetAllDiscount(ctx context.Context, param entity.GetAllDiscountParam) (*[]entity.GetDiscountQuery, *int64, error) {
	admin := new(entity.Admin)
	totalRow := new(int64)
	discounts := new([]entity.GetDiscountQuery)
	user := new(entity.User)
	if err := ds.ur.FindByID(ctx, param.AdminID, user); err != nil {
		if user.Status == constants.UnverifiedStatus {
			return nil, nil, customerrors.NewError(
				"unauthenticate",
				errors.New("admin is not verified"),
				customerrors.Unauthenticate,
			)
		}
	}
	if err := ds.ar.FindByID(ctx, param.AdminID, admin); err != nil {
		return nil, nil, err
	}
	if err := ds.dr.GetAllDiscount(ctx, param.Limit, param.Page, totalRow, discounts); err != nil {
		return nil, nil, err
	}
	return discounts, totalRow, nil
}

func (ds *DiscountServiceImpl) DeleteDiscount(ctx context.Context, param entity.DeleteDiscountParam) error {
	admin := new(entity.Admin)
	discount := new(entity.Discount)
	user := new(entity.User)
	if err := ds.ur.FindByID(ctx, param.AdminID, user); err != nil {
		return err
	}
	if user.Status == constants.UnverifiedStatus {
		return customerrors.NewError(
			"unauthenticate",
			errors.New("admin is not verified"),
			customerrors.Unauthenticate,
		)
	}
	if err := ds.ar.FindByID(ctx, param.AdminID, admin); err != nil {
		return err
	}
	if err := ds.dr.FindByID(ctx, param.DiscountID, discount); err != nil {
		return err
	}
	if err := ds.dr.DeleteDiscount(ctx, param.DiscountID); err != nil {
		return err
	}
	return nil
}

func (ds *DiscountServiceImpl) UpdateDiscountAmount(ctx context.Context, param entity.UpdateDiscountParam) error {
	admin := new(entity.Admin)
	discount := new(entity.Discount)
	user := new(entity.User)
	if err := ds.ur.FindByID(ctx, param.AdminID, user); err != nil {
		return err
	}
	if user.Status == constants.UnverifiedStatus {
		return customerrors.NewError(
			"unauthenticate",
			errors.New("admin is not verified"),
			customerrors.Unauthenticate,
		)
	}
	if err := ds.ar.FindByID(ctx, param.AdminID, admin); err != nil {
		return err
	}
	if err := ds.dr.FindByID(ctx, param.DiscountID, discount); err != nil {
		return err
	}
	if err := ds.dr.UpdateAmount(ctx, param.DiscountID, param.Amount); err != nil {
		return err
	}
	return nil
}

func (ds *DiscountServiceImpl) CreateNewDiscount(ctx context.Context, param entity.CreateNewDiscountParam) (int, error) {
	admin := new(entity.Admin)
	id := new(int)
	discount := new(entity.Discount)
	user := new(entity.User)
	if err := ds.ur.FindByID(ctx, param.AdminID, user); err != nil {
		return 0, err
	}
	if user.Status == constants.UnverifiedStatus {
		return 0, customerrors.NewError(
			"unauthenticate",
			errors.New("admin is not verified"),
			customerrors.Unauthenticate,
		)
	}
	if err := ds.ar.FindByID(ctx, param.AdminID, admin); err != nil {
		return 0, err
	}
	if err := ds.dr.GetDiscountByNumberOfParticipant(ctx, param.NumberOfParticipant, discount); err != nil {
		var parsedErr *customerrors.CustomError
		if errors.As(err, &parsedErr) {
			if parsedErr.ErrCode != customerrors.ItemNotExist {
				return 0, err
			}
		}
	} else {
		return 0, customerrors.NewError(
			"discount already exist",
			errors.New("discount already exist"),
			customerrors.ItemAlreadyExist,
		)
	}
	if err := ds.dr.CreateNewDiscount(ctx, param.NumberOfParticipant, param.Amount, id); err != nil {
		return 0, err
	}
	return *id, nil
}
