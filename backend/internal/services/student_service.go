package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/utils"
	"strings"
)

type StudentServiceImpl struct {
	ur     *repositories.UserRepositoryImpl
	sr     *repositories.StudentRepositoryImpl
	ar     *repositories.AdminRepositoryImpl
	tmr    *repositories.TransactionManagerRepositories
	bu     *utils.BcryptUtil
	gu     *utils.GomailUtil
	cu     *utils.CloudinaryUtil
	ju     *utils.JWTUtil
	goauth *utils.GoogleOauth
}

func CreateStudentService(
	ur *repositories.UserRepositoryImpl,
	sr *repositories.StudentRepositoryImpl,
	ar *repositories.AdminRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
	bu *utils.BcryptUtil,
	gu *utils.GomailUtil,
	cu *utils.CloudinaryUtil,
	ju *utils.JWTUtil,
	goauth *utils.GoogleOauth,
) *StudentServiceImpl {
	return &StudentServiceImpl{ur, sr, ar, tmr, bu, gu, cu, ju, goauth}
}

func (us *StudentServiceImpl) GoogleLogin(oauthState string) string {
	return us.goauth.Config.AuthCodeURL(oauthState)
}

func (us *StudentServiceImpl) RefreshToken(ctx context.Context, param entity.RefreshTokenParam) (string, error) {
	user := new(entity.User)
	if err := us.ur.FindByID(ctx, param.UserID, user); err != nil {
		return "", err
	}

	token, err := us.ju.GenerateJWT(param.UserID, param.Role, constants.ForLogin, user.Status, constants.AUTH_AGE)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (us *StudentServiceImpl) GetStudentProfile(ctx context.Context, param entity.StudentProfileParam) (*entity.StudentProfileQuery, error) {
	user := new(entity.User)
	student := new(entity.Student)
	profile := new(entity.StudentProfileQuery)
	if err := us.ur.FindByID(ctx, param.ID, user); err != nil {
		return nil, err
	}
	if err := us.sr.FindByID(ctx, user.ID, student); err != nil {
		return nil, err
	}
	profile.ID = param.ID
	profile.Name = user.Name
	profile.Bio = user.Bio
	profile.ProfileImage = user.ProfileImage
	profile.Email = user.Email

	return profile, nil
}

func (us *StudentServiceImpl) ChangePassword(ctx context.Context, param entity.StudentChangePasswordParam) error {
	user := new(entity.User)
	student := new(entity.Student)

	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.FindByID(ctx, param.ID, user); err != nil {
			return err
		}
		if err := us.sr.FindByID(ctx, param.ID, student); err != nil {
			return err
		}
		if match := us.bu.ComparePassword(param.NewPassword, user.Password); match {
			return customerrors.NewError(
				"cannot change into same password",
				errors.New("new password same as previous password"),
				customerrors.InvalidAction,
			)
		}
		hashedPass, err := us.bu.HashPassword(param.NewPassword)
		if err != nil {
			return err
		}
		if err := us.ur.UpdateUserPassword(ctx, hashedPass, param.ID); err != nil {
			return err
		}
		return nil
	})
}

func (us *StudentServiceImpl) GoogleLoginCallback(ctx context.Context, code string) (string, string, string, error) {
	oauthToken, err := us.goauth.Config.Exchange(context.Background(), code)
	if err != nil {
		return "", "", "", customerrors.NewError(
			"failed to login",
			err,
			customerrors.CommonErr,
		)
	}
	response, err := http.Get(us.goauth.UrlAPI + oauthToken.AccessToken)
	if err != nil {
		return "", "", "", customerrors.NewError(
			"failed to login",
			err,
			customerrors.CommonErr,
		)
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", "", customerrors.NewError(
			"failed to login",
			err,
			customerrors.CommonErr,
		)
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(contents, &userInfo); err != nil {
		return "", "", "", customerrors.NewError(
			"failed to login",
			err,
			customerrors.CommonErr,
		)
	}
	email, ok := userInfo["email"].(string)
	if !ok {
		return "", "", "", customerrors.NewError(
			"no email found",
			errors.New("email not found in user info"),
			customerrors.ItemNotExist,
		)
	}
	user := new(entity.User)
	student := new(entity.Student)
	if err := us.ur.FindByEmail(ctx, email, user); err != nil {
		var parsedErr *customerrors.CustomError
		errors.As(err, &parsedErr)
		if parsedErr.ErrCode == customerrors.ItemNotExist {
			pass, err := generateRandomPassword()
			if err != nil {
				return "", "", "", customerrors.NewError(
					"error when creating account",
					err,
					customerrors.CommonErr,
				)
			}
			hashed, err := us.bu.HashPassword(pass)
			if err != nil {
				return "", "", "", err
			}
			newUser := &entity.User{
				Email:        email,
				Name:         strings.Split(email, "@")[0],
				Status:       constants.VerifiedStatus,
				Password:     hashed,
				ProfileImage: constants.DefaultAvatar,
			}
			log.Println(newUser.Name)
			if err := us.ur.AddNewUser(ctx, newUser); err != nil {
				return "", "", "", err
			}
			newStudent := &entity.Student{
				ID: newUser.ID,
			}
			if err := us.sr.AddNewStudent(ctx, newStudent); err != nil {
				return "", "", "", err
			}
			authToken, err := us.ju.GenerateJWT(newUser.ID, constants.StudentRole, constants.ForLogin, constants.VerifiedStatus, constants.AUTH_AGE)
			if err != nil {
				return "", "", "", err
			}
			refreshToken, err := us.ju.GenerateJWT(newStudent.ID, constants.StudentRole, constants.ForRefresh, constants.VerifiedStatus, constants.REFRESH_AGE)
			if err != nil {
				return "", "", "", err
			}
			return authToken, refreshToken, newUser.Status, nil
		}
		return "", "", "", err
	}
	if err := us.sr.FindByID(ctx, user.ID, student); err != nil {
		return "", "", "", err
	}
	authToken, err := us.ju.GenerateJWT(student.ID, constants.StudentRole, constants.ForLogin, constants.VerifiedStatus, constants.AUTH_AGE)
	if err != nil {
		return "", "", "", err
	}
	refreshToken, err := us.ju.GenerateJWT(student.ID, constants.StudentRole, constants.ForRefresh, constants.VerifiedStatus, constants.REFRESH_AGE)
	if err != nil {
		return "", "", "", err
	}
	return authToken, refreshToken, user.Status, nil
}

