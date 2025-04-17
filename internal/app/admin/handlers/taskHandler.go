package handlers

import (
	"server/internal/app/admin/dto"
	"server/internal/app/admin/services"
	"server/internal/common"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *services.TaskService
}

var taskHandler *TaskHandler

func NewTaskHandler() *TaskHandler {
	if taskHandler == nil {
		taskHandler = &TaskHandler{
			taskService: services.NewTaskService(),
		}
	}
	return taskHandler
}

func (t TaskHandler) GetProjectTask(ctx *gin.Context) {
	var request dto.TaskGetDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	data, err := t.taskService.GetProjectTask(request)
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
