package services

import (
	"server/internal/app/admin/dto"
	"server/internal/repositories"
)

type TaskService struct {
	taskRepo         *repositories.TaskRepo
	taskAssigneeRepo *repositories.TaskAssigneeRepo
	userRepo         *repositories.UserRepo
	resourceRepo     *repositories.ResourceRepo
}

var taskService *TaskService

func NewTaskService() *TaskService {
	if taskService == nil {
		taskService = &TaskService{
			taskRepo:         repositories.NewTaskRepo(),
			taskAssigneeRepo: repositories.NewTaskAssigneeRepo(),
			userRepo:         repositories.NewUserRepo(),
			resourceRepo:     repositories.NewResourceRepo(),
		}
	}
	return taskService
}

func (t *TaskService) GetProjectTask(request dto.TaskGetDto) (*dto.TaskPageResponse, error) {
	task, err := t.taskRepo.GetTaskByProjectId(request.Id)
	if err != nil {
		return nil, err
	}
	total, err := t.taskRepo.GetTaskCountByProjectId(request.Id)
	if err != nil {
		return nil, err
	}
	data := []dto.TaskResponse{}
	for _, task := range *task {
		taskAssignees, err := t.taskAssigneeRepo.GetTaskAssigneesByProjectIdAndTankId(request.Id, task.ID)
		if err != nil {
			return nil, err
		}

		taskAssigneeWithAvatars := []dto.TaskAssigneeWithAvatar{}
		for _, assignee := range *taskAssignees {
			user, err := t.userRepo.GetUserById(assignee.UserID)
			if err != nil {
				return nil, err
			}
			resource, err := t.resourceRepo.GetResourceById(user.Avatar)
			if err != nil {
				return nil, err
			}
			taskAssigneeWithAvatar := dto.TaskAssigneeWithAvatar{
				Avatar:       resource.StaticPath,
				TaskAssignee: assignee,
			}
			taskAssigneeWithAvatars = append(taskAssigneeWithAvatars, taskAssigneeWithAvatar)
		}

		user, err := t.userRepo.GetUserById(task.CreatorID)
		if err != nil {
			return nil, err
		}
		resource, err := t.resourceRepo.GetResourceById(user.Avatar)
		if err != nil {
			return nil, err
		}
		var UserResponse dto.UserResponse
		creator := UserResponse.Set(user, resource)

		var taskResponse dto.TaskResponse
		response := taskResponse.Set(&task, creator, &taskAssigneeWithAvatars)
		data = append(data, *response)
	}
	var taskResponse dto.TaskPageResponse
	return taskResponse.Set(total, request.Page, request.PageSize, data), err
}
