package handlers

import (
	"server/internal/app/admin/dto"
	"server/internal/app/admin/services"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
)

type WSHandler struct {
	wsService *services.WsService
}

var wsHandler *WSHandler

func NewWSHandler() *WSHandler {
	if wsHandler == nil {
		wsHandler = &WSHandler{
			wsService: services.NewWsService(),
		}
	}
	return wsHandler
}

func (w WSHandler) WebSocket(ctx *gin.Context) {
	var request dto.UserIDRequest

	if err := utils.BindUri(ctx, &request); err != nil {
		return
	}

	w.wsService.HandleWebsocket(ctx, request.ID)
}
