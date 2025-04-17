package services

import (
	"errors"

	"server/internal/app/admin/dto"
	"server/internal/models"
	"server/internal/repositories"
)

type ProjectService struct {
	projectRepo       *repositories.ProjectRepo
	projectMemberRepo *repositories.ProjectMemberRepo
	userRepo          *repositories.UserRepo
	resourceRepo      *repositories.ResourceRepo
}

var projectService *ProjectService

func NewProjectService() *ProjectService {
	if projectService == nil {
		projectService = &ProjectService{
			projectRepo:       repositories.NewProjectRepo(),
			projectMemberRepo: repositories.NewProjectMemberRepo(),
			userRepo:          repositories.NewUserRepo(),
			resourceRepo:      repositories.NewResourceRepo(),
		}
	}
	return projectService
}

func (p *ProjectService) GetProjectList(request *dto.PageRequest) (*dto.ProjectPageResponse, error) {
	total := p.projectRepo.GetProjectCount()
	projects, err := p.projectRepo.GetProjectList(request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}

	var pageResponse *dto.ProjectPageResponse
	var projectResponse *dto.ProjectResponse
	projectList := []dto.ProjectResponse{}

	for _, project := range *projects {
		members, err := p.projectMemberRepo.GetMemberListByProjectId(project.ID)
		if err != nil {
			return nil, err
		}

		projectResponseWitheAvatars := []dto.ProjectResponseWitheAvatar{}
		for _, member := range *members {
			user, err := p.userRepo.GetUserById(member.UserID)
			if err != nil {
				return nil, err
			}
			resource, err := p.resourceRepo.GetResourceById(user.Avatar)
			if err != nil {
				return nil, err
			}
			projectResponseWitheAvatar := dto.ProjectResponseWitheAvatar{
				Avatar:        resource.StaticPath,
				ProjectMember: member,
			}

			projectResponseWitheAvatars = append(projectResponseWitheAvatars, projectResponseWitheAvatar)
		}

		projectResponse = projectResponse.Set(&project, &projectResponseWitheAvatars)
		projectList = append(projectList, *projectResponse)
	}

	return pageResponse.Set(total, request.Page, request.PageSize, projectList), nil
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

	var users []dto.MemberResponse
	for _, member := range *members {
		user, err := p.userRepo.GetUserById(member.UserID)
		if err != nil {
			return nil, err
		}

		var userResponse *dto.MemberResponse
		resource, err := p.resourceRepo.GetResourceById(user.Avatar)
		if err != nil {
			return nil, err
		}

		users = append(users, *userResponse.Set(user, member.Assignee, resource))
	}

	var projectResponse *dto.ProjectWithUserResponse
	return projectResponse.Set(project, users), nil
}

func (p *ProjectService) CreateProject(request *dto.ProjectCreateDto) error {
	var newProject models.Project

	newProject.Name = request.Name
	newProject.Desc = request.Desc

	if p.projectRepo.CheckProjectExistByName(newProject.Name) {
		return errors.New("项目名称已存在")
	}
	project, err := p.projectRepo.CreateProject(newProject)
	if err != nil {
		return err
	}
	if request.Members != nil {
		if err := p.projectMemberRepo.AddProjectMember(*request.Members, project.ID); err != nil {
			return err
		}
	}
	if request.Assignees != nil {
		var members []models.Member
		var ids []uint
		for _, member := range *request.Assignees {
			if p.projectMemberRepo.CheckProjectMemberExist(project.ID, member.UserID) {
				ids = append(ids, member.UserID)
			} else {
				members = append(members, member)
			}
		}

		if err := p.projectMemberRepo.AddProjectAssignee(members, project.ID); err != nil {
			return err
		}

		if err := p.projectMemberRepo.ChangeAssignees(ids, project.ID, true); err != nil {
			return err
		}
	}

	return nil
}

func (p *ProjectService) UpdateProject(request *dto.ProjectUpdateDto) error {
	var project models.Project

	project.ID = request.Id
	if request.Name != nil {
		project.Name = *request.Name
	}
	project.Desc = request.Desc

	err := p.projectRepo.UpdateProjectById(project)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProjectService) DeleteProject(request *dto.ProjectIdDto) error {
	if err := p.projectMemberRepo.DeleteProjectMember(request.Id); err != nil {
		return err
	}

	if err := p.projectRepo.DeleteProjectById(request.Id); err != nil {
		return err
	}

	return nil
}

func (p *ProjectService) SetProjectAssignee(request *dto.ProjectAssigneeDto) error {
	if !p.projectRepo.CheckProjectExistById(request.ProjectId) {
		return errors.New("项目不存在")
	}
	projectId := request.ProjectId
	members := request.Members
	value := *request.Value
	err := p.projectMemberRepo.ChangeAssignees(members, projectId, value)
	return err
}

func (p *ProjectService) AddProjectMember(request *dto.ProjectAddMemberDto) error {
	if !p.projectRepo.CheckProjectExistById(request.ProjectId) {
		return errors.New("项目不存在")
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

	if request.Assignees != nil {
		var members []models.Member
		var ids []uint
		for _, member := range *request.Assignees {
			if p.projectMemberRepo.CheckProjectMemberExist(projectId, member.UserID) {
				ids = append(ids, member.UserID)
			} else {
				members = append(members, member)
			}
		}

		if err := p.projectMemberRepo.AddProjectAssignee(members, projectId); err != nil {
			return err
		}

		if err := p.projectMemberRepo.ChangeAssignees(ids, projectId, true); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProjectService) RemoveProjectMember(request *dto.ProjectRemoveMemberDto) error {
	if !p.projectRepo.CheckProjectExistById(request.ProjectId) {
		return errors.New("项目不存在")
	}
	projectId := request.ProjectId
	members := request.Members
	err := p.projectMemberRepo.RemoveProjectMember(members, projectId)
	return err
}
