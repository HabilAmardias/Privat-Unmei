package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
	"privat-unmei/internal/utils"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

type StudentServiceImpl struct {
	ur     *repositories.UserRepositoryImpl
	sr     *repositories.StudentRepositoryImpl
	ar     *repositories.AdminRepositoryImpl
	crr    *repositories.CourseRequestRepositoryImpl
	chr    *repositories.ChatRepositoryImpl
	tmr    *repositories.TransactionManagerRepositories
	bu     *utils.BcryptUtil
	gu     *utils.GomailUtil
	cu     *utils.CloudinaryUtil
	ju     *utils.JWTUtil
	ogu    *utils.OTPGenUtil
	goauth *utils.GoogleOauth
}

func CreateStudentService(
	ur *repositories.UserRepositoryImpl,
	sr *repositories.StudentRepositoryImpl,
	ar *repositories.AdminRepositoryImpl,
	crr *repositories.CourseRequestRepositoryImpl,
	chr *repositories.ChatRepositoryImpl,
	tmr *repositories.TransactionManagerRepositories,
	bu *utils.BcryptUtil,
	gu *utils.GomailUtil,
	cu *utils.CloudinaryUtil,
	ju *utils.JWTUtil,
	ogu *utils.OTPGenUtil,
	goauth *utils.GoogleOauth,
) *StudentServiceImpl {
	return &StudentServiceImpl{ur, sr, ar, crr, chr, tmr, bu, gu, cu, ju, ogu, goauth}
}

func (us *StudentServiceImpl) GoogleLogin(oauthState string) string {
	return us.goauth.Config.AuthCodeURL(oauthState)
}

func (us *StudentServiceImpl) DeleteStudent(ctx context.Context, param entity.DeleteStudentParam) error {
	userAdmin := new(entity.User)
	userStudent := new(entity.User)
	admin := new(entity.Admin)
	student := new(entity.Student)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return us.ur.FindByID(ctx, param.AdminID, userAdmin)
	})
	g.Go(func() error {
		return us.ar.FindByID(ctx, param.AdminID, admin)
	})
	g.Go(func() error {
		return us.ur.FindByID(ctx, param.StudentID, userStudent)
	})
	g.Go(func() error {
		return us.sr.FindByID(ctx, param.StudentID, student)
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.chr.DeleteUserMessages(ctx, param.StudentID); err != nil {
			return err
		}
		if err := us.chr.DeleteStudentChatrooms(ctx, param.StudentID); err != nil {
			return err
		}
		if err := us.crr.DeleteAllStudentOrders(ctx, param.StudentID); err != nil {
			return err
		}
		if err := us.sr.DeleteStudent(ctx, param.StudentID); err != nil {
			return err
		}
		return us.ur.DeleteUser(ctx, param.StudentID)
	})
}

func (us *StudentServiceImpl) RefreshToken(ctx context.Context, param entity.RefreshTokenParam) (string, error) {
	user := new(entity.User)
	if err := us.ur.FindByID(ctx, param.UserID, user); err != nil {
		return "", err
	}

	token, err := us.ju.GenerateJWT(param.UserID, param.Role, constants.ForAuth, user.Status, constants.AUTH_AGE)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (us *StudentServiceImpl) GetStudentProfile(ctx context.Context, param entity.StudentProfileParam) (*entity.StudentProfileQuery, error) {

	user := new(entity.User)
	student := new(entity.Student)
	profile := new(entity.StudentProfileQuery)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return us.ur.FindByID(ctx, param.ID, user)
	})
	g.Go(func() error {
		return us.sr.FindByID(ctx, param.ID, student)
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}

	profile.ID = param.ID
	profile.Name = user.Name
	profile.Bio = user.Bio
	profile.ProfileImage = user.ProfileImage
	profile.PublicID = user.PublicID
	profile.Status = user.Status

	return profile, nil
}

