package services

import (
	"time"

	"server/internal/app/admin/dto"
	"server/internal/constant"
	"server/internal/global"
	"server/internal/repositories"
	"server/pkg/captcha"
)

type PublicService struct {
	captcha     *captcha.Captcha
	projectRepo *repositories.ProjectRepo
	taskRepo    *repositories.TaskRepo
	userRepo    *repositories.UserRepo
}

var publicService *PublicService

func NewPublicService() *PublicService {
	if publicService == nil {
		redisStore := captcha.NewRedisStore(constant.ADMIN_CAPTCHA, 1*time.Minute, global.Redis)
		publicService = &PublicService{
			captcha:     captcha.NewCaptcha(redisStore),
			projectRepo: repositories.NewProjectRepo(),
			taskRepo:    repositories.NewTaskRepo(),
			userRepo:    repositories.NewUserRepo(),
		}
	}

	return publicService
}

func (p *PublicService) GetCaptcha() (*dto.CaptchaResponse, error) {
	id, b64s, _, err := p.captcha.Generate()
	if err != nil {
		return nil, err
	}
	response := &dto.CaptchaResponse{
		Captcha: b64s,
		Id:      id,
	}
	return response, nil
}

func (p *PublicService) VerifyCaptcha(id, answer string) bool {
	return p.captcha.Verify(id, answer)
}

func (p *PublicService) GetRecentTasks() ([]dto.TaskRecentResponse, error) {
	recentTasks, err := p.taskRepo.GetRecentTasks()
	if err != nil {
		return nil, err
	}
	data := []dto.TaskRecentResponse{}
	for index, task := range *recentTasks {
		var taskRecentResponse dto.TaskRecentResponse
		project, err := p.projectRepo.GetProjectById(task.ProjectID)
		if err != nil {
			return nil, err
		}
		taskRecentResponse.Set(index, &task, project)
		data = append(data, taskRecentResponse)
	}
	return data, nil
}

func (p *PublicService) GetDashboard() (*dto.DashboardResponse, error) {
	projectCount := p.projectRepo.GetProjectCount()

	createAtData, err := p.projectRepo.GetAllProjectCreatedAt()
	if err != nil {
		return nil, err
	}

	projectCreatedAt := []string{}
	for _, createAt := range *createAtData {
		projectCreatedAt = append(projectCreatedAt, createAt.Local().Format(time.DateTime))
	}

	taskCount := p.taskRepo.GetTaskCount()
	taskUndoCount := p.taskRepo.GetAllUndoTaskCount()
	taskDoneCount := p.taskRepo.GetAllDoneTaskCount()
	taskHighCount := p.taskRepo.GetAllHighTaskCount()
	taskMediumCount := p.taskRepo.GetAllMediumTaskCount()
	taskLowCount := p.taskRepo.GetAllLowTaskCount()

	taskCreateAtData, err := p.taskRepo.GetAllTaskCreatedAt()
	if err != nil {
		return nil, err
	}

	taskCreateAt := []string{}
	for _, createAt := range *taskCreateAtData {
		taskCreateAt = append(taskCreateAt, createAt.Local().Format(time.DateTime))
	}

	userCount := p.userRepo.GetUserCount()
	userMaleCount := p.userRepo.GetMaleCount()
	userFemaleCount := p.userRepo.GetFemaleCount()
	userCreateFromKanboardCount := p.userRepo.GetCreateFromKanboardCount()
	userCreateFromAdmin := p.userRepo.GetCreateFromAdminCount()
	userAdminCount := p.userRepo.GetAdminCount()
	userLoginableCount := p.userRepo.GetLoginableCount()
	userUnLoginableCount := p.userRepo.GetUnLoginableCount()

	userCreateAtData, err := p.userRepo.GetAllUserCreateAt()
	if err != nil {
		return nil, err
	}

	userCreateAt := []string{}
	for _, createAt := range *userCreateAtData {
		userCreateAt = append(userCreateAt, createAt.Local().Format(time.DateTime))
	}

	project := dto.DashboardProject{}
	project.Count = projectCount
	project.CreatedAt = projectCreatedAt

	var task dto.DashboardTask
	task.Count = taskCount
	task.Undo = taskUndoCount
	task.Done = taskDoneCount
	task.High = taskHighCount
	task.Medium = taskMediumCount
	task.Low = taskLowCount
	task.CreatedAt = taskCreateAt

	var user dto.DashboardUser
	user.Count = userCount
	user.Male = userMaleCount
	user.Female = userFemaleCount
	user.CreateFromKanboard = userCreateFromKanboardCount
	user.CreateFromAdmin = userCreateFromAdmin
	user.Admin = userAdminCount
	user.Loginable = userLoginableCount
	user.Unloginable = userUnLoginableCount
	user.CreatedAt = userCreateAt

	response := &dto.DashboardResponse{
		Project: project,
		Task:    task,
		User:    user,
	}

	return response, nil
}
