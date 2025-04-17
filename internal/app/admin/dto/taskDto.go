package dto

import (
	"time"

	"server/internal/models"
)

type TaskGetDto struct {
	Id uint `json:"id" form:"id" binding:"required"`
	PageRequest
}

type TaskAssigneeWithAvatar struct {
	Avatar string `json:"avatar"`
	models.TaskAssignee
}

type TaskResponse struct {
	Id        uint                     `json:"id"`
	CreatedAt string                   `json:"created_at"`
	UpdatedAt string                   `json:"updated_at"`
	Title     string                   `json:"title"`
	Desc      string                   `json:"desc"`
	Status    uint                     `json:"status"`
	DueDate   string                   `json:"due_date"`
	Priority  int                      `json:"priority"`
	ProjectId uint                     `json:"project_id"`
	CreatorId UserResponse             `json:"creator"`
	Members   []TaskAssigneeWithAvatar `json:"members"`
}

func (t *TaskResponse) Set(task *models.Task, creator *UserResponse, members *[]TaskAssigneeWithAvatar) *TaskResponse {
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
