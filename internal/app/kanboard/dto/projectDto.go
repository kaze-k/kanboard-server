package dto

import (
	"time"

	"server/internal/models"
)

type ProjectIdDto struct {
	Id uint `json:"id" form:"id" binding:"required"`
}

type ProjectAddMemberDto struct {
	ProjectId uint            `json:"project_id" form:"project_id" binding:"required"`
	UserId    uint            `json:"user_id" form:"user_id" binding:"required"`
	Members   []models.Member `json:"members" form:"members" binding:"required"`
}

type ProjectRemoveMemberDto struct {
	ProjectId uint `json:"project_id" form:"project_id" binding:"required"`
	UserId    uint `json:"user_id" form:"user_id" binding:"required"`
	MemberId  uint `json:"member_id" form:"member_id" binding:"required"`
}

type MemberResponse struct {
	Assignee bool `json:"assignee"`
	UserResponse
}

func (m *MemberResponse) Set(user *models.User, assignee bool, resource *models.Resource) *MemberResponse {
	return &MemberResponse{
		Assignee: assignee,
		UserResponse: UserResponse{
			ID:         user.ID,
			CreatedAt:  user.CreatedAt.Local().Format(time.DateTime),
			Username:   user.Username,
			Avatar:     resource.StaticPath,
			Gender:     user.Gender,
			Email:      user.Email,
			Mobile:     user.Mobile,
			CreateFrom: user.CreateFrom,
			Position:   user.Position,
		},
	}
}

type Statistics struct {
	TaskTotal      int64 `json:"task_total"`
	TaskInProgress int64 `json:"task_in_progress"`
	TaskDone       int64 `json:"task_done"`
	TaskHigh       int64 `json:"task_high"`
	TaskMedium     int64 `json:"task_medium"`
	TaskLow        int64 `json:"task_low"`
}

type ProjectWithUserResponse struct {
	Id         uint           `json:"id"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
	Name       string         `json:"name"`
	Desc       string         `json:"desc"`
	Members    []UserResponse `json:"members"`
	Statistics Statistics     `json:"statistics"`
}

func (r *ProjectWithUserResponse) Set(project *models.Project, users []UserResponse, statistics *Statistics) *ProjectWithUserResponse {
	var projectResponse ProjectWithUserResponse
	projectResponse.Id = project.ID
	projectResponse.CreatedAt = project.CreatedAt.Local().Format(time.DateTime)
	projectResponse.UpdatedAt = project.UpdatedAt.Local().Format(time.DateTime)
	projectResponse.Name = project.Name
	if project.Desc != nil {
		projectResponse.Desc = *project.Desc
	}
	if users != nil {
		projectResponse.Members = users
	}
	if statistics != nil {
		projectResponse.Statistics = *statistics
	}
	return &projectResponse
}
