package services

import (
	"context"
	"log"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/utils"
)

type StudentServiceImpl struct {
	ur  *repositories.UserRepositoryImpl
	sr  *repositories.StudentRepository
	tmr *repositories.TransactionManagerRepositories
	bu  *utils.BcryptUtil
	gu  *utils.GomailUtil
}

func (us *StudentServiceImpl) Register(ctx context.Context, param entity.StudentRegisterParam) error {
	user := new(entity.User)
	student := new(entity.Student)

	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.FindByEmail(ctx, param.Email, user); err != nil {
			if err.Error() != customerrors.ErrRecordNotFound.Error() {
				return err
			}
		} else {
			return customerrors.NewError(
				customerrors.ErrItemAlreadyExist,
				customerrors.ErrItemAlreadyExist,
				customerrors.ItemAlreadyExist,
			)
		}

		user.Bio = param.Bio
		user.Email = param.Email
		user.Name = param.Name
		user.Status = param.Status

		//TODO: Add upload profile photo logic (still thinking on where to save all files uploaded to the platform)

		hashed, err := us.bu.HashPassword(param.Password)
		if err != nil {
			return err
		}
		user.Password = hashed
		if err := us.ur.AddNewUser(ctx, user); err != nil {
			return err
		}
		student.ID = user.ID
		if err := us.sr.AddNewStudent(ctx, student); err != nil {
			return err
		}

		// wrapped in go func to make other request does not get blocked by this
		// TODO: Create Email Body Template
		go func() {
			param := entity.SendEmailParams{
				Receiver:  param.Email,
				Subject:   "Verify account",
				EmailBody: "",
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