func (us *StudentServiceImpl) ChangePassword(ctx context.Context, param entity.StudentChangePasswordParam) error {
	g, ctx := errgroup.WithContext(ctx)
	user := new(entity.User)
	student := new(entity.Student)

	g.Go(func() error {
		return us.ur.FindByID(ctx, param.ID, user)
	})
	g.Go(func() error {
		return us.sr.FindByID(ctx, param.ID, student)
	})
	if err := g.Wait(); err != nil {
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
}

func (us *StudentServiceImpl) GoogleLoginCallback(ctx context.Context, code string) (string, string, string, error) {
	generateAuthAndRefreshToken := func(userID string, status string) (string, string, error) {
		usedFor := constants.ForAuth
		if status != constants.VerifiedStatus {
			usedFor = constants.ForVerification
		}
		authToken, err := us.ju.GenerateJWT(userID, constants.StudentRole, usedFor, status, constants.AUTH_AGE)
		if err != nil {
			return "", "", err
		}
		refreshToken, err := us.ju.GenerateJWT(userID, constants.StudentRole, constants.ForRefresh, status, constants.REFRESH_AGE)
		if err != nil {
			return "", "", err
		}
		return authToken, refreshToken, nil
	}

	var (
		authToken    string
		refreshToken string
	)

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
		if !errors.As(err, &parsedErr) {
			return "", "", "", customerrors.NewError(
				"something went wrong",
				errors.New("fail to parse error"),
				customerrors.CommonErr,
			)
		}
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
				Status:       constants.UnverifiedStatus,
				Password:     hashed,
				ProfileImage: constants.DefaultAvatar,
			}
			if err := us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
				if err := us.ur.AddNewUser(ctx, newUser); err != nil {
					return err
				}
				authToken, refreshToken, err = generateAuthAndRefreshToken(newUser.ID, newUser.Status)
				if err != nil {
					return err
				}
				newStudent := &entity.Student{
					ID:          newUser.ID,
					VerifyToken: &authToken,
				}
				return us.sr.AddNewStudent(ctx, newStudent)
			}); err != nil {
				return "", "", "", err
			}

			return authToken, refreshToken, newUser.Status, nil
		}
		return "", "", "", err
	}
	if err := us.sr.FindByID(ctx, user.ID, student); err != nil {
		return "", "", "", err
	}
	authToken, refreshToken, err = generateAuthAndRefreshToken(user.ID, user.Status)
	if err != nil {
		return "", "", "", err
	}
	if user.Status != constants.VerifiedStatus {
		if err := us.sr.UpdateVerifyToken(ctx, user.ID, &authToken); err != nil {
			return "", "", "", err
		}
	}
	return authToken, refreshToken, user.Status, nil
}

func (us *StudentServiceImpl) UpdateStudentProfile(ctx context.Context, param entity.UpdateStudentParam) error {
	g, ctx := errgroup.WithContext(ctx)
	user := new(entity.User)
	student := new(entity.Student)
	updateQuery := new(entity.UpdateUserQuery)
	g.Go(func() error {
		return us.ur.FindByID(ctx, param.ID, user)
	})
	g.Go(func() error {
		return us.sr.FindByID(ctx, param.ID, student)
	})
	if err := g.Wait(); err != nil {
		return err
	}
	if param.ProfileImage != nil {
		filename := param.ID
		res, err := us.cu.UploadFile(context.Background(), param.ProfileImage, filename, constants.AvatarFolder)
		if err != nil {
			return err
		}
		imageURL := res.SecureURL
		updateQuery.ProfileImage = &imageURL
	}

	updateQuery.Name = param.Name
	updateQuery.Bio = param.Bio

	if err := us.ur.UpdateUserProfile(ctx, updateQuery, param.ID); err != nil {
		return err
	}
	return nil
}

func (us *StudentServiceImpl) GetStudentList(ctx context.Context, param entity.ListStudentParam) (*[]entity.ListStudentQuery, *int64, error) {
	g, ctx := errgroup.WithContext(ctx)

	students := new([]entity.ListStudentQuery)
	totalRow := new(int64)
	user := new(entity.User)
	admin := new(entity.Admin)

	g.Go(func() error {
		if err := us.ur.FindByID(ctx, param.AdminID, user); err != nil {
			return err
		}
		if user.Status == constants.UnverifiedStatus {
			return customerrors.NewError(
				"unauthorized",
				errors.New("admin is not verified"),
				customerrors.Unauthenticate,
			)
		}
		return nil
	})
	g.Go(func() error {
		return us.ar.FindByID(ctx, param.AdminID, admin)
	})
	g.Go(func() error {
		return us.sr.GetStudentList(ctx, totalRow, param.Limit, param.Page, students, param.Search)
	})
	if err := g.Wait(); err != nil {
		return nil, nil, err
	}
	return students, totalRow, nil
}

func (us *StudentServiceImpl) SendVerificationEmail(ctx context.Context, id string) error {
	g, ctx := errgroup.WithContext(ctx)

	user := new(entity.User)
	student := new(entity.Student)

	g.Go(func() error {
		if err := us.ur.FindByID(ctx, id, user); err != nil {
			return err
		}
		if user.Status == constants.VerifiedStatus {
			return customerrors.NewError(
				"user already verified",
				errors.New("user already verified"),
				customerrors.InvalidAction,
			)
		}
		return nil
	})
	g.Go(func() error {
		return us.sr.FindByID(ctx, id, student)
	})
	if err := g.Wait(); err != nil {
		return err
	}
	jwt, err := us.ju.GenerateJWT(id, constants.StudentRole, constants.ForVerification, constants.UnverifiedStatus, constants.AUTH_AGE)
	if err != nil {
		return err
	}
	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
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
		return nil
	})
}

