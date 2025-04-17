package handlers

import (
	"server/internal/app/kanboard/dto"
	"server/internal/app/kanboard/services"
	"server/internal/common"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService *services.MessageService
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		messageService: services.NewMessageService(),
	}
}

func (m MessageHandler) GetUnReadMsgs(ctx *gin.Context) {
	var request dto.UserIDRequest

	if err := utils.BindUri(ctx, &request); err != nil {
		return
	}

	data := m.messageService.GetUnReadMsgs(request.ID)

	common.Ok(ctx, common.RspOpts{
		Data: data,
	})
}

func (m MessageHandler) GetReadedMsgs(ctx *gin.Context) {
	var request dto.UserIDRequest

	if err := utils.BindUri(ctx, &request); err != nil {
		return
	}

	data := m.messageService.GetReadedMsgs(request.ID)

	common.Ok(ctx, common.RspOpts{
		Data: data,
	})
}

func (m MessageHandler) MarkReadMsg(ctx *gin.Context) {
	var request dto.MsgMarkReadDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	err := m.messageService.MarkReadMsg(request.UserId, request.MsgId)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{})
}

func (m MessageHandler) GetMsgsByProjectId(ctx *gin.Context) {
	var request dto.ProjectIdDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	data, err := m.messageService.GetMsgsByProjectId(request.Id)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Data: data,
	})
}
