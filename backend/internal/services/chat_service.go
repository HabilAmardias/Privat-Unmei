package services

import (
	"context"
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/repositories"

	"golang.org/x/sync/errgroup"
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

func (chs *ChatServiceImpl) GetChatroom(ctx context.Context, param entity.GetChatroomParam) (string, error) {
	chatroom := new(entity.Chatroom)
	user := new(entity.User)
	mentor := new(entity.Mentor)
	student := new(entity.Student)
	secondUser := new(entity.User)

	if param.Role != constants.MentorRole && param.Role != constants.StudentRole {
		return "", customerrors.NewError(
			"invalid user",
			errors.New("role does not belong to student or mentor"),
			customerrors.Unauthenticate,
		)
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
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
		return nil
	})

	g.Go(func() error {
		if err := chs.ur.FindByID(ctx, param.SecondUserID, secondUser); err != nil {
			return err
		}
		if secondUser.Status != constants.VerifiedStatus {
			return customerrors.NewError(
				"user is not verified",
				errors.New("user is not verified"),
				customerrors.Unauthenticate,
			)
		}
		return nil
	})

	if param.Role == constants.MentorRole {
		secondUserStudent := new(entity.Student)
		g.Go(func() error {
			return chs.mr.FindByID(ctx, param.UserID, mentor, false)
		})
		g.Go(func() error {
			return chs.sr.FindByID(ctx, param.SecondUserID, secondUserStudent)
		})
		g.Go(func() error {
			if err := chs.chr.GetChatroom(ctx, param.UserID, param.SecondUserID, chatroom); err != nil {
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
				if err := chs.chr.CreateChatroom(ctx, param.UserID, param.SecondUserID, chatroom); err != nil {
					return err
				}
			}
			return nil
		})
	} else {
		secondUserMentor := new(entity.Mentor)
		g.Go(func() error {
			return chs.sr.FindByID(ctx, param.UserID, student)
		})
		g.Go(func() error {
			return chs.mr.FindByID(ctx, param.SecondUserID, secondUserMentor, false)
		})
		g.Go(func() error {
			if err := chs.chr.GetChatroom(ctx, param.SecondUserID, param.UserID, chatroom); err != nil {
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
				if err := chs.chr.CreateChatroom(ctx, param.SecondUserID, param.UserID, chatroom); err != nil {
					return err
				}
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return "", err
	}
	return chatroom.ID, nil
}

func (chs *ChatServiceImpl) GetChatroomInfo(ctx context.Context, param entity.GetChatroomInfoParam) (*entity.ChatroomDetailQuery, error) {
	chatroom := new(entity.Chatroom)
	user := new(entity.User)
	secondUser := new(entity.User)
	query := new(entity.ChatroomDetailQuery)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
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
		return nil
	})
	g.Go(func() error {
		if err := chs.chr.FindByID(ctx, param.ChatroomID, chatroom); err != nil {
			return err
		}
		if chatroom.MentorID != param.UserID && chatroom.StudentID != param.UserID {
			return customerrors.NewError(
				"unauthorized",
				errors.New("chatroom does not belong to the user"),
				customerrors.Unauthenticate,
			)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return nil, err
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
	query.UserPublicID = secondUser.PublicID
	query.UserProfileImage = secondUser.ProfileImage

	return query, nil
}

func (chs *ChatServiceImpl) GetUserChatrooms(ctx context.Context, param entity.GetUserChatroomsParam) (*[]entity.ChatroomDetailQuery, *int64, error) {
	chatrooms := new([]entity.ChatroomDetailQuery)
	user := new(entity.User)
	totalRow := new(int64)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
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
		return nil
	})
	g.Go(func() error {
		return chs.chr.GetUserChatrooms(ctx, param.UserID, param.Role, param.Limit, param.Page, totalRow, chatrooms)
	})
	if err := g.Wait(); err != nil {
		return nil, nil, err
	}
	return chatrooms, totalRow, nil
}

func (chs *ChatServiceImpl) SendMessage(ctx context.Context, param entity.SendMessageParam) (*entity.MessageDetailQuery, error) {
	user := new(entity.User)
	message := new(entity.Message)
	chatroom := new(entity.Chatroom)
	messageDetail := new(entity.MessageDetailQuery)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
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
		return nil
	})
	g.Go(func() error {
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
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if err := chs.tmr.WithTransaction(ctx, func(ctx context.Context) error {
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
	messageDetail.SenderPublicID = user.PublicID
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

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
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
		return nil
	})
	g.Go(func() error {
		if err := chs.chr.FindByID(ctx, param.ChatroomID, chatroom); err != nil {
			return customerrors.NewError(
				"unauthorized",
				errors.New("chatroom does not belong to the user"),
				customerrors.Unauthenticate,
			)
		}
		if chatroom.StudentID != param.UserID && chatroom.MentorID != param.UserID {
			return customerrors.NewError(
				"unauthorized",
				errors.New("chatroom does not belong to the user"),
				customerrors.Unauthenticate,
			)
		}
		return nil
	})
	g.Go(func() error {
		return chs.chr.GetMessages(ctx, param.ChatroomID, param.Limit, param.LastID, totalRow, messages)
	})
	if err := g.Wait(); err != nil {
		return nil, nil, err
	}
	return messages, totalRow, nil
}
