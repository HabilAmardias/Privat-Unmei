package services

import (
	"context"
	"errors"
	"log"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/utils"
)

type StudentServiceImpl struct {
	ur  *repositories.UserRepositoryImpl
	sr  *repositories.StudentRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
	bu  *utils.BcryptUtil
	gu  *utils.GomailUtil
	cu  *utils.CloudinaryUtil
	ju  *utils.JWTUtil
}

func CreateStudentService(
	ur *repositories.UserRepositoryImpl,
	sr *repositories.StudentRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
	bu *utils.BcryptUtil,
	gu *utils.GomailUtil,
	cu *utils.CloudinaryUtil,
	ju *utils.JWTUtil,
) *StudentServiceImpl {
	return &StudentServiceImpl{ur, sr, tmr, bu, gu, cu, ju}
}

func (us *StudentServiceImpl) Login(ctx context.Context, param entity.StudentLoginParam) (string, error) {
	user := new(entity.User)
	student := new(entity.Student)
	token := new(string)

	if err := us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.FindByEmail(ctx, param.Email, user); err != nil {
			return err
		}
		if err := us.sr.FindByID(ctx, user.ID, student); err != nil {
			return err
		}
		if match := us.bu.ComparePassword(param.Password, user.Password); !match {
			return customerrors.NewError("invalid email or password", errors.New("invalid email or password"), customerrors.InvalidAction)
		}
		jwt, err := us.ju.GenerateJWT(student.ID, constants.StudentRole, constants.ForLogin, user.Status)
		if err != nil {
			return err
		}
		*token = jwt

		return nil

	}); err != nil {
		return "", err
	}
	return *token, nil
}

func (us *StudentServiceImpl) Register(ctx context.Context, param entity.StudentRegisterParam) error {
	user := new(entity.User)
	student := new(entity.Student)

	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.FindByEmail(ctx, param.Email, user); err != nil {
			if err.Error() != customerrors.UserNotFound {
				return err
			}
		} else {
			return customerrors.NewError(
				"user already exist",
				customerrors.ErrItemAlreadyExist,
				customerrors.ItemAlreadyExist,
			)
		}

		user.Bio = param.Bio
		user.Email = param.Email
		user.Name = param.Name
		user.Status = param.Status
		user.ProfileImage = constants.DefaultAvatar

		hashed, err := us.bu.HashPassword(param.Password)
		if err != nil {
			return err
		}
		user.Password = hashed
		if err := us.ur.AddNewUser(ctx, user); err != nil {
			return err
		}
		token, err := us.ju.GenerateJWT(user.ID, constants.StudentRole, constants.ForVerification, user.Status)
		if err != nil {
			return err
		}
		student.ID = user.ID
		student.VerifyToken = &token
		if err := us.sr.AddNewStudent(ctx, student); err != nil {
			return err
		}

		// wrapped this with go func to make other request does not get blocked by this
		go func() {
			param := entity.SendEmailParams{
				Receiver:  param.Email,
				Subject:   "Verify your account",
				EmailBody: constants.VerificationEmailBody(token),
			}
			if err := us.gu.SendEmail(param); err != nil {
				log.Println(err.Error())
				return
			}
			log.Println("Send Email Success")
		}()

		return nil
	})
}
