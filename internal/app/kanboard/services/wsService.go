package services

import (
	"sync"

	"server/internal/common"
	"server/internal/constant"
	"server/internal/global"
	"server/internal/repositories"
	"server/pkg/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsService struct {
	upgrader       websocket.Upgrader
	clients        map[uint]*websocket.Conn
	mutex          sync.Mutex
	msgRepo        *repositories.MessageRepo
	messageService *MessageService
}

var wsService *WsService

func NewWsService() *WsService {
	return &WsService{
		upgrader:       *ws.NewWsUpgrader(),
		clients:        make(map[uint]*websocket.Conn),
		msgRepo:        repositories.NewMessageRepo(),
		messageService: NewMessageService(),
	}
}

func (w *WsService) PushUnReadMsg(conn *websocket.Conn, userID uint) {
	data := w.messageService.GetUnReadMsgs(userID)
	json := common.MsgRsp{
		Type:        constant.UNREAD_MESSAGE,
		UnReadCount: len(data),
		Payload:     data,
	}
	conn.WriteJSON(json)
}

func (w *WsService) ListenAndPush(conn *websocket.Conn, userID uint) {
	pubsub := w.msgRepo.SubscribeMsg(userID, constant.KANBOARD_MESSAGE_CHANNEL)
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(global.Redis.Ctx)
		if err != nil {
			global.Logger.Errorw("listen message error", "error", err)
			break
		}

		data := w.messageService.GetUnReadMsgs(userID)
		var json common.MsgRsp
		if msg.Payload == constant.PUBLISH_MESSAGE {
			json = common.MsgRsp{
				Type:        constant.PUBLISH_MESSAGE,
				UnReadCount: len(data),
			}
		} else if msg.Payload == constant.UPDATE_USER {
			json = common.MsgRsp{
				Type:        constant.UPDATE_USER,
				UnReadCount: len(data),
			}
		} else if msg.Payload == constant.NEW_TASK_STATUS {
			json = common.MsgRsp{
				Type:        constant.NEW_TASK_STATUS,
				UnReadCount: len(data),
			}
		} else {
			json = common.MsgRsp{
				Type:        constant.NEW_MESSAGE,
				UnReadCount: len(data),
				Payload:     msg.Payload,
			}
		}

		conn.WriteJSON(json)
	}
}

func (w *WsService) AddClient(conn *websocket.Conn, userID uint) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.clients[userID] = conn
	global.Logger.Infow("add client", "client", conn.RemoteAddr().String())
}

func (w *WsService) RemoveClient(conn *websocket.Conn, userID uint) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	delete(w.clients, userID)
	conn.Close()
	global.Logger.Infow("remove client", "client", conn.RemoteAddr().String())
}

func (w *WsService) HandleWebsocket(ctx *gin.Context, userID uint) {
	conn, err := w.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		global.Logger.Errorw("upgrade error", "error", err)
		return
	}
	defer w.RemoveClient(conn, userID)

	w.AddClient(conn, userID)

	// 上线推送未读消息
	w.PushUnReadMsg(conn, userID)

	// 监听Redis并推送消息
	w.ListenAndPush(conn, userID)

	delete(w.clients, userID)
}
