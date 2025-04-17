package services

import (
	"errors"
	"strconv"
	"time"

	"server/internal/app/kanboard/dto"
	"server/internal/constant"
	"server/internal/global"
	"server/internal/models"
	"server/internal/repositories"
	"server/pkg/crypto"
)

type UserService struct {
	userRepo          *repositories.UserRepo
	resourceRepo      *repositories.ResourceRepo
	projectMemberRepo *repositories.ProjectMemberRepo
	projectRepo       *repositories.ProjectRepo
	taskRepo          *repositories.TaskRepo
	taskAssigneeRepo  *repositories.TaskAssigneeRepo
}

var userService *UserService

func NewUserService() *UserService {
	if userService == nil {
		userService = &UserService{
			userRepo:          repositories.NewUserRepo(),
			resourceRepo:      repositories.NewResourceRepo(),
			projectMemberRepo: repositories.NewProjectMemberRepo(),
			projectRepo:       repositories.NewProjectRepo(),
			taskRepo:          repositories.NewTaskRepo(),
			taskAssigneeRepo:  repositories.NewTaskAssigneeRepo(),
		}
	}
	return userService
}

func (u *UserService) checkUserInfoExist(username string, email, mobile *string) (bool, error) {
	isUsernameExist := u.userRepo.CheckUserExistByName(username)
	if isUsernameExist {
		return isUsernameExist, errors.New("用户名已存在")
	}

	if email != nil {
		isEmailExist := u.userRepo.CheckEmailExist(*email)
		if isEmailExist {
			return isEmailExist, errors.New("邮箱已存在")
		}
	}

	if mobile != nil {
		isMobileExist := u.userRepo.CheckMobileExist(*mobile)
		if isMobileExist {
			return isMobileExist, errors.New("手机号已存在")
		}
	}

	return false, nil
}

func GenerateToken(id uint, username string) (string, error) {
	token, err := crypto.GenerateJWTToAdmin(id, username)
	global.Redis.Set(constant.KANBOARD_TOKEN, strconv.Itoa(int(id)), token, constant.JWTConfig.KanboardTokenExpiration)
	return token, err
}

func (u *UserService) Login(request dto.UserLoginRequest) (*dto.UserResponse, *string, error) {
	user, err := u.userRepo.GetUserByName(request.Username)
	comp := crypto.CheckPasswordHash(user.Password, request.Password)

	if !comp {
		return nil, nil, errors.New("用户名或密码错误")
	}

	token, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, nil, errors.New("token生成失败")
	}

	var userResponse *dto.UserResponse
	resource, err := u.resourceRepo.GetResourceById(user.Avatar)
	if err != nil {
		return nil, nil, err
	}
	projectMembers, err := u.projectMemberRepo.GetProjectByUserId(user.ID)
	if err != nil {
		return nil, nil, err
	}
	projects := []dto.ProjectsWithIdAndAssignee{}
	for _, projectMember := range projectMembers {
		projectModel, err := u.projectRepo.GetProjectById(projectMember.ProjectID)
		if err != nil {
			return nil, nil, err
		}
		project := dto.ProjectsWithIdAndAssignee{
			ProjectID:   projectMember.ProjectID,
			Assignee:    projectMember.Assignee,
			ProjectName: projectModel.Name,
			JoinedAt:    projectMember.JoinedAt.Local().Format(time.DateTime),
		}
		projects = append(projects, project)
	}

	return userResponse.Set(user, resource, projects, nil), &token, nil
}

func (u *UserService) RegisterKanboard(request dto.UserRegisterRequest) (*dto.UserResponse, error) {
	isExist, err := u.checkUserInfoExist(request.Username, request.Email, request.Mobile)
	if isExist {
		return nil, err
	}

	var createUser models.User

	createUser.Username = request.Username
	createUser.Password = request.Password
	if request.Email != nil {
		createUser.Email = *request.Email
	}
	if request.Mobile != nil {
		createUser.Mobile = *request.Mobile
	}
	createUser.Gender = request.Gender
	createUser.CreateFrom = constant.KANBOARD
	createUser.Loginable = true

	user, err := u.userRepo.CreateUser(createUser)

	var userResponse *dto.UserResponse
	resource, err := u.resourceRepo.GetResourceById(user.Avatar)
	if err != nil {
		return nil, err
	}
	projectMembers, err := u.projectMemberRepo.GetProjectByUserId(user.ID)
	if err != nil {
		return nil, err
	}
	projects := []dto.ProjectsWithIdAndAssignee{}
	for _, projectMember := range projectMembers {
		projectModel, err := u.projectRepo.GetProjectById(projectMember.ProjectID)
		if err != nil {
			return nil, err
		}
		project := dto.ProjectsWithIdAndAssignee{
			ProjectID:   projectMember.ProjectID,
			Assignee:    projectMember.Assignee,
			ProjectName: projectModel.Name,
		}
		projects = append(projects, project)
	}

	return userResponse.Set(user, resource, projects, nil), err
}

func (u *UserService) GetUserInfo(userIdDto dto.UserIDRequest) (*dto.UserResponse, error) {
	user, err := u.userRepo.GetUserById(userIdDto.ID)
	var userResponse *dto.UserResponse
	resource, err := u.resourceRepo.GetResourceById(user.Avatar)
	if err != nil {
		return nil, err
	}
	projectMembers, err := u.projectMemberRepo.GetProjectByUserId(user.ID)
	if err != nil {
		return nil, err
	}
	projects := []dto.ProjectsWithIdAndAssignee{}
	for _, projectMember := range projectMembers {
		projectModel, err := u.projectRepo.GetProjectById(projectMember.ProjectID)
		if err != nil {
			return nil, err
		}
		project := dto.ProjectsWithIdAndAssignee{
			ProjectID:   projectMember.ProjectID,
			Assignee:    projectMember.Assignee,
			ProjectName: projectModel.Name,
		}
		projects = append(projects, project)
	}
	return userResponse.Set(user, resource, projects, nil), err
}

