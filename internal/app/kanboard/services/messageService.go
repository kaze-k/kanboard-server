package services

import (
	"strconv"
	"time"

	"server/internal/app/kanboard/dto"
	"server/internal/constant"
	"server/internal/event"
	"server/internal/global"
	"server/internal/models"
	"server/internal/repositories"
	"server/internal/utils"
)

func init() {
	event.KanboardSubscribe(func(event event.Event) {
		if event.UserID != nil {
			msgRepo := repositories.NewMessageRepo()
			msgRepo.PublishMsg(constant.UPDATE_USER, *event.UserID, constant.KANBOARD_MESSAGE_CHANNEL)
		}

		if event.EventType == nil && event.UserID != nil && event.ProjectID != nil && event.Content != nil {
			msgRepo := repositories.NewMessageRepo()
			msgRepo.PublishMsg(*event.Content, *event.UserID, constant.KANBOARD_MESSAGE_CHANNEL)
		}

		if event.EventType == nil {
			return
		}

		msgService := NewMessageService()

		projectMemberRepo := repositories.NewProjectMemberRepo()
		memberIds, err := projectMemberRepo.GetAllMemberIdByProjectId(*event.ProjectID)
		if err != nil {
			global.Logger.Errorw("get all member id error", "error", err)
		}

		project_event := constant.PROJECT_EVENT
		task_event := constant.TASK_EVENT
		if event.EventType != nil && *event.EventType == project_event {
			if memberIds != nil {
				msgService.SendMsg(*event.Content, *memberIds, 0, *event.ProjectID, "")
			}
		} else if event.EventType != nil && *event.EventType == task_event {
			taskAssigneeRepo := repositories.NewTaskAssigneeRepo()
			assigneeIds, err := taskAssigneeRepo.GetAllAssigneeIdByTaskId(*event.TaskID)
			if err != nil {
				global.Logger.Errorw("get all assignee id error", "error", err)
			}

			taskRepo := repositories.NewTaskRepo()
			creatorId, err := taskRepo.GetCreatorIdByTaskId(*event.TaskID)
			if err != nil {
				global.Logger.Errorw("get creator id error", "error", err)
			}

			merged := append(*assigneeIds, *memberIds...)
			merged = append(merged, *creatorId)
			unique := utils.UniqueUintSlice(merged)

			if len(unique) > 0 {
				msgService.SendMsg(*event.Content, unique, *event.TaskID, *event.ProjectID, constant.NEW_TASK_STATUS)
			}
		}
	})
}

type MessageService struct {
	msgRepo           *repositories.MessageRepo
	userRepo          *repositories.UserRepo
	resource          *repositories.ResourceRepo
	project           *repositories.ProjectRepo
	projectMemberRepo *repositories.ProjectMemberRepo
}

var messageService *MessageService

func NewMessageService() *MessageService {
	if messageService == nil {
		messageService = &MessageService{
			msgRepo:           repositories.NewMessageRepo(),
			userRepo:          repositories.NewUserRepo(),
			resource:          repositories.NewResourceRepo(),
			project:           repositories.NewProjectRepo(),
			projectMemberRepo: repositories.NewProjectMemberRepo(),
		}
	}
	return messageService
}

func (m *MessageService) SendMsg(content string, userIDs []uint, taskId uint, projectId uint, publish string) string {
	msgID := strconv.Itoa(int(time.Now().UnixNano()))
	msg := models.Message{
		ID:        msgID,
		To:        userIDs,
		ProjectID: &projectId,
		TaskID:    &taskId,
		Content:   content,
		CreatedAt: time.Now().Local().Format(time.DateTime),
	}

	if err := m.msgRepo.SaveKanboardMsg(msg, constant.KANBOARD_NOTIFICATION, taskId, projectId); err != nil {
		global.Logger.Errorw("save message error", "error", err)
	}

	for _, userID := range userIDs {
		if err := m.msgRepo.AddProjectMsg(strconv.Itoa(int(projectId)), msgID, constant.KANBOARD_MESSAGE_UNREADED); err != nil {
			global.Logger.Errorw("add message error", "error", err)
		}

		err := m.msgRepo.AddUnReadMsg(strconv.Itoa(int(userID)), msgID, constant.KANBOARD_MESSAGE_UNREADED)
		if err != nil {
			global.Logger.Errorw("add message error", "error", err)
		}

		if err := m.msgRepo.PublishMsg(content, userID, constant.KANBOARD_MESSAGE_CHANNEL); err != nil {
			global.Logger.Errorw("publish message error", "error", err)
		}

		if publish != "" {
			if err := m.msgRepo.PublishMsg(publish, userID, constant.KANBOARD_MESSAGE_CHANNEL); err != nil {
				global.Logger.Errorw("publish message error", "error", err)
			}
		}
	}

	return msgID
}

func (m *MessageService) GetUnReadMsgs(userID uint) []models.Message {
	projectMembers, err := m.projectMemberRepo.GetProjectByUserId(userID)
	if err != nil {
		return nil
	}

	messages := []models.Message{}
	for _, projectMember := range projectMembers {
		msgs := m.msgRepo.GetKanboardUnReadMsgs(strconv.Itoa(int(userID)), projectMember.ProjectID)
		messages = append(messages, msgs...)
	}

	return messages
}

func (m *MessageService) GetReadedMsgs(userID uint) []models.Message {
	projectMembers, err := m.projectMemberRepo.GetProjectByUserId(userID)
	if err != nil {
		return nil
	}

	messages := []models.Message{}
	for _, projectMember := range projectMembers {
		msgs := m.msgRepo.GetKanboardReadedMsgs(strconv.Itoa(int(userID)), projectMember.ProjectID)
		messages = append(messages, msgs...)
	}

	return messages
}

func (m *MessageService) MarkReadMsg(userID uint, msgID string) error {
	if err := m.msgRepo.MarkKanboardMsgAsRead(strconv.Itoa(int(userID)), msgID); err != nil {
		return err
	}

	if err := m.msgRepo.PublishMsg(constant.PUBLISH_MESSAGE, userID, constant.KANBOARD_MESSAGE_CHANNEL); err != nil {
		return err
	}

	return nil
}

func (m *MessageService) GetMsgsByProjectId(projectId uint) ([]dto.MsgResponse, error) {
	messages, err := m.msgRepo.GetMsgsByProjectId(10, projectId)
	if err != nil {
		return nil, err
	}

	response := []dto.MsgResponse{}
	for _, message := range messages {
		response = append(response, dto.MsgResponse{
			ID:        message.ID,
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
		})
	}

	return response, nil
}
