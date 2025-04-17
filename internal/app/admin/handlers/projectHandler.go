package handlers

import (
	"server/internal/app/admin/dto"
	"server/internal/app/admin/services"
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

func (p ProjectHandler) GetProjectList(ctx *gin.Context) {
	var request dto.PageRequest

	if err := utils.BindQuery(ctx, &request); err != nil {
		return
	}

	data, err := p.projectService.GetProjectList(&request)
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

func (p ProjectHandler) GetProject(ctx *gin.Context) {
	var request dto.ProjectIdDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	data, err := p.projectService.GetProjectById(request.Id)
	if err != nil {
		common.NotFound(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Data: data,
	})
}

func (p ProjectHandler) CreateProject(ctx *gin.Context) {
	var request dto.ProjectCreateDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	err := p.projectService.CreateProject(&request)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "创建项目成功",
	})
}

func (p ProjectHandler) UpdateProject(ctx *gin.Context) {
	var request dto.ProjectUpdateDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	err := p.projectService.UpdateProject(&request)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "更新项目成功",
	})
}

func (p ProjectHandler) DeleteProject(ctx *gin.Context) {
	var request dto.ProjectIdDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	err := p.projectService.DeleteProject(&request)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "删除项目成功",
	})
}

func (p ProjectHandler) SetProjectAssignee(ctx *gin.Context) {
	var request dto.ProjectAssigneeDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	err := p.projectService.SetProjectAssignee(&request)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "添加负责人成功",
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
