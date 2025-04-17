package handlers

import (
	"server/internal/app/kanboard/dto"
	"server/internal/app/kanboard/services"
	"server/internal/common"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService *services.ProjectService
}

var projectHandler *ProjectHandler

func NewProjectHandler() *ProjectHandler {
	if projectHandler == nil {
		projectHandler = &ProjectHandler{
			projectService: services.NewProjectService(),
		}
	}

	return projectHandler
}

func (p ProjectHandler) GetProjectById(ctx *gin.Context) {
	var projectIdRequest dto.ProjectIdDto

	if err := utils.BindQuery(ctx, &projectIdRequest); err != nil {
		return
	}

	data, err := p.projectService.GetProjectById(projectIdRequest.Id)
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

func (p ProjectHandler) AddProjectMember(ctx *gin.Context) {
	var request dto.ProjectAddMemberDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	err := p.projectService.AddProjectMember(&request)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "添加项目成员成功",
	})
}

func (p ProjectHandler) RemoveProjectMember(ctx *gin.Context) {
	var request dto.ProjectRemoveMemberDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	err := p.projectService.RemoveProjectMember(&request)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "移除项目成员成功",
	})
}

func (p ProjectHandler) GetProjectMembers(ctx *gin.Context) {
	var request dto.ProjectIdDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	var userIdRequest dto.UserIDRequest
	if err := utils.BindUri(ctx, &userIdRequest); err != nil {
		return
	}

	members, err := p.projectService.GetProjectMembers(request.Id, userIdRequest.ID)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Data: members,
	})
}

func (p ProjectHandler) GetMembers(ctx *gin.Context) {
	var request dto.ProjectIdDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	var userIdRequest dto.UserIDRequest
	if err := utils.BindUri(ctx, &userIdRequest); err != nil {
		return
	}

	members, err := p.projectService.GetMembers(request.Id, userIdRequest.ID)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Data: members,
	})
}