func (us *StudentServiceImpl) UpdateStudentProfile(ctx context.Context, param entity.UpdateStudentParam) error {
	user := new(entity.User)
	student := new(entity.Student)
	updateQuery := new(entity.UpdateUserQuery)

	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.FindByID(ctx, param.ID, user); err != nil {
			return err
		}
		if err := us.sr.FindByID(ctx, param.ID, student); err != nil {
			return err
		}

		updateQuery.Name = param.Name
		updateQuery.Bio = param.Bio

		if param.ProfileImage != nil {
			filename := param.ID
			res, err := us.cu.UploadFile(ctx, param.ProfileImage, filename, constants.AvatarFolder)
			if err != nil {
				return err
			}
			updateQuery.ProfileImage = &res.SecureURL
		}
		if err := us.ur.UpdateUserProfile(ctx, updateQuery, param.ID); err != nil {
			return err
		}
		return nil
	})
}

func (us *StudentServiceImpl) GetStudentList(ctx context.Context, param entity.ListStudentParam) (*[]entity.ListStudentQuery, *int64, error) {
	students := new([]entity.ListStudentQuery)
	totalRow := new(int64)
	user := new(entity.User)
	admin := new(entity.Admin)
	if err := us.ur.FindByID(ctx, param.AdminID, user); err != nil {
		return nil, nil, err
	}
	if user.Status == constants.UnverifiedStatus {
		return nil, nil, customerrors.NewError(
			"unauthorized",
			errors.New("admin is not verified"),
			customerrors.Unauthenticate,
		)
	}
	if err := us.ar.FindByID(ctx, param.AdminID, admin); err != nil {
		return nil, nil, err
	}
	if err := us.sr.GetStudentList(ctx, totalRow, param.Limit, param.Page, students); err != nil {
		return nil, nil, err
	}
	return students, totalRow, nil
}

func (us *StudentServiceImpl) SendVerificationEmail(ctx context.Context, id string) error {
	user := new(entity.User)
	student := new(entity.Student)
	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.FindByID(ctx, id, user); err != nil {
			return err
		}
		if err := us.sr.FindByID(ctx, user.ID, student); err != nil {
			return err
		}
		if user.Status == constants.VerifiedStatus {
			return customerrors.NewError(
				"user already verified",
				errors.New("user already verified"),
				customerrors.InvalidAction,
			)
		}
		jwt, err := us.ju.GenerateJWT(student.ID, constants.StudentRole, constants.ForVerification, user.Status, constants.AUTH_AGE)
		if err != nil {
			return err
		}
		if err := us.sr.UpdateVerifyToken(ctx, student.ID, &jwt); err != nil {
			return err
		}
		emailParam := entity.SendEmailParams{
			Receiver:  user.Email,
			Subject:   "Verify your account",
			EmailBody: constants.VerificationEmailBody(jwt),
		}
		if err := us.gu.SendEmail(emailParam); err != nil {
			return err
		}
		log.Println(jwt)
		return nil
	})
}

