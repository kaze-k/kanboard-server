package repositories

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"server/internal/constant"
	"server/internal/global"
	"server/internal/models"

	"github.com/redis/go-redis/v9"
)

type MessageRepo struct {
	redis *global.RedisClient
}

var messageRepo *MessageRepo

func NewMessageRepo() *MessageRepo {
	if messageRepo == nil {
		messageRepo = &MessageRepo{
			redis: global.Redis,
		}
	}
	return messageRepo
}

func (m *MessageRepo) SaveAdminMsg(msg models.Message, namespace string) error {
	jsonData, err := json.Marshal(msg.To)
	if err != nil {
		return err
	}
	data := map[string]string{
		"id":        msg.ID,
		"to":        string(jsonData),
		"content":   msg.Content,
		"createdAt": msg.CreatedAt,
	}
	if err := m.redis.HSet(namespace, msg.ID, data); err != nil {
		return err
	}
	return nil
}

func (m *MessageRepo) SaveKanboardMsg(msg models.Message, namespace string, taskId uint, projectId uint) error {
	jsonData, err := json.Marshal(msg.To)
	if err != nil {
		return err
	}
	data := map[string]string{
		"id":        msg.ID,
		"to":        string(jsonData),
		"projectID": fmt.Sprintf("%d", projectId),
		"taskID":    fmt.Sprintf("%d", taskId),
		"content":   msg.Content,
		"createdAt": msg.CreatedAt,
	}
	namespace = fmt.Sprintf("%s/%d", namespace, projectId)
	if err := m.redis.HSet(namespace, msg.ID, data); err != nil {
		return err
	}
	return nil
}

func (m *MessageRepo) AddProjectMsg(projectId, msgId string, namespace string) error {
	return m.redis.SAdd(namespace, projectId, msgId, true)
}

func (m *MessageRepo) AddUnReadMsg(userID, msgeID string, namespace string) error {
	return m.redis.SAdd(namespace, userID, msgeID, true)
}

func (m *MessageRepo) GetUnReadCount(userID string, namespace string) int64 {
	return m.redis.SCard(namespace, userID)
}

func (m *MessageRepo) GetAdminUnReadMsgs(userID string) []models.Message {
	msgIDs := m.redis.SMembers(constant.ADMIN_MESSAGE_UNREADED, userID)

	messages := []models.Message{}
	for _, msgID := range msgIDs {
		data := m.redis.HGetAll(constant.ADMIN_NOTIFICATION, msgID)
		to := []uint{}
		if err := json.Unmarshal([]byte(data["to"]), &to); err != nil {
			continue
		}
		if len(data) > 0 {
			messages = append(messages, models.Message{
				ID:        data["id"],
				To:        to,
				Content:   data["content"],
				CreatedAt: data["createdAt"],
			})
		}
	}

	sort.Slice(messages, func(i, j int) bool {
		t1 := messages[i].CreatedAt
		t2 := messages[j].CreatedAt
		return t1 > t2
	})

	return messages
}

func (m *MessageRepo) GetKanboardUnReadMsgs(userID string, projectId uint) []models.Message {
	msgIDs := m.redis.SMembers(constant.KANBOARD_MESSAGE_UNREADED, userID)

	messages := []models.Message{}
	for _, msgID := range msgIDs {
		namespace := fmt.Sprintf("%s/%d", constant.KANBOARD_NOTIFICATION, projectId)
		data := m.redis.HGetAll(namespace, msgID)
		to := []uint{}
		if err := json.Unmarshal([]byte(data["to"]), &to); err != nil {
			continue
		}
		var taskID uint
		if data["taskID"] != "" {
			id, err := strconv.Atoi(data["taskID"])
			if err != nil {
				continue
			}
			taskID = uint(id)
		}
		var projectID uint
		if data["projectID"] != "" {
			id, err := strconv.Atoi(data["projectID"])
			if err != nil {
				projectID = 0
			}
			projectID = uint(id)
		}
		if len(data) > 0 {
			messages = append(messages, models.Message{
				ID:        data["id"],
				To:        to,
				ProjectID: &projectID,
				TaskID:    &taskID,
				Content:   data["content"],
				CreatedAt: data["createdAt"],
			})
		}
	}

	sort.Slice(messages, func(i, j int) bool {
		t1 := messages[i].CreatedAt
		t2 := messages[j].CreatedAt
		return t1 > t2
	})

	return messages
}

