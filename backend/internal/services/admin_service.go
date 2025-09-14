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
	mr  *repositories.MentorRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
	cu  *utils.CloudinaryUtil
	bu  *utils.BcryptUtil
	ju  *utils.JWTUtil
	gu  *utils.GomailUtil
}

func CreateAdminService(
	ur *repositories.UserRepositoryImpl,
	ar *repositories.AdminRepositoryImpl,
	sr *repositories.StudentRepositoryImpl,
	mr *repositories.MentorRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
	cu *utils.CloudinaryUtil,
	bu *utils.BcryptUtil,
	ju *utils.JWTUtil,
	gu *utils.GomailUtil,
) *AdminServiceImpl {
	return &AdminServiceImpl{ur, ar, sr, mr, tmr, cu, bu, ju, gu}
}

func (as *AdminServiceImpl) UpdatePassword(ctx context.Context, param entity.AdminUpdatePasswordParam) error {
	user := new(entity.User)
	admin := new(entity.Admin)
	if err := as.ur.FindByID(ctx, param.AdminID, user); err != nil {
		return err
	}
	if user.Status != constants.VerifiedStatus {
		return customerrors.NewError(
			"admin is unverified",
			errors.New("admin is unverified"),
			customerrors.Unauthenticate,
		)
	}
	if err := as.ar.FindByID(ctx, param.AdminID, admin); err != nil {
		return err
	}
	if match := as.bu.ComparePassword(param.Password, user.Password); match {
		return customerrors.NewError(
			"cannot change into same password",
			errors.New("cannot change into same password"),
			customerrors.InvalidAction,
		)
	}
	hashedPass, err := as.bu.HashPassword(param.Password)
	if err != nil {
		return err
	}
	if err := as.ar.ChangePassword(ctx, param.AdminID, hashedPass); err != nil {
		return err
	}
	return nil
}

func (as *AdminServiceImpl) GenerateRandomPassword() (string, error) {
	pass, err := generateRandomPassword()
	if err != nil {
		return "", customerrors.NewError("failed to generate password", err, customerrors.CommonErr)
	}
	return pass, nil
}

func (as *AdminServiceImpl) VerifyAdmin(ctx context.Context, param entity.AdminVerificationParam) error {
	user := new(entity.User)
	admin := new(entity.Admin)
	if err := as.ur.FindByID(ctx, param.AdminID, user); err != nil {
		return err
	}
	if user.Status == constants.VerifiedStatus {
		return customerrors.NewError(
			"admin is already verified",
			errors.New("admin is already verified"),
			customerrors.InvalidAction,
		)
	}
	if err := as.ar.FindByID(ctx, param.AdminID, admin); err != nil {
		return err
	}
	if match := as.bu.ComparePassword(param.Password, user.Password); match {
		return customerrors.NewError(
			"cannot change into same password",
			errors.New("cannot change into same password"),
			customerrors.InvalidAction,
		)
	}
	hashedPass, err := as.bu.HashPassword(param.Password)
	if err != nil {
		return err
	}
	if err := as.ar.VerifyAdmin(ctx, param.AdminID, param.Email, hashedPass); err != nil {
		return err
	}
	return nil
}

func (as *AdminServiceImpl) Login(ctx context.Context, param entity.AdminLoginParam) (*string, *string, error) {
	user := new(entity.User)
	admin := new(entity.Admin)
	status := new(string)
	token := new(string)

	if err := as.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := as.ur.FindByEmail(ctx, param.Email, user); err != nil {
			var parsedErr *customerrors.CustomError
			if errors.As(err, &parsedErr) {
				if parsedErr.ErrCode == customerrors.ItemNotExist {
					return customerrors.NewError(
						"invalid email or password",
						parsedErr.ErrLog,
						customerrors.InvalidAction,
					)
				}
			}
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
		*status = user.Status

		return nil
	}); err != nil {
		return nil, nil, err
	}
	return token, status, nil
}
