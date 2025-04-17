package handlers

import (
	"server/internal/app/kanboard/dto"
	"server/internal/app/kanboard/services"
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

func (t TaskHandler) GetTaskListByProjectId(ctx *gin.Context) {
	var taskIdRequest dto.TaskGetDto

	if err := utils.BindQuery(ctx, &taskIdRequest); err != nil {
		return
	}

	var userIdRequest dto.UserIDRequest
	if err := utils.BindUri(ctx, &userIdRequest); err != nil {
		return
	}

	data, err := t.taskService.GetProjectTaskList(taskIdRequest, userIdRequest.ID)
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

func (t TaskHandler) GetTaskByProjectId(ctx *gin.Context) {
	var taskIdRequest dto.TasksDTO

	if err := utils.BindQuery(ctx, &taskIdRequest); err != nil {
		return
	}

	var userIdRequest dto.UserIDRequest
	if err := utils.BindUri(ctx, &userIdRequest); err != nil {
		return
	}

	data, err := t.taskService.GetProjectTask(taskIdRequest, userIdRequest.ID)
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

func (t TaskHandler) GetTaskListByUserId(ctx *gin.Context) {
	var taskIdRequest dto.TaskGetDto

	if err := utils.BindQuery(ctx, &taskIdRequest); err != nil {
		return
	}

	data, err := t.taskService.GetTaskListByUserId(taskIdRequest)
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

func (t TaskHandler) GetTaskByUserId(ctx *gin.Context) {
	var taskIdRequest dto.TasksDTO

	if err := utils.BindQuery(ctx, &taskIdRequest); err != nil {
		return
	}

	data, err := t.taskService.GetTaskByUserId(taskIdRequest)
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

func (t TaskHandler) CreateTask(ctx *gin.Context) {
	var createRequest dto.TaskCreateDto

	if err := utils.BindRequest(ctx, &createRequest); err != nil {
		return
	}

	data, err := t.taskService.CreateTask(createRequest)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg:  "创建任务成功",
		Data: data,
	})
}

func (t TaskHandler) UpdateTask(ctx *gin.Context) {
	var updateRequest dto.TaskUpdateDto

	if err := utils.BindRequest(ctx, &updateRequest); err != nil {
		return
	}

	var userIdRequest dto.UserIDRequest
	if err := utils.BindUri(ctx, &userIdRequest); err != nil {
		return
	}

	err := t.taskService.UpdateTask(updateRequest, userIdRequest.ID)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "更新任务成功",
	})
}

func (t TaskHandler) UpdateTaskStatus(ctx *gin.Context) {
	var updateRequest dto.TaskChangeStatusDto

	if err := utils.BindRequest(ctx, &updateRequest); err != nil {
		return
	}

	var userIdRequest dto.UserIDRequest
	if err := utils.BindUri(ctx, &userIdRequest); err != nil {
		return
	}

	err := t.taskService.UpdateTaskStatus(updateRequest, userIdRequest.ID)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "更新任务状态成功",
	})
}

func (t TaskHandler) DeleteTask(ctx *gin.Context) {
	var deleteRequest dto.TaskDeleteDto

	if err := utils.BindRequest(ctx, &deleteRequest); err != nil {
		return
	}

	var userIdRequest dto.UserIDRequest
	if err := utils.BindUri(ctx, &userIdRequest); err != nil {
		return
	}

	err := t.taskService.DeleteTask(deleteRequest, userIdRequest.ID)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "删除任务成功",
	})
}

func (t TaskHandler) GetTask(ctx *gin.Context) {
	var getTaskRequest dto.TaskGetInfoDto

	if err := utils.BindQuery(ctx, &getTaskRequest); err != nil {
		return
	}

	var userIdRequest dto.UserIDRequest
	if err := utils.BindUri(ctx, &userIdRequest); err != nil {
		return
	}

	data, err := t.taskService.GetTaskInfo(getTaskRequest, userIdRequest.ID)
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

func (t TaskHandler) AddTaskAssignee(ctx *gin.Context) {
	var request dto.TaskAddAssigneeDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	var userIdRequest dto.UserIDRequest
	if err := utils.BindUri(ctx, &userIdRequest); err != nil {
		return
	}

	err := t.taskService.AddTaskAssignee(request, userIdRequest.ID)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "添加成功",
	})
}

func (t TaskHandler) RemoveTaskAssignee(ctx *gin.Context) {
	var request dto.TaskRemoveAssigneeDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	var userIdRequest dto.UserIDRequest
	if err := utils.BindUri(ctx, &userIdRequest); err != nil {
		return
	}

	err := t.taskService.RemoveTaskAssignee(request, userIdRequest.ID)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "移除成功",
	})
}

func (t TaskHandler) SearchTask(ctx *gin.Context) {
	var request dto.TaskSearchDto

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	var userIdRequest dto.UserIDRequest
	if err := utils.BindUri(ctx, &userIdRequest); err != nil {
		return
	}

	tasks, err := t.taskService.SearchTask(request, userIdRequest.ID)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Data: tasks,
	})
}