func (m *MessageRepo) GetAdminReadedMsgs(userID string) []models.Message {
	msgIDs := m.redis.SMembers(constant.ADMIN_MESSAGE_READED, userID)

	messages := []models.Message{}
	for _, msgID := range msgIDs {
		data := m.redis.HGetAll(constant.ADMIN_NOTIFICATION, msgID)
		to := []uint{}
		if err := json.Unmarshal([]byte(data["to"]), &to); err != nil {
			continue
		}
		if len(data) > 0 {
			messages = append(messages, models.Message{
				ID:        data["id"],
				To:        to,
				Content:   data["content"],
				CreatedAt: data["createdAt"],
			})
		}
	}

	sort.Slice(messages, func(i, j int) bool {
		t1 := messages[i].CreatedAt
		t2 := messages[j].CreatedAt
		return t1 > t2
	})

	return messages
}

func (m *MessageRepo) GetKanboardReadedMsgs(userID string, projectId uint) []models.Message {
	msgIDs := m.redis.SMembers(constant.KANBOARD_MESSAGE_READED, userID)

	messages := []models.Message{}
	for _, msgID := range msgIDs {
		namespace := fmt.Sprintf("%s/%d", constant.KANBOARD_NOTIFICATION, projectId)
		data := m.redis.HGetAll(namespace, msgID)
		to := []uint{}
		if err := json.Unmarshal([]byte(data["to"]), &to); err != nil {
			continue
		}
		var taskID uint
		if data["taskID"] != "" {
			id, err := strconv.Atoi(data["taskID"])
			if err != nil {
				continue
			}
			taskID = uint(id)
		}
		var projectID uint
		if data["projectID"] != "" {
			id, err := strconv.Atoi(data["projectID"])
			if err != nil {
				projectID = 0
			}
			projectID = uint(id)
		}
		if len(data) > 0 {
			messages = append(messages, models.Message{
				ID:        data["id"],
				To:        to,
				ProjectID: &projectID,
				TaskID:    &taskID,
				Content:   data["content"],
				CreatedAt: data["createdAt"],
			})
		}
	}

	sort.Slice(messages, func(i, j int) bool {
		t1 := messages[i].CreatedAt
		t2 := messages[j].CreatedAt
		return t1 > t2
	})

	return messages
}

func (m *MessageRepo) MarkAdminMsgAsRead(userID, msgID string) error {
	if err := m.redis.SRem(constant.ADMIN_MESSAGE_UNREADED, userID, msgID); err != nil {
		return err
	}
	if err := m.redis.SAdd(constant.ADMIN_MESSAGE_READED, userID, msgID, false); err != nil {
		return err
	}
	return nil
}

func (m *MessageRepo) MarkKanboardMsgAsRead(userID, msgID string) error {
	if err := m.redis.SRem(constant.KANBOARD_MESSAGE_UNREADED, userID, msgID); err != nil {
		return err
	}
	if err := m.redis.SAdd(constant.KANBOARD_MESSAGE_READED, userID, msgID, false); err != nil {
		return err
	}
	return nil
}

func (m *MessageRepo) PublishMsg(content string, userID uint, channelName string) error {
	channel := fmt.Sprintf("%d_%s", userID, channelName)
	return m.redis.Publish(channel, content)
}

func (m *MessageRepo) SubscribeMsg(userID uint, channelName string) *redis.PubSub {
	channel := fmt.Sprintf("%d_%s", userID, channelName)
	return m.redis.Subscribe(channel)
}

func (m *MessageRepo) DeleteMsg(msgID string, namespace string) error {
	return m.redis.Delete(namespace, msgID)
}

func (m *MessageRepo) GetMsgField(msgID string, field string) (string, []string) {
	iter := m.redis.ScanByMatch(msgID, 1)
	for iter.Next(m.redis.Ctx) {
		namespace := strings.Split(iter.Val(), "/")
		return m.redis.HGetBykey(iter.Val(), field), namespace
	}
	if iter.Err() == nil {
		return "", nil
	}
	return "", nil
}

func (m *MessageRepo) RemoveAdminMsg(msgID string, userID string) error {
	if err := m.redis.SRem(constant.ADMIN_MESSAGE_UNREADED, userID, msgID); err != nil {
		return err
	}
	if err := m.redis.SRem(constant.ADMIN_MESSAGE_READED, userID, msgID); err != nil {
		return err
	}
	return nil
}

