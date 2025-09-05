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
	mr  *repositories.MentorRepositoryImpl
	tmr *repositories.TransactionManagerRepositories
}

func CreateChatService(chr *repositories.ChatRepositoryImpl, ur *repositories.UserRepositoryImpl, mr *repositories.MentorRepositoryImpl, tmr *repositories.TransactionManagerRepositories) *ChatServiceImpl {
	return &ChatServiceImpl{chr, ur, mr, tmr}
}

func (chs *ChatServiceImpl) CreateChatroom(ctx context.Context, param entity.CreateChatroomParam) (*entity.Chatroom, error) {
	chatroom := new(entity.Chatroom)
	user := new(entity.User)
	mentor := new(entity.Mentor)
	if err := chs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := chs.ur.FindByID(ctx, param.UserID, user); err != nil {
			return err
		}
		if user.Status != constants.VerifiedStatus {
			return customerrors.NewError(
				"user is unverified",
				errors.New("user is unverified"),
				customerrors.Unauthenticate,
			)
		}
		if err := chs.mr.FindByID(ctx, param.MentorID, mentor, false); err != nil {
			return err
		}
		if err := chs.chr.GetChatroom(ctx, param.MentorID, param.UserID, chatroom); err != nil {
			if err.Error() != customerrors.ChatroomNotFound {
				return err
			}
		} else {
			return customerrors.NewError(
				"chatroom already exist",
				errors.New("chatroom already exist"),
				customerrors.ItemAlreadyExist,
			)
		}
		if err := chs.chr.CreateChatroom(ctx, param.MentorID, param.UserID, chatroom); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return chatroom, nil
}

func (chs *ChatServiceImpl) GetChatroom(ctx context.Context, param entity.GetChatroomParam) (*entity.Chatroom, error) {
	chatroom := new(entity.Chatroom)
	user := new(entity.User)
	if err := chs.ur.FindByID(ctx, param.UserID, user); err != nil {
		return nil, err
	}
	if param.Role == constants.StudentRole {
		if user.Status != constants.VerifiedStatus {
			return nil, customerrors.NewError(
				"user is not verified",
				errors.New("user is not verified"),
				customerrors.Unauthenticate,
			)
		}
	}
	if err := chs.chr.FindByID(ctx, param.ChatroomID, chatroom); err != nil {
		return nil, err
	}
	if chatroom.MentorID != param.UserID && chatroom.UserID != param.UserID {
		return nil, customerrors.NewError(
			"unauthorized",
			errors.New("chatroom does not belong to the user"),
			customerrors.Unauthenticate,
		)
	}
	return chatroom, nil
}

func (chs *ChatServiceImpl) GetUserChatrooms(ctx context.Context, param entity.GetUserChatroomsParam) (*[]entity.Chatroom, *int64, error) {
	chatrooms := new([]entity.Chatroom)
	user := new(entity.User)
	totalRow := new(int64)
	if err := chs.ur.FindByID(ctx, param.UserID, user); err != nil {
		return nil, nil, err
	}
	if param.Role == constants.StudentRole {
		if user.Status != constants.VerifiedStatus {
			return nil, nil, customerrors.NewError(
				"user is not verified",
				errors.New("user is not verified"),
				customerrors.Unauthenticate,
			)
		}
	}
	if err := chs.chr.GetUserChatrooms(ctx, param.UserID, param.Limit, param.Page, totalRow, chatrooms); err != nil {
		return nil, nil, err
	}
	return chatrooms, totalRow, nil
}

func (chs *ChatServiceImpl) SendMessage(ctx context.Context, param entity.SendMessageParam) (*entity.Message, error) {
	user := new(entity.User)
	message := new(entity.Message)
	chatroom := new(entity.Chatroom)

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
		if chatroom.UserID != param.UserID && chatroom.MentorID != param.UserID {
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
	return message, nil
}

func (chs *ChatServiceImpl) GetMessages(ctx context.Context, param entity.GetMessagesParam) (*[]entity.Message, *int64, error) {
	user := new(entity.User)
	messages := new([]entity.Message)
	chatroom := new(entity.Chatroom)
	totalRow := new(int64)
	if err := chs.ur.FindByID(ctx, param.UserID, user); err != nil {
		return nil, nil, err
	}
	if param.Role == constants.StudentRole {
		if user.Status != constants.VerifiedStatus {
			return nil, nil, customerrors.NewError(
				"user is not verified",
				errors.New("user is not verified"),
				customerrors.Unauthenticate,
			)
		}
	}
	if err := chs.chr.FindByID(ctx, param.ChatroomID, chatroom); err != nil {
		return nil, nil, customerrors.NewError(
			"unauthorized",
			errors.New("chatroom does not belong to the user"),
			customerrors.Unauthenticate,
		)
	}
	if chatroom.UserID != param.UserID && chatroom.MentorID != param.UserID {
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