func (u *UserService) UpdateInfo(request dto.UserUpdateRequest) (*dto.UserResponse, error) {
	updateData := make(map[string]any)
	if request.Email != nil {
		email := *request.Email
		if isExist := u.userRepo.CheckEmailExist(email); isExist {
			return nil, errors.New("邮箱已存在")
		}
		updateData["email"] = request.Email
	}

	if request.Mobile != nil {
		mobile := *request.Mobile
		if isMobileExist := u.userRepo.CheckMobileExist(mobile); isMobileExist {
			return nil, errors.New("手机号已存在")
		}
		updateData["mobile"] = request.Mobile
	}

	if request.Avatar != nil {
		updateData["avatar"] = request.Avatar
	}

	user, err := u.userRepo.UpdateUserById(updateData, request.ID)
	if err != nil {
		return nil, err
	}

	var userResponse *dto.UserResponse
	resource, err := u.resourceRepo.GetResourceById(user.Avatar)
	if err != nil {
		return nil, err
	}
	projectMembers, err := u.projectMemberRepo.GetProjectByUserId(user.ID)
	if err != nil {
		return nil, err
	}
	projects := []dto.ProjectsWithIdAndAssignee{}
	for _, projectMember := range projectMembers {
		projectModel, err := u.projectRepo.GetProjectById(projectMember.ProjectID)
		if err != nil {
			return nil, err
		}
		project := dto.ProjectsWithIdAndAssignee{
			ProjectID:   projectMember.ProjectID,
			Assignee:    projectMember.Assignee,
			ProjectName: projectModel.Name,
		}
		projects = append(projects, project)
	}

	return userResponse.Set(user, resource, projects, nil), nil
}

func (u *UserService) UpdatePassword(request dto.UserUpdatePasswordRequest) error {
	user, err := u.userRepo.GetUserById(request.ID)
	comp := crypto.CheckPasswordHash(user.Password, request.Current)
	if !comp {
		return errors.New("密码错误")
	}

	if err := u.userRepo.UpdatePasswordById(request.New, request.ID); err != nil {
		return err
	}

	return err
}

func (u *UserService) GetStatistics(userId uint) (*dto.UserStatsResponse, error) {
	totalTasks, err := u.taskAssigneeRepo.GetTaskCountByUserId(userId)
	if err != nil {
		return nil, err
	}
	projects, err := u.projectMemberRepo.GetProjectByUserId(userId)
	if err != nil {
		return nil, err
	}
	var totalProjects int64 = int64(len(projects))

	var doneTasks int64 = 0
	var inProgressTasks int64 = 0
	var lastWeekDoneTasks int64 = 0
	var thisWeekDoneTasks int64 = 0
	var lastMonthDoneTasks int64 = 0
	var thisMonthDoneTasks int64 = 0
	taskAssignees, err := u.taskAssigneeRepo.GetTaskByUserId(userId)
	for _, taskAssignee := range *taskAssignees {
		if u.taskRepo.CheckTaskDoneById(taskAssignee.TaskID) {
			doneTasks++

			task, err := u.taskRepo.GetTaskById(taskAssignee.TaskID)
			if err != nil {
				return nil, err
			}

			now := time.Now()
			lastMonth := now.AddDate(0, -1, 0)
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7 // 周日为 0，转换为 7
			}
			thisWeekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -weekday+1)
			lastWeekStart := thisWeekStart.AddDate(0, 0, -7)
			lastWeekEnd := thisWeekStart.AddDate(0, 0, -1).Add(time.Hour*23 + time.Minute*59 + time.Second*59)

			if task.CreatedAt.Year() == now.Year() && task.CreatedAt.Month() == now.Month() {
				thisMonthDoneTasks++
			}
			if task.CreatedAt.Year() == lastMonth.Year() && task.CreatedAt.Month() == lastMonth.Month() {
				lastMonthDoneTasks++
			}
			if !task.CreatedAt.Before(lastWeekStart) && !task.CreatedAt.After(lastWeekEnd) {
				lastWeekDoneTasks++
			}
			if !task.CreatedAt.Before(thisWeekStart) && !task.CreatedAt.After(now) {
				thisWeekDoneTasks++
			}
			continue
		}
		if u.taskRepo.CheckTaskInProgressById(taskAssignee.TaskID) {
			inProgressTasks++
		}
	}

	userStatsResponse := &dto.UserStatsResponse{
		TotalTasks:      totalTasks,
		DoneTasks:       doneTasks,
		InProgressTasks: inProgressTasks,
		Projects:        totalProjects,
		LastWeekTasks:   lastWeekDoneTasks,
		ThisWeekTasks:   thisWeekDoneTasks,
		LastMonthTasks:  lastMonthDoneTasks,
		ThisMonthTasks:  thisMonthDoneTasks,
	}
	return userStatsResponse, nil
}

func (u *UserService) GetCalendar(userId uint) ([]dto.UserCalendarResponse, error) {
	taskAssignees, err := u.taskAssigneeRepo.GetTaskByUserId(userId)
	if err != nil {
		return nil, err
	}

	responses := []dto.UserCalendarResponse{}
	for _, taskAssignee := range *taskAssignees {
		task, err := u.taskRepo.GetDueDateByTaskId(taskAssignee.TaskID)
		if err != nil {
			continue
		}

		response := dto.UserCalendarResponse{
			Id:    task.ID,
			Title: task.Title,
			Date:  task.DueDate.Local().Format(time.DateTime),
		}
		responses = append(responses, response)
	}

	return responses, nil
}