func (us *StudentServiceImpl) ResetPassword(ctx context.Context, param entity.ResetPasswordParam) error {
	user := new(entity.User)
	student := new(entity.Student)

	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.FindByID(ctx, param.ID, user); err != nil {
			return err
		}
		if err := us.sr.FindByID(ctx, user.ID, student); err != nil {
			return err
		}
		if us.bu.ComparePassword(param.NewPassword, user.Password) {
			return customerrors.NewError(
				"cannot update to same password",
				errors.New("old password and new password is the same"),
				customerrors.InvalidAction,
			)
		}
		if student.ResetToken == nil || *student.ResetToken != param.Token {
			return customerrors.NewError("wrong credentials", errors.New("reset token does not match"), customerrors.Unauthenticate)
		}
		newHashedPass, err := us.bu.HashPassword(param.NewPassword)
		if err != nil {
			return err
		}
		if err := us.ur.UpdateUserPassword(ctx, newHashedPass, user.ID); err != nil {
			return err
		}
		if err := us.sr.UpdateResetToken(ctx, student.ID, nil); err != nil {
			return err
		}
		return nil
	})
}

func (us *StudentServiceImpl) SendResetTokenEmail(ctx context.Context, email string) error {
	user := new(entity.User)
	student := new(entity.Student)

	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.FindByEmail(ctx, email, user); err != nil {
			return err
		}
		if err := us.sr.FindByID(ctx, user.ID, student); err != nil {
			return err
		}
		resetToken, err := us.ju.GenerateJWT(student.ID, constants.StudentRole, constants.ForReset, user.Status, constants.RESET_AGE)
		if err != nil {
			return err
		}
		// keep it to test reset password feature
		log.Println(resetToken)
		if err := us.sr.UpdateResetToken(ctx, student.ID, &resetToken); err != nil {
			return err
		}

		param := entity.SendEmailParams{
			Receiver:  user.Email,
			Subject:   "Reset Password",
			EmailBody: constants.ResetEmailBody(resetToken),
		}
		if err := us.gu.SendEmail(param); err != nil {
			return customerrors.NewError("failed to send email", err, customerrors.CommonErr)
		}
		return nil
	})
}

func (us *StudentServiceImpl) Verify(ctx context.Context, param entity.VerifyStudentParam) error {
	user := new(entity.User)
	student := new(entity.Student)
	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.FindByID(ctx, param.ID, user); err != nil {
			return err
		}
		if err := us.sr.FindByID(ctx, user.ID, student); err != nil {
			return err
		}
		if student.VerifyToken == nil || param.Token != *student.VerifyToken {
			if student.VerifyToken != nil {
				log.Println("stored token: ", *student.VerifyToken)
			}
			log.Println("sent token: ", param.Token)
			return customerrors.NewError(
				"invalid credential",
				errors.New("invalid verify token"),
				customerrors.Unauthenticate,
			)
		}
		if err := us.ur.UpdateUserStatus(ctx, constants.VerifiedStatus, user.ID); err != nil {
			return err
		}
		if err := us.sr.UpdateVerifyToken(ctx, student.ID, nil); err != nil {
			return err
		}
		return nil
	})
}

func (us *StudentServiceImpl) Login(ctx context.Context, param entity.StudentLoginParam) (*string, *string, *string, error) {
	user := new(entity.User)
	student := new(entity.Student)
	authToken := new(string)
	refreshToken := new(string)
	status := new(string)

	if err := us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.FindByEmail(ctx, param.Email, user); err != nil {
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
		if err := us.sr.FindByID(ctx, user.ID, student); err != nil {
			return err
		}
		if match := us.bu.ComparePassword(param.Password, user.Password); !match {
			return customerrors.NewError("invalid email or password", errors.New("password does not match"), customerrors.InvalidAction)
		}
		atoken, err := us.ju.GenerateJWT(student.ID, constants.StudentRole, constants.ForLogin, user.Status, constants.AUTH_AGE)
		if err != nil {
			return err
		}
		rtoken, err := us.ju.GenerateJWT(student.ID, constants.StudentRole, constants.ForRefresh, user.Status, constants.REFRESH_AGE)
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

func (us *StudentServiceImpl) Register(ctx context.Context, param entity.StudentRegisterParam) error {
	user := new(entity.User)
	student := new(entity.Student)

	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.FindByEmail(ctx, param.Email, user); err != nil {
			var parsedErr *customerrors.CustomError
			if errors.As(err, &parsedErr) {
				if parsedErr.ErrCode != customerrors.ItemNotExist {
					return err
				}
			}
		} else {
			return customerrors.NewError(
				"user already exist",
				errors.New("user already exist"),
				customerrors.ItemAlreadyExist,
			)
		}

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
		token, err := us.ju.GenerateJWT(user.ID, constants.StudentRole, constants.ForVerification, user.Status, constants.AUTH_AGE)
		if err != nil {
			return err
		}
		// keep it for testing verify functionality
		log.Println(token)

		student.ID = user.ID
		student.VerifyToken = &token
		if err := us.sr.AddNewStudent(ctx, student); err != nil {
			return err
		}

		// wrapped this with go func to make other request does not get blocked when this func running
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
