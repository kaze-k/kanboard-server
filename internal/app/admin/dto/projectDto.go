package dto

import (
	"time"

	"server/internal/models"
)

type ProjectIdDto struct {
	Id uint `json:"id" binding:"required" uri:"id"`
}

type ProjectCreateDto struct {
	Name      string           `json:"name" form:"name" binding:"required"`
	Desc      *string          `json:"desc" form:"desc"`
	Assignees *[]models.Member `json:"assignees" form:"assignees" binding:"required"`
	Members   *[]models.Member `json:"members" form:"members"`
}

type ProjectUpdateDto struct {
	Id   uint    `json:"id" form:"id" binding:"required"`
	Name *string `json:"name" form:"name"`
	Desc *string `json:"desc" form:"desc"`
}

type ProjectAddMemberDto struct {
	ProjectId uint             `json:"project_id" form:"project_id" binding:"required"`
	Assignees *[]models.Member `json:"assignees" form:"assignees" binding:"required"`
	Members   []models.Member  `json:"members" form:"members" binding:"required"`
}

type ProjectRemoveMemberDto struct {
	ProjectId uint   `json:"project_id" form:"project_id" binding:"required"`
	Members   []uint `json:"members" form:"members" binding:"required"`
}

type ProjectAssigneeDto struct {
	ProjectId uint   `json:"project_id" form:"project_id" binding:"required"`
	Members   []uint `json:"members" form:"members" binding:"required"`
	Value     *bool  `json:"value" form:"value" binding:"required"`
}

type ProjectPageResponse struct {
	Total     int               `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
	TotalPage int               `json:"total_page"`
	Data      []ProjectResponse `json:"data"`
}

func (r *ProjectPageResponse) Set(total int64, page int, pageSize int, data []ProjectResponse) *ProjectPageResponse {
	var projectPageResponse ProjectPageResponse

	projectPageResponse.Total = int(total)
	projectPageResponse.Page = page
	projectPageResponse.PageSize = pageSize
	totalPage := int(float64(total) / float64(pageSize))
	if total%int64(pageSize) != 0 {
		totalPage++
	}
	projectPageResponse.TotalPage = totalPage
	projectPageResponse.Data = data

	return &projectPageResponse
}

type ProjectResponseWitheAvatar struct {
	Avatar string `json:"avatar"`
	models.ProjectMember
}

type ProjectResponse struct {
	Id        uint                         `json:"id"`
	CreatedAt string                       `json:"created_at"`
	UpdatedAt string                       `json:"updated_at"`
	Name      string                       `json:"name"`
	Desc      string                       `json:"desc"`
	Members   []ProjectResponseWitheAvatar `json:"members"`
}

func (r *ProjectResponse) Set(project *models.Project, members *[]ProjectResponseWitheAvatar) *ProjectResponse {
	var projectResponse ProjectResponse
	projectResponse.Id = project.ID
	projectResponse.CreatedAt = project.CreatedAt.Local().Format(time.DateTime)
	projectResponse.UpdatedAt = project.UpdatedAt.Local().Format(time.DateTime)
	projectResponse.Name = project.Name
	if project.Desc != nil {
		projectResponse.Desc = *project.Desc
	}
	if members != nil {
		projectResponse.Members = *members
	}
	return &projectResponse
}

type MemberResponse struct {
	Assignee bool `json:"assignee"`
	UserResponse
}

func (m *MemberResponse) Set(user *models.User, assignee bool, resource *models.Resource) *MemberResponse {
	return &MemberResponse{
		Assignee: assignee,
		UserResponse: UserResponse{
			CreateAt:   user.CreatedAt.Local().Format(time.DateTime),
			UpdateAt:   user.UpdatedAt.Local().Format(time.DateTime),
			ID:         user.ID,
			Username:   user.Username,
			Avatar:     resource.StaticPath,
			Gender:     user.Gender,
			Email:      user.Email,
			Mobile:     user.Mobile,
			CreateFrom: user.CreateFrom,
			Loginable:  user.Loginable,
			IsAdmin:    user.IsAdmin,
			Position:   user.Position,
		},
	}
}

type ProjectWithUserResponse struct {
	Id        uint             `json:"id"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
	Name      string           `json:"name"`
	Desc      string           `json:"desc"`
	Members   []MemberResponse `json:"members"`
}

func (r *ProjectWithUserResponse) Set(project *models.Project, users []MemberResponse) *ProjectWithUserResponse {
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
	return &projectResponse
}
