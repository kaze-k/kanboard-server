package services

import (
	"errors"
	"time"

	"server/internal/app/kanboard/dto"
	"server/internal/models"
	"server/internal/repositories"
)

type TaskService struct {
	taskRepo          *repositories.TaskRepo
	taskAssigneeRepo  *repositories.TaskAssigneeRepo
	userRepo          *repositories.UserRepo
	resourceRepo      *repositories.ResourceRepo
	projectRepo       *repositories.ProjectRepo
	projectMemberRepo *repositories.ProjectMemberRepo
}

var taskService *TaskService

func NewTaskService() *TaskService {
	if taskService == nil {
		taskService = &TaskService{
			taskRepo:          repositories.NewTaskRepo(),
			taskAssigneeRepo:  repositories.NewTaskAssigneeRepo(),
			userRepo:          repositories.NewUserRepo(),
			resourceRepo:      repositories.NewResourceRepo(),
			projectRepo:       repositories.NewProjectRepo(),
			projectMemberRepo: repositories.NewProjectMemberRepo(),
		}
	}
	return taskService
}

func (t *TaskService) GetProjectTaskList(request dto.TaskGetDto, userId uint) (*dto.TaskPageResponse, error) {
	if !t.projectMemberRepo.CheckProjectMemberExist(request.Id, userId) {
		return nil, errors.New("没有权限")
	}
	task, err := t.taskRepo.GetTaskByProjectIdLimt(request.Id, request.Page, request.PageSize)
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
		var userResponse *dto.UserResponse
		creator := userResponse.Set(user, resource, nil, nil)

		project, err := t.projectRepo.GetProjectById(task.ProjectID)
		if err != nil {
			return nil, err
		}

		var taskResponse dto.TaskResponse
		response := taskResponse.Set(&task, project, creator, &taskAssigneeWithAvatars)
		data = append(data, *response)
	}
	var taskPageResponse dto.TaskPageResponse
	return taskPageResponse.Set(total, request.Page, request.PageSize, data), err
}

func (t *TaskService) GetProjectTask(request dto.TasksDTO, userId uint) ([]dto.TaskResponse, error) {
	if !t.projectMemberRepo.CheckProjectMemberExist(request.Id, userId) {
		return nil, errors.New("没有权限")
	}
	task, err := t.taskRepo.GetTaskByProjectId(request.Id)
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
		var userResponse *dto.UserResponse
		creator := userResponse.Set(user, resource, nil, nil)

		project, err := t.projectRepo.GetProjectById(task.ProjectID)
		if err != nil {
			return nil, err
		}

		var taskResponse dto.TaskResponse
		response := taskResponse.Set(&task, project, creator, &taskAssigneeWithAvatars)
		data = append(data, *response)
	}
	return data, err
}

func (t *TaskService) GetTaskListByUserId(request dto.TaskGetDto) (*dto.TaskPageResponse, error) {
	taskAssignees, err := t.taskAssigneeRepo.GetTaskByUserIdLimt(request.Id, request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}
	total, err := t.taskAssigneeRepo.GetTaskCountByUserId(request.Id)
	if err != nil {
		return nil, err
	}
	data := []dto.TaskResponse{}
	for _, taskAssignee := range *taskAssignees {
		task, err := t.taskRepo.GetTaskById(taskAssignee.TaskID)
		if err != nil {
			return nil, err
		}
		taskAssignees, err := t.taskAssigneeRepo.GetTaskAssigneesByProjectIdAndTankId(task.ProjectID, task.ID)
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
		var userResponse dto.UserResponse
		creator := userResponse.Set(user, resource, nil, nil)

		project, err := t.projectRepo.GetProjectById(task.ProjectID)
		if err != nil {
			return nil, err
		}

		var taskResponse dto.TaskResponse
		response := taskResponse.Set(task, project, creator, &taskAssigneeWithAvatars)
		data = append(data, *response)
	}

	var taskPageResponse dto.TaskPageResponse
	return taskPageResponse.Set(total, request.Page, request.PageSize, data), err
}

func (t *TaskService) GetTaskByUserId(request dto.TasksDTO) ([]dto.TaskResponse, error) {
	taskAssignees, err := t.taskAssigneeRepo.GetTaskByUserId(request.Id)
	if err != nil {
		return nil, err
	}
	data := []dto.TaskResponse{}
	for _, taskAssignee := range *taskAssignees {
		task, err := t.taskRepo.GetTaskById(taskAssignee.TaskID)
		if err != nil {
			return nil, err
		}
		taskAssignees, err := t.taskAssigneeRepo.GetTaskAssigneesByProjectIdAndTankId(task.ProjectID, task.ID)
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
		var userResponse dto.UserResponse
		creator := userResponse.Set(user, resource, nil, nil)

		project, err := t.projectRepo.GetProjectById(task.ProjectID)
		if err != nil {
			return nil, err
		}

		var taskResponse dto.TaskResponse
		response := taskResponse.Set(task, project, creator, &taskAssigneeWithAvatars)
		data = append(data, *response)
	}

	return data, err
}

