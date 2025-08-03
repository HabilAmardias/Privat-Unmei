package services

import (
	"context"
	"fmt"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/utils"
)

type MentorServiceImpl struct {
	tmr *repositories.TransactionManagerRepositories
	ur  *repositories.UserRepositoryImpl
	mr  *repositories.MentorRepositoryImpl
	bu  *utils.BcryptUtil
	ju  *utils.JWTUtil
	cu  *utils.CloudinaryUtil
	gu  *utils.GomailUtil
}

func CreateMentorService(
	tmr *repositories.TransactionManagerRepositories,
	ur *repositories.UserRepositoryImpl,
	mr *repositories.MentorRepositoryImpl,
	bu *utils.BcryptUtil,
	ju *utils.JWTUtil,
	cu *utils.CloudinaryUtil,
	gu *utils.GomailUtil,
) *MentorServiceImpl {
	return &MentorServiceImpl{tmr, ur, mr, bu, ju, cu, gu}
}

func (ms *MentorServiceImpl) GetMentorList(ctx context.Context, param entity.ListMentorParam) (*[]entity.ListMentorQuery, *int64, error) {
	mentors := new([]entity.ListMentorQuery)
	totalRow := new(int64)
	if err := ms.mr.GetMentorList(ctx, mentors, totalRow, param); err != nil {
		return nil, nil, err
	}
	return mentors, totalRow, nil
}

func (ms *MentorServiceImpl) DeleteMentor(ctx context.Context, param entity.DeleteMentorParam) error {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	return ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ms.ur.FindByID(ctx, param.ID, user); err != nil {
			return err
		}
		if err := ms.mr.FindByID(ctx, user.ID, mentor); err != nil {
			return err
		}
		if err := ms.mr.DeleteMentor(ctx, mentor.ID); err != nil {
			return err
		}
		if err := ms.ur.DeleteUser(ctx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (ms *MentorServiceImpl) UpdateMentorForAdmin(ctx context.Context, param entity.UpdateMentorParam) error {
	user := new(entity.User)
	mentor := new(entity.Mentor)
	return ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ms.ur.FindByID(ctx, param.ID, user); err != nil {
			return err
		}
		if err := ms.mr.FindByID(ctx, user.ID, mentor); err != nil {
			return err
		}
		if err := ms.mr.UpdateMentor(ctx, mentor.ID, &param.UpdateMentorQuery); err != nil {
			return err
		}
		return nil
	})
}

func (ms *MentorServiceImpl) AddNewMentor(ctx context.Context, param entity.AddNewMentorParam) error {
	user := new(entity.User)
	mentor := new(entity.Mentor)

	return ms.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := ms.ur.FindByEmail(ctx, param.Email, user); err != nil {
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

		user.Email = param.Email
		user.Name = param.Name
		user.Bio = param.Bio
		user.ProfileImage = constants.DefaultAvatar

		hashedPass, err := ms.bu.HashPassword(param.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPass

		// Will be changed to verified after mentor change their password for the first time
		user.Status = constants.UnverifiedStatus

		if err := ms.ur.AddNewUser(ctx, user); err != nil {
			return err
		}

		mentor.ID = user.ID
		mentor.Campus = param.Campus
		mentor.Degree = param.Degree
		mentor.Major = param.Major
		mentor.WhatsappNumber = param.WhatsappNumber
		mentor.YearsOfExperience = param.YearsOfExperience

		newFilename := fmt.Sprintf("%s.pdf", mentor.ID)
		uploadRes, err := ms.cu.UploadFile(ctx, param.ResumeFile, newFilename, constants.ResumeFolder)
		if err != nil {
			return err
		}
		mentor.Resume = uploadRes.URL
		if err := ms.mr.AddNewMentor(ctx, mentor); err != nil {
			return err
		}

		emailParam := entity.SendEmailParams{
			Receiver:  param.Email,
			Subject:   "Login Credentials - Privat Unmei",
			EmailBody: constants.SendMentorAccEmailBody(param.Email, param.Password),
		}
		if err := ms.gu.SendEmail(emailParam); err != nil {
			return err
		}

		return nil
	})
}
