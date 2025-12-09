package services

import (
	"context"
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"
)

type ChatServiceImpl struct {
	chr *repositories.ChatRepositoryImpl
	ur  *repositories.UserRepositoryImpl
	sr  *repositories.StudentRepositoryImpl
	mr  *repositories.MentorRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func CreateChatService(chr *repositories.ChatRepositoryImpl, ur *repositories.UserRepositoryImpl, sr *repositories.StudentRepositoryImpl, mr *repositories.MentorRepositoryImpl, tmr *repositories.TransactionManagerRepositories) *ChatServiceImpl {
	return &ChatServiceImpl{chr, ur, sr, mr, tmr}
}

func (chs *ChatServiceImpl) GetChatroom(ctx context.Context, param entity.GetChatroomParam) (int, error) {
	chatroom := new(entity.Chatroom)
	user := new(entity.User)
	mentor := new(entity.Mentor)
	student := new(entity.Student)

	if err := chs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := chs.ur.FindByID(ctx, param.UserID, user); err != nil {
			return err
		}
		if user.Status != constants.VerifiedStatus {
			return customerrors.NewError(
				"user is not verified",
				errors.New("user is not verified"),
				customerrors.Unauthenticate,
			)
		}
		// case: logged in user is a mentor
		if param.Role == constants.MentorRole {
			secondUserStudent := new(entity.Student)
			if err := chs.mr.FindByID(ctx, param.UserID, mentor, false); err != nil {
				return err
			}
			if err := chs.sr.FindByID(ctx, param.SecondUserID, secondUserStudent); err != nil {
				return err
			}
			if err := chs.chr.GetChatroom(ctx, param.UserID, param.SecondUserID, chatroom); err != nil {
				var parsedErr *customerrors.CustomError
				if errors.As(err, &parsedErr) {
					if parsedErr.ErrCode != customerrors.ItemNotExist {
						return err
					}
				}
				if err := chs.chr.CreateChatroom(ctx, param.UserID, param.SecondUserID, chatroom); err != nil {
					return err
				}
			}
			return nil
		}
		// case: logged in user is a student
		if param.Role == constants.StudentRole {
			secondUserMentor := new(entity.Mentor)
			if err := chs.sr.FindByID(ctx, param.UserID, student); err != nil {
				return err
			}
			if err := chs.mr.FindByID(ctx, param.SecondUserID, secondUserMentor, false); err != nil {
				return err
			}
			if err := chs.chr.GetChatroom(ctx, param.SecondUserID, param.UserID, chatroom); err != nil {
				var parsedErr *customerrors.CustomError
				if errors.As(err, &parsedErr) {
					if parsedErr.ErrCode != customerrors.ItemNotExist {
						return err
					}
				}
				if err := chs.chr.CreateChatroom(ctx, param.SecondUserID, param.UserID, chatroom); err != nil {
					return err
				}
			}
			return nil
		}
		// default case: logged in user is neither a mentor nor a student
		return customerrors.NewError(
			"invalid user",
			errors.New("role does not belong to student or mentor"),
			customerrors.Unauthenticate,
		)
	}); err != nil {
		return 0, err
	}
	return chatroom.ID, nil
}

func (chs *ChatServiceImpl) GetChatroomInfo(ctx context.Context, param entity.GetChatroomInfoParam) (*entity.ChatroomDetailQuery, error) {
	chatroom := new(entity.Chatroom)
	user := new(entity.User)
	secondUser := new(entity.User)
	query := new(entity.ChatroomDetailQuery)

	if err := chs.ur.FindByID(ctx, param.UserID, user); err != nil {
		return nil, err
	}
	if user.Status != constants.VerifiedStatus {
		return nil, customerrors.NewError(
			"user is not verified",
			errors.New("user is not verified"),
			customerrors.Unauthenticate,
		)
	}
	if err := chs.chr.FindByID(ctx, param.ChatroomID, chatroom); err != nil {
		return nil, err
	}
	if chatroom.MentorID != param.UserID && chatroom.StudentID != param.UserID {
		return nil, customerrors.NewError(
			"unauthorized",
			errors.New("chatroom does not belong to the user"),
			customerrors.Unauthenticate,
		)
	}
	secondUserID := chatroom.MentorID
	if param.Role == constants.MentorRole {
		secondUserID = chatroom.StudentID
	}
	if err := chs.ur.FindByID(ctx, secondUserID, secondUser); err != nil {
		return nil, err
	}

	query.ID = chatroom.ID
	query.UserID = secondUser.ID
	query.Username = secondUser.Name
	query.UserEmail = secondUser.Email
	query.UserProfileImage = secondUser.ProfileImage

	return query, nil
}

