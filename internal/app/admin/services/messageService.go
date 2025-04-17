package services

import (
	"encoding/json"
	"strconv"
	"time"

	"server/internal/app/admin/dto"
	"server/internal/constant"
	"server/internal/event"
	"server/internal/global"
	"server/internal/models"
	"server/internal/repositories"

	"gorm.io/gorm"
)

func init() {
	event.AdminSubscribe(func(event event.Event) {
		msgService := NewMessageService()
		userRepo := repositories.NewUserRepo()
		userIDs, err := userRepo.GetAllAdminID()
		if err != nil {
			global.Logger.Errorw("get all admin id error", "error", err)
		}
		if userIDs != nil {
			msgService.SendMsg(*event.Content, *userIDs)
		}
	})
}

type MessageService struct {
	msgRepo  *repositories.MessageRepo
	userRepo *repositories.UserRepo
	resource *repositories.ResourceRepo
}

var messageService *MessageService

func NewMessageService() *MessageService {
	if messageService == nil {
		messageService = &MessageService{
			msgRepo:  repositories.NewMessageRepo(),
			userRepo: repositories.NewUserRepo(),
			resource: repositories.NewResourceRepo(),
		}
	}
	return messageService
}

func (m *MessageService) SendMsg(content string, userIDs []uint) string {
	msgID := strconv.Itoa(int(time.Now().UnixNano()))
	msg := models.Message{
		ID:        msgID,
		To:        userIDs,
		Content:   content,
		CreatedAt: time.Now().Local().Format(time.DateTime),
	}

	if err := m.msgRepo.SaveAdminMsg(msg, constant.ADMIN_NOTIFICATION); err != nil {
		global.Logger.Errorw("save message error", "error", err)
	}

	for _, userID := range userIDs {
		err := m.msgRepo.AddUnReadMsg(strconv.Itoa(int(userID)), msgID, constant.ADMIN_MESSAGE_UNREADED)
		if err != nil {
			global.Logger.Errorw("add message error", "error", err)
		}

		if err := m.msgRepo.PublishMsg(content, userID, constant.ADMIN_MESSAGE_CHANNEL); err != nil {
			global.Logger.Errorw("publish message error", "error", err)
		}
	}

	return msgID
}

func (m *MessageService) GetUnReadMsgs(userID uint) []models.Message {
	return m.msgRepo.GetAdminUnReadMsgs(strconv.Itoa(int(userID)))
}

func (m *MessageService) GetReadedMsgs(userID uint) []models.Message {
	return m.msgRepo.GetAdminReadedMsgs(strconv.Itoa(int(userID)))
}

func (m *MessageService) MarkReadMsg(userID uint, msgID string) error {
	if err := m.msgRepo.MarkAdminMsgAsRead(strconv.Itoa(int(userID)), msgID); err != nil {
		return err
	}

	if err := m.msgRepo.PublishMsg("", userID, constant.ADMIN_MESSAGE_CHANNEL); err != nil {
		return err
	}

	return nil
}

func (m *MessageService) DeleteMsg(msgID string) error {
	userIDs := []uint{}
	result := m.msgRepo.GetMsgField(msgID, "to", constant.ADMIN_NOTIFICATION)
	if err := json.Unmarshal([]byte(result), &userIDs); err != nil {
		return err
	}

	if err := m.msgRepo.DeleteMsg(msgID, constant.ADMIN_NOTIFICATION); err != nil {
		return err
	}

	for _, userID := range userIDs {
		if err := m.msgRepo.RemoveAdminMsg(msgID, strconv.Itoa(int(userID))); err != nil {
			return err
		}

		if err := m.msgRepo.PublishMsg("", userID, constant.ADMIN_MESSAGE_CHANNEL); err != nil {
			return err
		}
	}

	return nil
}

func (m *MessageService) GetAllMsgs() ([]dto.MsgGetAllResponse, error) {
	messages, err := m.msgRepo.GetAllAdminMsgs(10)
	if err != nil {
		return nil, err
	}

	response := []dto.MsgGetAllResponse{}
	for _, message := range messages {
		userIDs := message.To
		to := []models.Member{}
		for _, userID := range userIDs {
			user, err := m.userRepo.GetUserById(userID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					continue
				}
				return nil, err
			}

			resource, err := m.resource.GetResourceById(user.Avatar)
			if err != nil {
				return nil, err
			}

			to = append(to, models.Member{
				UserID:   user.ID,
				Username: user.Username,
				Avatar:   resource.StaticPath,
			})
		}

		response = append(response, dto.MsgGetAllResponse{
			ID:        message.ID,
			To:        to,
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
		})
	}

	return response, nil
}