func (m *MessageRepo) RemoveKanboardMsg(msgID string, userID string) error {
	if err := m.redis.SRem(constant.KANBOARD_MESSAGE_UNREADED, userID, msgID); err != nil {
		return err
	}
	if err := m.redis.SRem(constant.KANBOARD_MESSAGE_READED, userID, msgID); err != nil {
		return err
	}
	return nil
}

func (m *MessageRepo) GetAllAdminMsgs(count int64) ([]models.Message, error) {
	admin_messages := []models.Message{}
	admin_iter := m.redis.Scan(constant.ADMIN_NOTIFICATION, "*", count)

	for admin_iter.Next(m.redis.Ctx) {
		key := admin_iter.Val()
		strings := strings.Split(key, "/")
		data := m.redis.HGetAll(strings[0], strings[1])
		to := []uint{}
		if err := json.Unmarshal([]byte(data["to"]), &to); err != nil {
			continue
		}
		message := models.Message{
			ID:        data["id"],
			To:        to,
			Content:   data["content"],
			CreatedAt: data["createdAt"],
		}
		admin_messages = append(admin_messages, message)
	}

	if err := admin_iter.Err(); err != nil {
		return nil, err
	}

	kanboard_messages := []models.Message{}
	kanboard_iter := m.redis.Scan(constant.KANBOARD_NOTIFICATION, "*/*", count)

	for kanboard_iter.Next(m.redis.Ctx) {
		key := kanboard_iter.Val()
		strings := strings.Split(key, "/")
		data := m.redis.HGetAll(fmt.Sprintf("%s/%s", strings[0], strings[1]), strings[2])
		to := []uint{}
		if err := json.Unmarshal([]byte(data["to"]), &to); err != nil {
			continue
		}
		var taskID uint
		if data["taskID"] != "" {
			id, err := strconv.Atoi(data["taskID"])
			if err != nil {
				continue
			}
			taskID = uint(id)
		}
		var projectID uint
		if data["projectID"] != "" {
			id, err := strconv.Atoi(data["projectID"])
			if err != nil {
				projectID = 0
			}
			projectID = uint(id)
		}
		message := models.Message{
			ID:        data["id"],
			To:        to,
			ProjectID: &projectID,
			TaskID:    &taskID,
			Content:   data["content"],
			CreatedAt: data["createdAt"],
		}
		kanboard_messages = append(kanboard_messages, message)
	}

	if err := kanboard_iter.Err(); err != nil {
		return nil, err
	}

	messages := append(admin_messages, kanboard_messages...)

	sort.Slice(messages, func(i, j int) bool {
		t1 := messages[i].CreatedAt
		t2 := messages[j].CreatedAt
		return t1 > t2
	})

	return messages, nil
}

func (m *MessageRepo) GetMsgsByProjectId(count int64, projectId uint) ([]models.Message, error) {
	messages := []models.Message{}
	match := fmt.Sprintf("%d/*", projectId)
	iter := m.redis.Scan(constant.KANBOARD_NOTIFICATION, match, count)

	for iter.Next(m.redis.Ctx) {
		key := iter.Val()
		strings := strings.Split(key, "/")
		data := m.redis.HGetAll(fmt.Sprintf("%s/%s", strings[0], strings[1]), strings[2])
		to := []uint{}
		if err := json.Unmarshal([]byte(data["to"]), &to); err != nil {
			continue
		}
		var taskID uint
		if data["taskID"] != "" {
			id, err := strconv.Atoi(data["taskID"])
			if err != nil {
				continue
			}
			taskID = uint(id)
		}
		var projectID uint
		if data["projectID"] != "" {
			id, err := strconv.Atoi(data["projectID"])
			if err != nil {
				projectID = 0
			}
			projectID = uint(id)
		}
		message := models.Message{
			ID:        data["id"],
			To:        to,
			ProjectID: &projectID,
			TaskID:    &taskID,
			Content:   data["content"],
			CreatedAt: data["createdAt"],
		}
		messages = append(messages, message)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	sort.Slice(messages, func(i, j int) bool {
		t1 := messages[i].CreatedAt
		t2 := messages[j].CreatedAt
		return t1 > t2
	})

	// 只取最新的10条
	if len(messages) > 10 {
		messages = messages[:10]
	}

	return messages, nil
}