func (us *StudentServiceImpl) ResetPassword(ctx context.Context, param entity.ResetPasswordParam) error {
	g, ctx := errgroup.WithContext(ctx)
	user := new(entity.User)
	student := new(entity.Student)
	g.Go(func() error {
		return us.ur.FindByID(ctx, param.ID, user)
	})
	g.Go(func() error {
		return us.sr.FindByID(ctx, param.ID, student)
	})
	if err := g.Wait(); err != nil {
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
	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
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

func (us *StudentServiceImpl) GoogleVerify(ctx context.Context, param entity.VerifyStudentParam) (string, string, string, error) {
	g, ctx := errgroup.WithContext(ctx)

	user := new(entity.User)
	student := new(entity.Student)

	g.Go(func() error {
		return us.ur.FindByID(ctx, param.ID, user)
	})
	g.Go(func() error {
		return us.sr.FindByID(ctx, param.ID, student)
	})
	if err := g.Wait(); err != nil {
		return "", "", "", err
	}
	if student.VerifyToken == nil || param.Token != *student.VerifyToken {
		return "", "", "", customerrors.NewError(
			"invalid credential",
			errors.New("invalid verify token"),
			customerrors.Unauthenticate,
		)
	}
	if err := us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.UpdateUserStatus(ctx, constants.VerifiedStatus, user.ID); err != nil {
			return err
		}
		if err := us.sr.UpdateVerifyToken(ctx, student.ID, nil); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", "", "", err
	}
	authToken, err := us.ju.GenerateJWT(param.ID, constants.StudentRole, constants.ForAuth, user.Status, constants.AUTH_AGE)
	if err != nil {
		return "", "", "", err
	}
	refreshToken, err := us.ju.GenerateJWT(param.ID, constants.StudentRole, constants.ForRefresh, user.Status, constants.REFRESH_AGE)
	if err != nil {
		return "", "", "", err
	}
	return authToken, refreshToken, user.Status, nil
}

func (us *StudentServiceImpl) Verify(ctx context.Context, param entity.VerifyStudentParam) error {
	g, ctx := errgroup.WithContext(ctx)

	user := new(entity.User)
	student := new(entity.Student)

	g.Go(func() error {
		return us.ur.FindByID(ctx, param.ID, user)
	})
	g.Go(func() error {
		return us.sr.FindByID(ctx, param.ID, student)
	})
	if err := g.Wait(); err != nil {
		return err
	}
	if student.VerifyToken == nil || param.Token != *student.VerifyToken {
		return customerrors.NewError(
			"invalid credential",
			errors.New("invalid verify token"),
			customerrors.Unauthenticate,
		)
	}
	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.UpdateUserStatus(ctx, constants.VerifiedStatus, user.ID); err != nil {
			return err
		}
		if err := us.sr.UpdateVerifyToken(ctx, student.ID, nil); err != nil {
			return err
		}
		return nil
	})
}

func (us *StudentServiceImpl) Login(ctx context.Context, param entity.StudentLoginParam) (string, error) {
	user := new(entity.User)
	student := new(entity.Student)
	loginToken := ""
	if err := us.ur.FindByEmail(ctx, param.Email, user); err != nil {
		var parsedErr *customerrors.CustomError
		if !errors.As(err, &parsedErr) {
			return "", customerrors.NewError(
				"something went wrong",
				errors.New("cannot parse error"),
				customerrors.CommonErr,
			)
		}
		if parsedErr.ErrCode == customerrors.ItemNotExist {
			return "", customerrors.NewError(
				"invalid email or password",
				parsedErr.ErrLog,
				customerrors.InvalidAction,
			)
		}
		return "", err
	}
	if err := us.sr.FindByID(ctx, user.ID, student); err != nil {
		return "", err
	}
	if match := us.bu.ComparePassword(param.Password, user.Password); !match {
		return "", customerrors.NewError(
			"invalid email or password",
			errors.New("password does not match"),
			customerrors.InvalidAction,
		)
	}
	otp, err := us.ogu.GenerateOTP()
	if err != nil {
		return "", err
	}
	loginToken, err = us.ju.GenerateJWT(user.ID, constants.StudentRole, constants.ForLogin, user.Status, constants.LOGIN_AGE)
	if err != nil {
		return "", err
	}
	now := time.Now()

	if err := us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.sr.UpdateOTP(ctx, user.ID, &now, &otp); err != nil {
			return err
		}
		if err := us.sr.UpdateLoginToken(ctx, user.ID, &loginToken); err != nil {
			return err
		}
		return nil

	}); err != nil {
		return "", err
	}
	go func() {
		if err := us.gu.SendEmail(entity.SendEmailParams{
			Receiver:  user.Email,
			Subject:   "One Time Password for Login - Privat Unmei",
			EmailBody: constants.OTPEmailBody(otp),
		}); err != nil {
			log.Println(err)
		}
	}()
	return loginToken, nil
}

