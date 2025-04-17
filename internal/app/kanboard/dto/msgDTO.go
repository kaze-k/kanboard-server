package dto

type MsgMarkReadDto struct {
	UserId uint   `json:"id" form:"id" binding:"required"`
	MsgId  string `json:"msg_id" form:"msg_id" binding:"required"`
}

type MsgResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
