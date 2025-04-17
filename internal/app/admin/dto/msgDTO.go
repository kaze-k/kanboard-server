package dto

import "server/internal/models"

type MsgMarkReadDto struct {
	UserId uint   `json:"id" form:"id" binding:"required"`
	MsgId  string `json:"msg_id" form:"msg_id" binding:"required"`
}

type MsgDeleteDto struct {
	MsgId string `json:"msg_id" form:"msg_id" binding:"required"`
}

type MsgGetAllResponse struct {
	ID        string          `json:"id"`
	To        []models.Member `json:"to"`
	Content   string          `json:"content"`
	CreatedAt string          `json:"created_at"`
}