func (t *TaskService) CreateTask(request dto.TaskCreateDto) (uint, error) {
	if !t.projectMemberRepo.CheckProjectMemberExist(request.ProjectId, request.UserId) {
		return 0, errors.New("没有权限")
	}
	var createTask models.Task

	createTask.CreatorID = request.UserId
	createTask.Title = request.Title
	createTask.Desc = request.Desc
	createTask.ProjectID = request.ProjectId
	if request.Priority != nil {
		createTask.Priority = *request.Priority
	}
	if request.DueDate != nil {
		createTask.DueDate = time.UnixMilli(*request.DueDate)
	}
	task, err := t.taskRepo.CreateTask(createTask)
	if err != nil {
		return 0, err
	}

	if request.Assignees != nil {
		for _, assignee := range *request.Assignees {
			createTaskAssignee := models.TaskAssignee{
				ProjectID: request.ProjectId,
				TaskID:    task.ID,
				UserID:    assignee.UserID,
				Username:  assignee.Username,
			}
			_, err := t.taskAssigneeRepo.CreateTaskAssignee(createTaskAssignee)
			if err != nil {
				return 0, err
			}
		}
	}

	return task.ID, nil
}

func (t *TaskService) DeleteTask(request dto.TaskDeleteDto, userId uint) error {
	task, err := t.taskRepo.GetTaskById(request.Id)
	if err != nil {
		return err
	}
	isAssignee := t.projectMemberRepo.CheckAssignee(task.ProjectID, userId)
	if task.CreatorID != userId && !isAssignee {
		return errors.New("没有权限")
	}
	if err := t.taskRepo.DeleteTaskById(request.Id); err != nil {
		return err
	}
	if err := t.taskAssigneeRepo.DeleteTaskAssigneeById(request.Id); err != nil {
		return err
	}
	return nil
}

