package dto

import (
	"time"

	"server/internal/models"
)

type DashboardProject struct {
	Count     int64    `json:"count"`
	CreatedAt []string `json:"created_at"`
}

type DashboardTask struct {
	Count     int64    `json:"count"`
	Undo      int64    `json:"undo"`
	Done      int64    `json:"done"`
	High      int64    `json:"high"`
	Medium    int64    `json:"medium"`
	Low       int64    `json:"low"`
	CreatedAt []string `json:"created_at"`
}

type DashboardUser struct {
	Count              int64    `json:"count"`
	Male               int64    `json:"male"`
	Female             int64    `json:"female"`
	CreateFromKanboard int64    `json:"create_from_kanboard"`
	CreateFromAdmin    int64    `json:"create_from_admin"`
	Admin              int64    `json:"admin"`
	Loginable          int64    `json:"loginable"`
	Unloginable        int64    `json:"unloginable"`
	CreatedAt          []string `json:"created_at"`
}

type DashboardResponse struct {
	Project DashboardProject `json:"project"`
	Task    DashboardTask    `json:"task"`
	User    DashboardUser    `json:"user"`
}

type TaskRecentResponse struct {
	Id          int    `json:"id"`
	CreatedAt   string `json:"created_at"`
	Title       string `json:"title"`
	Status      uint   `json:"status"`
	Priority    int    `json:"priority"`
	ProjectName string `json:"project_name"`
}

func (t *TaskRecentResponse) Set(id int, task *models.Task, project *models.Project) *TaskRecentResponse {
	t.Id = id
	t.CreatedAt = task.CreatedAt.Local().Format(time.DateTime)
	t.Title = task.Title
	t.Status = task.Status
	t.Priority = task.Priority
	t.ProjectName = project.Name
	return t
}
