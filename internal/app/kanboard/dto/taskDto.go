package dto

import (
	"time"

	"server/internal/models"
)

type TasksDTO struct {
	Id uint `json:"id" form:"id" uri:"id" binding:"required"`
}

type TaskGetDto struct {
	Id uint `json:"id" form:"id" uri:"id" binding:"required"`
	PageRequest
}

type TaskCreateDto struct {
	UserId    uint             `json:"user_id" form:"user_id" binding:"required"`
	ProjectId uint             `json:"project_id" form:"project_id" binding:"required"`
	Title     string           `json:"title" form:"title" binding:"required"`
	Desc      string           `json:"desc" form:"desc" binding:"required"`
	Priority  *int             `json:"priority" form:"priority" binding:"required"`
	DueDate   *int64           `json:"due_date" form:"due_date"`
	Assignees *[]models.Member `json:"assignees" form:"assignees" binding:"required"`
}

type TaskDeleteDto struct {
	Id        uint `json:"id" form:"id" binding:"required"`
	ProjectId uint `json:"project_id" form:"project_id" binding:"required"`
}

type TaskUpdateDto struct {
	Id        uint    `json:"id" form:"id" binding:"required"`
	ProjectId uint    `json:"project_id" form:"project_id" binding:"required"`
	Desc      *string `json:"desc" form:"desc"`
	Priority  *int    `json:"priority" form:"priority"`
	DueDate   *int64  `json:"due_date" form:"due_date"`
}

type TaskChangeStatusDto struct {
	Id        uint  `json:"id" form:"id" binding:"required"`
	ProjectId uint  `json:"project_id" form:"project_id" binding:"required"`
	Status    *uint `json:"status" form:"status" binding:"required"`
}

type TaskAddAssigneeDto struct {
	Id        uint            `json:"id" form:"id" binding:"required"`
	ProjectId uint            `json:"project_id" form:"project_id" binding:"required"`
	Assignees []models.Member `json:"assignees" form:"assignees" binding:"required"`
}

type TaskRemoveAssigneeDto struct {
	Id        uint `json:"id" form:"id" binding:"required"`
	ProjectId uint `json:"project_id" form:"project_id" binding:"required"`
	UserId    uint `json:"user_id" form:"user_id" binding:"required"`
}

type TaskGetInfoDto struct {
	ProjectId uint `json:"project_id" form:"project_id" uri:"project_id" binding:"required"`
	TaskId    uint `json:"task_id" form:"task_id" uri:"task_id" binding:"required"`
}

type TaskSearchDto struct {
	ProjectId *uint   `json:"project_id" form:"project_id"`
	Title     *string `json:"title" form:"title"`
	Priority  *int    `json:"priority" form:"priority"`
	UserId    *uint   `json:"user_id" form:"user_id"`
	CreatorId *uint   `json:"creator_id" form:"creator_id"`
}

type TaskAssigneeWithAvatar struct {
	Avatar string `json:"avatar"`
	models.TaskAssignee
}

type TaskResponse struct {
	Id          uint                     `json:"id"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
	Title       string                   `json:"title"`
	Desc        string                   `json:"desc"`
	Status      uint                     `json:"status"`
	DueDate     string                   `json:"due_date"`
	Priority    int                      `json:"priority"`
	ProjectId   uint                     `json:"project_id"`
	ProjectName string                   `json:"project_name"`
	CreatorId   UserResponse             `json:"creator"`
	Members     []TaskAssigneeWithAvatar `json:"members"`
}

func (t *TaskResponse) Set(task *models.Task, project *models.Project, creator *UserResponse, members *[]TaskAssigneeWithAvatar) *TaskResponse {
	t.Id = task.ID
	t.CreatedAt = task.CreatedAt.Local().Format(time.DateTime)
	t.UpdatedAt = task.UpdatedAt.Local().Format(time.DateTime)
	t.Title = task.Title
	t.Desc = task.Desc
	t.Status = task.Status
	if !task.DueDate.IsZero() {
		t.DueDate = task.DueDate.Local().Format(time.DateTime)
	}
	t.Priority = task.Priority
	t.ProjectId = task.ProjectID
	t.ProjectName = project.Name
	t.CreatorId = *creator
	t.Members = *members
	return t
}

type TaskWithMemberResponse struct {
	Id          uint           `json:"id"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
	Title       string         `json:"title"`
	Desc        string         `json:"desc"`
	Status      uint           `json:"status"`
	DueDate     string         `json:"due_date"`
	Priority    int            `json:"priority"`
	ProjectId   uint           `json:"project_id"`
	ProjectName string         `json:"project_name"`
	CreatorId   UserResponse   `json:"creator"`
	Members     []UserResponse `json:"members"`
}

func (t *TaskWithMemberResponse) Set(task *models.Task, project *models.Project, creator *UserResponse, members *[]UserResponse) *TaskWithMemberResponse {
	t.Id = task.ID
	t.CreatedAt = task.CreatedAt.Local().Format(time.DateTime)
	t.UpdatedAt = task.UpdatedAt.Local().Format(time.DateTime)
	t.Title = task.Title
	t.Desc = task.Desc
	t.Status = task.Status
	if !task.DueDate.IsZero() {
		t.DueDate = task.DueDate.Local().Format(time.DateTime)
	}
	t.Priority = task.Priority
	t.ProjectId = task.ProjectID
	t.ProjectName = project.Name
	t.CreatorId = *creator
	t.Members = *members
	return t
}

type TaskPageResponse struct {
	Total     int            `json:"total"`
	Page      int            `json:"page"`
	PageSize  int            `json:"size"`
	TotalPage int            `json:"total_page"`
	Data      []TaskResponse `json:"data"`
}

func (t *TaskPageResponse) Set(total int64, page int, pageSize int, data []TaskResponse) *TaskPageResponse {
	var taskPageResponse TaskPageResponse

	taskPageResponse.Total = int(total)
	taskPageResponse.Page = page
	taskPageResponse.PageSize = pageSize
	totalPage := int(float64(total) / float64(pageSize))
	if total%int64(pageSize) != 0 {
		totalPage++
	}
	taskPageResponse.TotalPage = totalPage
	taskPageResponse.Data = data
	return &taskPageResponse
}
