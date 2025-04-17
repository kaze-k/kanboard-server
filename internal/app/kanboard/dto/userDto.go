package dto

import (
	"time"

	"server/internal/constant"
	"server/internal/models"
)

type UserIDRequest struct {
	ID uint `uri:"id" binding:"required"`
}

type UserLoginRequest struct {
	Username      string `json:"username" binding:"required" from:"username"`
	Password      string `json:"password" binding:"required" from:"password"`
	CaptchaID     string `json:"captchaId" binding:"required" from:"captchaId"`
	CaptchaAnswer string `json:"captchaAnswer" binding:"required" from:"captchaAnswer"`
}

type UserRegisterRequest struct {
	Username      string          `json:"username" binding:"required" form:"username"`
	Password      string          `json:"password" binding:"required" form:"password"`
	Email         *string         `json:"email" form:"email"`
	Mobile        *string         `json:"mobile" form:"mobile"`
	Gender        constant.Gender `json:"gender" form:"gender"`
	CaptchaID     string          `json:"captchaId" binding:"required" from:"captchaId"`
	CaptchaAnswer string          `json:"captchaAnswer" binding:"required" from:"captchaAnswer"`
}

type UserUpdateRequest struct {
	ID     uint    `json:"id" form:"id" binding:"required"`
	Avatar *uint   `json:"avatar" form:"avatar"`
	Email  *string `json:"email" form:"email"`
	Mobile *string `json:"mobile" form:"mobile"`
}

type UserUpdatePasswordRequest struct {
	ID      uint   `json:"id" form:"id" required:"true"`
	Current string `json:"current" form:"current" required:"true"`
	New     string `json:"new" form:"new" required:"true"`
}

type ProjectsWithIdAndAssignee struct {
	ProjectID   uint   `json:"project_id"`
	Assignee    bool   `json:"assignee"`
	ProjectName string `json:"project_name"`
	JoinedAt    string `json:"joined_at"`
}

type UserResponse struct {
	ID         uint                        `json:"id"`
	CreatedAt  string                      `json:"created_at"`
	Username   string                      `json:"username"`
	Avatar     string                      `json:"avatar"`
	Gender     constant.Gender             `json:"gender"`
	Email      string                      `json:"email"`
	Mobile     string                      `json:"mobile"`
	CreateFrom constant.From               `json:"create_from"`
	Position   string                      `json:"position"`
	Projects   []ProjectsWithIdAndAssignee `json:"projects"`
	Assignee   *bool                       `json:"assignee,omitempty"`
}

func (r *UserResponse) Set(user *models.User, resource *models.Resource, Projects []ProjectsWithIdAndAssignee, assignee *bool) *UserResponse {
	return &UserResponse{
		ID:         user.ID,
		CreatedAt:  user.CreatedAt.Local().Format(time.DateTime),
		Username:   user.Username,
		Avatar:     resource.StaticPath,
		Gender:     user.Gender,
		Email:      user.Email,
		Mobile:     user.Mobile,
		CreateFrom: user.CreateFrom,
		Position:   user.Position,
		Projects:   Projects,
		Assignee:   assignee,
	}
}

type UserStatsResponse struct {
	TotalTasks      int64 `json:"total_tasks"`
	DoneTasks       int64 `json:"done_tasks"`
	InProgressTasks int64 `json:"in_progress_tasks"`
	Projects        int64 `json:"projects"`
	LastWeekTasks   int64 `json:"last_week_tasks"`
	ThisWeekTasks   int64 `json:"this_week_tasks"`
	LastMonthTasks  int64 `json:"last_month_tasks"`
	ThisMonthTasks  int64 `json:"this_month_tasks"`
}

type UserCalendarResponse struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
}
