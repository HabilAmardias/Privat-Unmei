package services

import (
	"context"
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/utils"
)

type AdminServiceImpl struct {
	ur  *repositories.UserRepositoryImpl
	ar  *repositories.AdminRepositoryImpl
	sr  *repositories.StudentRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
	bu  *utils.BcryptUtil
	ju  *utils.JWTUtil
	gu  *utils.GomailUtil
}

func CreateAdminService(
	ur *repositories.UserRepositoryImpl,
	ar *repositories.AdminRepositoryImpl,
	sr *repositories.StudentRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
	bu *utils.BcryptUtil,
	ju *utils.JWTUtil,
	gu *utils.GomailUtil,
) *AdminServiceImpl {
	return &AdminServiceImpl{ur, ar, sr, tmr, bu, ju, gu}
}

func (as *AdminServiceImpl) GetStudentList(ctx context.Context, param entity.ListStudentParam) (*[]entity.ListStudentQuery, *int64, error) {
	students := new([]entity.ListStudentQuery)
	totalRow := new(int64)

	if err := as.sr.GetStudentList(ctx, totalRow, param.Limit, param.Page, students); err != nil {
		return nil, nil, err
	}
	return students, totalRow, nil
}

func (as *AdminServiceImpl) Login(ctx context.Context, param entity.AdminLoginParam) (string, error) {
	user := new(entity.User)
	admin := new(entity.Admin)
	token := new(string)

	if err := as.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := as.ur.FindByEmail(ctx, param.Email, user); err != nil {
			return err
		}
		if err := as.ar.FindByID(ctx, user.ID, admin); err != nil {
			return err
		}
		if match := as.bu.ComparePassword(param.Password, user.Password); !match {
			return customerrors.NewError(
				"invalid email or password",
				errors.New("password does not match"),
				customerrors.InvalidAction,
			)
		}
		loginToken, err := as.ju.GenerateJWT(admin.ID, constants.AdminRole, constants.ForLogin, user.Status)
		if err != nil {
			return err
		}

		*token = loginToken

		return nil
	}); err != nil {
		return "", err
	}
	return *token, nil
}
