package handlers

import (
	"server/internal/app/admin/services"
	"server/internal/common"

	"github.com/gin-gonic/gin"
)

type PublicHandler struct {
	publicService  *services.PublicService
	projectService *services.ProjectService
	taskService    *services.TaskService
	userService    *services.UserService
}

var publicHandler *PublicHandler

func NewPublicHandler() *PublicHandler {
	if publicHandler == nil {
		publicHandler = &PublicHandler{
			publicService:  services.NewPublicService(),
			projectService: services.NewProjectService(),
			taskService:    services.NewTaskService(),
			userService:    services.NewUserService(),
		}
	}

	return publicHandler
}

func (p PublicHandler) GetCaptcha(ctx *gin.Context) {
	data, err := p.publicService.GetCaptcha()
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

func (p PublicHandler) Dashboard(ctx *gin.Context) {
	data, err := p.publicService.GetDashboard()
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

func (p PublicHandler) GetRecentTasks(ctx *gin.Context) {
	data, err := p.publicService.GetRecentTasks()
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
