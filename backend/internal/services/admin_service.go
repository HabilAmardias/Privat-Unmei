package services

import (
	"context"
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/utils"

	"golang.org/x/sync/errgroup"
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

func (as *AdminServiceImpl) AdminProfile(ctx context.Context, param entity.AdminProfileParam) (*entity.AdminProfileQuery, error) {
	g, ctx := errgroup.WithContext(ctx)
	user := new(entity.User)
	admin := new(entity.Admin)
	query := new(entity.AdminProfileQuery)
	g.Go(func() error {
		return as.ur.FindByID(ctx, param.AdminID, user)
	})
	g.Go(func() error {
		return as.ar.FindByID(ctx, param.AdminID, admin)
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}
	query.Name = user.Name
	query.Email = user.Email
	query.Bio = user.Bio
	query.ProfileImage = user.ProfileImage
	query.Status = user.Status

	return query, nil
}

func (as *AdminServiceImpl) UpdatePassword(ctx context.Context, param entity.AdminUpdatePasswordParam) error {
	g, ctx := errgroup.WithContext(ctx)
	user := new(entity.User)
	admin := new(entity.Admin)
	g.Go(func() error {
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
		return nil
	})
	g.Go(func() error {
		return as.ar.FindByID(ctx, param.AdminID, admin)
	})
	if err := g.Wait(); err != nil {
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
	g, ctx := errgroup.WithContext(ctx)
	user := new(entity.User)
	admin := new(entity.Admin)
	g.Go(func() error {
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
		return nil
	})
	g.Go(func() error {
		return as.ar.FindByID(ctx, param.AdminID, admin)
	})
	if err := g.Wait(); err != nil {
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

func (as *AdminServiceImpl) Login(ctx context.Context, param entity.AdminLoginParam) (*string, *string, *string, error) {
	user := new(entity.User)
	admin := new(entity.Admin)
	status := new(string)
	authToken := new(string)
	refreshToken := new(string)

	if err := as.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := as.ur.FindByEmail(ctx, param.Email, user); err != nil {
			var parsedErr *customerrors.CustomError
			if !errors.As(err, &parsedErr) {
				return customerrors.NewError(
					"something went wrong",
					errors.New("cannot parse error"),
					customerrors.CommonErr,
				)
			}
			if parsedErr.ErrCode == customerrors.ItemNotExist {
				return customerrors.NewError(
					"invalid email or password",
					parsedErr.ErrLog,
					customerrors.InvalidAction,
				)
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
		atoken, err := as.ju.GenerateJWT(admin.ID, constants.AdminRole, constants.ForLogin, user.Status, constants.AUTH_AGE)
		if err != nil {
			return err
		}

		rtoken, err := as.ju.GenerateJWT(admin.ID, constants.AdminRole, constants.ForRefresh, user.Status, constants.REFRESH_AGE)
		if err != nil {
			return err
		}

		*authToken = atoken
		*refreshToken = rtoken
		*status = user.Status

		return nil
	}); err != nil {
		return nil, nil, nil, err
	}
	return authToken, refreshToken, status, nil
}
