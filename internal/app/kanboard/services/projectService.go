package services

import (
	"errors"

	"server/internal/app/kanboard/dto"
	"server/internal/models"
	"server/internal/repositories"
)

type ProjectService struct {
	projectRepo       *repositories.ProjectRepo
	projectMemberRepo *repositories.ProjectMemberRepo
	userRepo          *repositories.UserRepo
	taskRepo          *repositories.TaskRepo
	resourceRepo      *repositories.ResourceRepo
}

var projectService *ProjectService

func NewProjectService() *ProjectService {
	if projectService == nil {
		projectService = &ProjectService{
			projectRepo:       repositories.NewProjectRepo(),
			projectMemberRepo: repositories.NewProjectMemberRepo(),
			userRepo:          repositories.NewUserRepo(),
			taskRepo:          repositories.NewTaskRepo(),
			resourceRepo:      repositories.NewResourceRepo(),
		}
	}
	return projectService
}

func (p *ProjectService) GetProjectById(id uint) (*dto.ProjectWithUserResponse, error) {
	project, err := p.projectRepo.GetProjectById(id)
	if err != nil {
		return nil, err
	}
	if project.ID == 0 {
		return nil, errors.New("项目不存在")
	}
	members, err := p.projectMemberRepo.GetMemberListByProjectId(project.ID)
	if err != nil {
		return nil, err
	}

	var users []dto.UserResponse
	for _, member := range *members {
		user, err := p.userRepo.GetUserById(member.UserID)
		if err != nil {
			return nil, err
		}

		var userResponse *dto.UserResponse
		resource, err := p.resourceRepo.GetResourceById(user.Avatar)
		if err != nil {
			return nil, err
		}

		assignee := p.projectMemberRepo.CheckAssignee(project.ID, member.UserID)

		users = append(users, *userResponse.Set(user, resource, nil, &assignee))
	}

	taskTotal, _ := p.taskRepo.GetTaskCountByProjectId(project.ID)
	taskInProgress := p.taskRepo.GetTaskInProgressCountByProjectId(project.ID)
	taskDone := p.taskRepo.GetTaskDoneCountByProjectId(project.ID)
	taskHighPriority := p.taskRepo.GetTaskHighPriorityCountByProjectId(project.ID)
	taskMediumPriority := p.taskRepo.GetTaskMediumPriorityCountByProjectId(project.ID)
	taskLowPriority := p.taskRepo.GetTaskLowPriorityCountByProjectId(project.ID)

	statistics := &dto.Statistics{
		TaskTotal:      taskTotal,
		TaskInProgress: taskInProgress,
		TaskDone:       taskDone,
		TaskHigh:       taskHighPriority,
		TaskMedium:     taskMediumPriority,
		TaskLow:        taskLowPriority,
	}

	var projectResponse *dto.ProjectWithUserResponse
	return projectResponse.Set(project, users, statistics), nil
}

func (p *ProjectService) AddProjectMember(request *dto.ProjectAddMemberDto) error {
	if !p.projectRepo.CheckProjectExistById(request.ProjectId) {
		return errors.New("项目不存在")
	}
	if !p.projectMemberRepo.CheckAssignee(request.ProjectId, request.UserId) {
		return errors.New("没有权限")
	}
	projectId := request.ProjectId
	var members []models.Member
	for _, member := range request.Members {
		if !p.projectMemberRepo.CheckProjectMemberExist(projectId, member.UserID) {
			members = append(members, member)
		}
	}
	if err := p.projectMemberRepo.AddProjectMember(members, projectId); err != nil {
		return err
	}

	return nil
}

func (p *ProjectService) RemoveProjectMember(request *dto.ProjectRemoveMemberDto) error {
	if !p.projectRepo.CheckProjectExistById(request.ProjectId) {
		return errors.New("项目不存在")
	}
	if !p.projectMemberRepo.CheckAssignee(request.ProjectId, request.UserId) {
		return errors.New("没有权限")
	}
	projectId := request.ProjectId
	memberId := request.MemberId
	err := p.projectMemberRepo.RemoveProjectMemberById(memberId, projectId)
	return err
}

func (p *ProjectService) GetProjectMembers(projectId uint, userId uint) ([]models.ProjectMember, error) {
	if !p.projectRepo.CheckProjectExistById(projectId) {
		return nil, errors.New("项目不存在")
	}
	if !p.projectMemberRepo.CheckProjectMemberExist(projectId, userId) {
		return nil, errors.New("没有权限")
	}
	_members, err := p.projectMemberRepo.GetMemberListByProjectId(projectId)
	if err != nil {
		return nil, err
	}
	members := []models.ProjectMember{}
	if _members != nil {
		members = *_members
	}
	return members, nil
}

func (p *ProjectService) GetMembers(projectId uint, useId uint) ([]models.Member, error) {
	isAssignee := p.projectMemberRepo.CheckAssignee(projectId, useId)
	if !isAssignee {
		return nil, errors.New("没有权限")
	}
	users, err := p.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var member models.Member
	var members []models.Member

	for _, user := range *users {
		member = member.Set(user)
		members = append(members, member)
	}

	return members, nil
}
