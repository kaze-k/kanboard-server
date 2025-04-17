package dto

import (
	"time"

	"server/internal/constant"
	"server/internal/models"
)

type UserIDRequest struct {
	ID uint `json:"id" binding:"required" uri:"id"`
}

type UserLoginRequest struct {
	Username      string `json:"username" binding:"required" from:"username"`
	Password      string `json:"password" binding:"required" from:"password"`
	CaptchaID     string `json:"captchaId" binding:"required" from:"captchaId"`
	CaptchaAnswer string `json:"captchaAnswer" binding:"required" from:"captchaAnswer"`
}

type UserRegisterRequest struct {
	Username string          `json:"username" binding:"required" form:"username"`
	Password string          `json:"password" binding:"required" form:"password"`
	Email    string          `json:"email" form:"email"`
	Mobile   string          `json:"mobile" form:"mobile"`
	Gender   constant.Gender `json:"gender" form:"gender"`
}

type UserCreateRequest struct {
	Username  string          `json:"username" binding:"required" form:"username"`
	Password  string          `json:"password" binding:"required" form:"password"`
	Email     *string         `json:"email" form:"email"`
	Mobile    *string         `json:"mobile" form:"mobile"`
	Gender    constant.Gender `json:"gender" form:"gender"`
	Avatar    *uint           `json:"avatar" form:"avatar"`
	IsAdmin   bool            `json:"is_admin" form:"is_admin"`
	Loginable bool            `json:"loginable" form:"loginable"`
	Position  *string         `json:"position" form:"position"`
}

type UserUpdateRequest struct {
	ID        uint             `json:"id" binding:"required" form:"id"`
	Avatar    *uint            `json:"avatar" form:"avatar"`
	Email     *string          `json:"email" form:"email"`
	Mobile    *string          `json:"mobile" form:"mobile"`
	Gender    *constant.Gender `json:"gender" form:"gender"`
	IsAdmin   *bool            `json:"is_admin" form:"is_admin"`
	Loginable *bool            `json:"loginable" form:"loginable"`
	Position  *string          `json:"position" form:"position"`
}

type UserSearchRequest struct {
	ID         *uint            `json:"id" form:"id"`
	Username   *string          `json:"username" form:"username"`
	CreateFrom *constant.From   `json:"create_from" form:"create_from"`
	IsAdmin    *bool            `json:"is_admin" form:"is_admin"`
	Loginable  *bool            `json:"loginable" form:"loginable"`
	Position   *string          `json:"position" form:"position"`
	Gender     *constant.Gender `json:"gender" form:"gender"`
	PageRequest
}

type UserChangePasswordRequest struct {
	ID          uint   `json:"id" binding:"required" form:"id"`
	NewPassword string `json:"newPassword" binding:"required" form:"newPassword"`
}

type UserDataResponse struct {
	ID         uint            `json:"id"`
	Username   string          `json:"username"`
	CreateFrom constant.From   `json:"create_from"`
	Loginable  bool            `json:"loginable"`
	IsAdmin    bool            `json:"is_admin"`
	Position   string          `json:"position"`
	Gender     constant.Gender `json:"gender"`
}

type UserPageResponse struct {
	Total     int                `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
	TotalPage int                `json:"total_page"`
	Data      []UserDataResponse `json:"data"`
}

func (r *UserPageResponse) Set(total int64, page int, pageSize int, users *[]models.User) *UserPageResponse {
	var data []UserDataResponse
	for _, user := range *users {
		userDataRespnose := UserDataResponse{
			ID:         user.ID,
			Username:   user.Username,
			CreateFrom: user.CreateFrom,
			Loginable:  user.Loginable,
			IsAdmin:    user.IsAdmin,
			Position:   user.Position,
			Gender:     user.Gender,
		}
		data = append(data, userDataRespnose)
	}

	if len(data) == 0 {
		data = []UserDataResponse{}
	}

	var userPageResponse UserPageResponse

	userPageResponse.Total = int(total)
	userPageResponse.Page = page
	userPageResponse.PageSize = pageSize
	totalPage := int(float64(total) / float64(pageSize))
	if total%int64(pageSize) != 0 {
		totalPage++
	}
	userPageResponse.TotalPage = totalPage
	userPageResponse.Data = data

	return &userPageResponse
}

type UserResponse struct {
	CreateAt   string          `json:"create_at"`
	UpdateAt   string          `json:"update_at"`
	ID         uint            `json:"id"`
	Username   string          `json:"username"`
	Avatar     string          `json:"avatar"`
	Gender     constant.Gender `json:"gender"`
	Email      string          `json:"email"`
	Mobile     string          `json:"mobile"`
	CreateFrom constant.From   `json:"create_from"`
	Loginable  bool            `json:"loginable"`
	IsAdmin    bool            `json:"is_admin"`
	Position   string          `json:"position"`
}

func (r *UserResponse) Set(user *models.User, resource *models.Resource) *UserResponse {
	return &UserResponse{
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
	}
}