func (us *StudentServiceImpl) LoginCallback(ctx context.Context, param entity.LoginCallbackParam) (string, string, string, error) {
	user := new(entity.User)
	student := new(entity.Student)
	now := time.Now()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return us.ur.FindByID(ctx, param.UserID, user)
	})
	g.Go(func() error {
		return us.sr.FindByID(ctx, param.UserID, student)
	})
	if err := g.Wait(); err != nil {
		return "", "", "", err
	}
	if student.OTP == nil || *student.OTP != param.OTP {
		return "", "", "", customerrors.NewError(
			"invalid code",
			errors.New("invalid OTP"),
			customerrors.Unauthenticate,
		)
	}
	if student.LoginToken == nil || *student.LoginToken != param.LoginToken {
		return "", "", "", customerrors.NewError(
			"unauthorized",
			errors.New("invalid token"),
			customerrors.Unauthenticate,
		)
	}
	if student.OTPLastUpdatedAt.Add(constants.OTP_DURATION).Before(now) {
		return "", "", "", customerrors.NewError(
			"Code Expired",
			fmt.Errorf("OTP Expired"),
			customerrors.InvalidAction,
		)
	}
	if err := us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.sr.UpdateLoginToken(ctx, param.UserID, nil); err != nil {
			return err
		}
		return us.sr.UpdateOTP(ctx, param.UserID, nil, nil)
	}); err != nil {
		return "", "", "", err
	}
	authToken, err := us.ju.GenerateJWT(user.ID, constants.StudentRole, constants.ForAuth, user.Status, constants.AUTH_AGE)
	if err != nil {
		return "", "", "", err
	}
	refreshToken, err := us.ju.GenerateJWT(user.ID, constants.StudentRole, constants.ForRefresh, user.Status, constants.REFRESH_AGE)
	if err != nil {
		return "", "", "", err
	}
	return authToken, refreshToken, user.Status, nil
}

func (us *StudentServiceImpl) ResendOTP(ctx context.Context, param entity.ResendOTPParam) (string, error) {
	user := new(entity.User)
	student := new(entity.Student)
	now := time.Now()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return us.ur.FindByID(ctx, param.UserID, user)
	})

	g.Go(func() error {
		return us.sr.FindByID(ctx, param.UserID, student)
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	if student.OTPLastUpdatedAt == nil {
		return "", customerrors.NewError(
			"unauthorized",
			errors.New("no otp last updated at"),
			customerrors.Unauthenticate,
		)
	}

	if student.OTPLastUpdatedAt.Add(constants.OTP_DURATION / 2).After(now) {
		return "", customerrors.NewError(
			"please try again later",
			fmt.Errorf("need to wait to send otp again"),
			customerrors.InvalidAction,
		)
	}

	otp, err := us.ogu.GenerateOTP()
	if err != nil {
		return "", err
	}
	loginToken, err := us.ju.GenerateJWT(param.UserID, constants.StudentRole, constants.ForLogin, user.Status, constants.LOGIN_AGE)
	if err != nil {
		return "", err
	}

	if err := us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.sr.UpdateLoginToken(ctx, param.UserID, &loginToken); err != nil {
			return err
		}
		now := time.Now()
		return us.sr.UpdateOTP(ctx, param.UserID, &now, &otp)
	}); err != nil {
		return "", err
	}
	go func() {
		if err := us.gu.SendEmail(entity.SendEmailParams{
			Receiver:  user.Email,
			Subject:   "One Time Password for Login - Privat Unmei",
			EmailBody: constants.OTPEmailBody(otp),
		}); err != nil {
			log.Println(err)
		}
	}()
	return loginToken, nil
}

func (us *StudentServiceImpl) Register(ctx context.Context, param entity.StudentRegisterParam) error {
	user := new(entity.User)
	student := new(entity.Student)

	return us.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := us.ur.FindByEmail(ctx, param.Email, user); err != nil {
			var parsedErr *customerrors.CustomError
			if !errors.As(err, &parsedErr) {
				return customerrors.NewError(
					"something went wrong",
					errors.New("cannot parse error"),
					customerrors.CommonErr,
				)
			}
			if parsedErr.ErrCode != customerrors.ItemNotExist {
				return err
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