func (chs *ChatServiceImpl) GetUserChatrooms(ctx context.Context, param entity.GetUserChatroomsParam) (*[]entity.ChatroomDetailQuery, *int64, error) {
	chatrooms := new([]entity.ChatroomDetailQuery)
	user := new(entity.User)
	totalRow := new(int64)
	if err := chs.ur.FindByID(ctx, param.UserID, user); err != nil {
		return nil, nil, err
	}
	if user.Status != constants.VerifiedStatus {
		return nil, nil, customerrors.NewError(
			"user is not verified",
			errors.New("user is not verified"),
			customerrors.Unauthenticate,
		)
	}
	if err := chs.chr.GetUserChatrooms(ctx, param.UserID, param.Role, param.Limit, param.Page, totalRow, chatrooms); err != nil {
		return nil, nil, err
	}
	return chatrooms, totalRow, nil
}

func (chs *ChatServiceImpl) SendMessage(ctx context.Context, param entity.SendMessageParam) (*entity.MessageDetailQuery, error) {
	user := new(entity.User)
	message := new(entity.Message)
	chatroom := new(entity.Chatroom)
	messageDetail := new(entity.MessageDetailQuery)

	if err := chs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := chs.ur.FindByID(ctx, param.UserID, user); err != nil {
			return err
		}
		if user.Status != constants.VerifiedStatus {
			return customerrors.NewError(
				"user is not verified",
				errors.New("user is not verified"),
				customerrors.Unauthenticate,
			)
		}
		if err := chs.chr.FindByID(ctx, param.ChatroomID, chatroom); err != nil {
			return err
		}
		if chatroom.StudentID != param.UserID && chatroom.MentorID != param.UserID {
			return customerrors.NewError(
				"unauthorized",
				errors.New("chatroom does not belong to the user"),
				customerrors.Unauthenticate,
			)
		}
		if err := chs.chr.SendMessage(ctx, param.UserID, param.ChatroomID, param.Content, message); err != nil {
			return err
		}
		if err := chs.chr.UpdateChatroom(ctx, param.ChatroomID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	messageDetail.ID = message.ID
	messageDetail.ChatroomID = message.ChatroomID
	messageDetail.SenderEmail = user.Email
	messageDetail.SenderID = user.ID
	messageDetail.SenderName = user.Name
	messageDetail.Content = message.Content

	return messageDetail, nil
}

func (chs *ChatServiceImpl) GetMessages(ctx context.Context, param entity.GetMessagesParam) (*[]entity.MessageDetailQuery, *int64, error) {
	user := new(entity.User)
	messages := new([]entity.MessageDetailQuery)
	chatroom := new(entity.Chatroom)
	totalRow := new(int64)
	if err := chs.ur.FindByID(ctx, param.UserID, user); err != nil {
		return nil, nil, err
	}
	if user.Status != constants.VerifiedStatus {
		return nil, nil, customerrors.NewError(
			"user is not verified",
			errors.New("user is not verified"),
			customerrors.Unauthenticate,
		)
	}
	if err := chs.chr.FindByID(ctx, param.ChatroomID, chatroom); err != nil {
		return nil, nil, customerrors.NewError(
			"unauthorized",
			errors.New("chatroom does not belong to the user"),
			customerrors.Unauthenticate,
		)
	}
	if chatroom.StudentID != param.UserID && chatroom.MentorID != param.UserID {
		return nil, nil, customerrors.NewError(
			"unauthorized",
			errors.New("chatroom does not belong to the user"),
			customerrors.Unauthenticate,
		)
	}

	if err := chs.chr.GetMessages(ctx, param.ChatroomID, param.Limit, param.LastID, totalRow, messages); err != nil {
		return nil, nil, err
	}
	return messages, totalRow, nil
}