func (t *TaskService) UpdateTaskStatus(request dto.TaskChangeStatusDto, userId uint) error {
	if !t.projectMemberRepo.CheckProjectMemberExist(request.ProjectId, userId) {
		return errors.New("没有权限")
	}
	values := make(map[string]any)
	if request.Status != nil {
		values["status"] = request.Status
	}
	err := t.taskRepo.UpdateTask(values, request.Id, request.ProjectId)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskService) UpdateTask(request dto.TaskUpdateDto, userId uint) error {
	task, err := t.taskRepo.GetTaskById(request.Id)
	if err != nil {
		return err
	}
	isAssignee := t.projectMemberRepo.CheckAssignee(task.ProjectID, userId)
	if task.CreatorID != userId && !isAssignee {
		return errors.New("没有权限")
	}
	values := make(map[string]any)
	if request.Desc != nil {
		values["desc"] = *request.Desc
	}
	if request.Priority != nil {
		values["priority"] = *request.Priority
	}
	if request.DueDate != nil {
		values["due_date"] = time.UnixMilli(*request.DueDate)
	}
	err = t.taskRepo.UpdateTask(values, request.Id, request.ProjectId)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskService) AddTaskAssignee(request dto.TaskAddAssigneeDto, userId uint) error {
	task, err := t.taskRepo.GetTaskById(request.Id)
	if err != nil {
		return err
	}
	isAssignee := t.projectMemberRepo.CheckAssignee(task.ProjectID, userId)
	if task.CreatorID != userId && !isAssignee {
		return errors.New("没有权限")
	}
	createTasks := []models.TaskAssignee{}
	for _, assignee := range request.Assignees {
		createTaskAssignee := models.TaskAssignee{
			ProjectID: request.ProjectId,
			TaskID:    request.Id,
			UserID:    assignee.UserID,
			Username:  assignee.Username,
		}
		createTasks = append(createTasks, createTaskAssignee)
	}

	err = t.taskAssigneeRepo.AddTaskAssignee(createTasks, request.ProjectId, request.Id)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskService) RemoveTaskAssignee(request dto.TaskRemoveAssigneeDto, userId uint) error {
	task, err := t.taskRepo.GetTaskById(request.Id)
	if err != nil {
		return err
	}
	isAssignee := t.projectMemberRepo.CheckAssignee(task.ProjectID, userId)
	if task.CreatorID != userId && !isAssignee {
		return errors.New("没有权限")
	}
	err = t.taskAssigneeRepo.RemoveAssignee(request.Id, request.ProjectId, request.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskService) GetTaskInfo(request dto.TaskGetInfoDto, userId uint) (*dto.TaskWithMemberResponse, error) {
	if !t.projectMemberRepo.CheckProjectMemberExist(request.ProjectId, userId) {
		return nil, errors.New("没有权限")
	}
	task, err := t.taskRepo.GetTaskByIdAndProjectId(request.TaskId, request.ProjectId)
	if err != nil {
		return nil, err
	}
	taskAssignees, err := t.taskAssigneeRepo.GetTaskAssigneesByProjectIdAndTankId(request.ProjectId, task.ID)
	if err != nil {
		return nil, err
	}

	taskAssigneeResponses := []dto.UserResponse{}
	for _, assignee := range *taskAssignees {
		user, err := t.userRepo.GetUserById(assignee.UserID)
		if err != nil {
			return nil, err
		}
		resource, err := t.resourceRepo.GetResourceById(user.Avatar)
		if err != nil {
			return nil, err
		}
		isAssignee := t.projectMemberRepo.CheckAssignee(request.ProjectId, assignee.UserID)
		var userResposne dto.UserResponse
		taskAssigneeResponse := userResposne.Set(user, resource, nil, &isAssignee)
		taskAssigneeResponses = append(taskAssigneeResponses, *taskAssigneeResponse)
	}

	user, err := t.userRepo.GetUserById(task.CreatorID)
	if err != nil {
		return nil, err
	}
	resource, err := t.resourceRepo.GetResourceById(user.Avatar)
	if err != nil {
		return nil, err
	}
	isAssignee := t.projectMemberRepo.CheckAssignee(request.ProjectId, user.ID)
	var UserResponse *dto.UserResponse
	creator := UserResponse.Set(user, resource, nil, &isAssignee)

	project, err := t.projectRepo.GetProjectById(task.ProjectID)
	if err != nil {
		return nil, err
	}

	var taskResponse dto.TaskWithMemberResponse
	response := taskResponse.Set(task, project, creator, &taskAssigneeResponses)

	return response, nil
}

func (t *TaskService) SearchTask(request dto.TaskSearchDto, userId uint) ([]dto.TaskResponse, error) {
	taskResponses := []dto.TaskResponse{}
	if request.ProjectId == nil {
		return taskResponses, nil
	}
	projectId := *request.ProjectId
	if !t.projectMemberRepo.CheckProjectMemberExist(projectId, userId) {
		return nil, errors.New("没有权限")
	}

	query := make(map[string]any)
	if request.Priority != nil {
		query["priority"] = *request.Priority
	}
	if request.Title != nil {
		query["title"] = *request.Title
	}
	if request.CreatorId != nil {
		query["creatorId"] = *request.CreatorId
	}

	tasks, err := t.taskRepo.SearchTask(query, projectId)
	if err != nil {
		return nil, err
	}

	filteredTasks := *tasks

	if request.UserId != nil {
		taskAssignees, err := t.taskAssigneeRepo.SearchTask(*request.UserId, projectId)
		if err != nil {
			return nil, err
		}

		tasksByTaskIds := []models.Task{}
		for _, taskAssignee := range *taskAssignees {
			tasksByTaskId, err := t.taskRepo.GetTaskById(taskAssignee.TaskID)
			if err != nil {
				return nil, err
			}

			if title, ok := query["title"].(string); ok {
				if tasksByTaskId.Title != title {
					continue
				}
			}

			if priority, ok := query["priority"].(int); ok {
				if tasksByTaskId.Priority != priority {
					continue
				}
			}

			if creatorId, ok := query["creatorId"].(uint); ok {
				if tasksByTaskId.CreatorID != creatorId {
					continue
				}
			}

			tasksByTaskIds = append(tasksByTaskIds, *tasksByTaskId)
		}

		filteredTasks = tasksByTaskIds
	}

	data := []dto.TaskResponse{}
	for _, task := range filteredTasks {
		taskAssignees, err := t.taskAssigneeRepo.GetTaskAssigneesByProjectIdAndTankId(task.ProjectID, task.ID)
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
		var userResponse *dto.UserResponse
		creator := userResponse.Set(user, resource, nil, nil)

		project, err := t.projectRepo.GetProjectById(task.ProjectID)
		if err != nil {
			return nil, err
		}

		var taskResponse dto.TaskResponse
		response := taskResponse.Set(&task, project, creator, &taskAssigneeWithAvatars)
		data = append(data, *response)
	}
	return data, err
}
