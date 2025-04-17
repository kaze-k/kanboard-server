package repositories

import (
	"server/internal/global"
	"server/internal/models"
	"server/internal/utils"

	"gorm.io/gorm"
)

type ProjectMemberRepo struct {
	db *gorm.DB
}

var projectMemberRepo *ProjectMemberRepo

func NewProjectMemberRepo() *ProjectMemberRepo {
	if projectMemberRepo == nil {
		projectMemberRepo = &ProjectMemberRepo{
			db: global.DB,
		}
	}
	return projectMemberRepo
}

func (p *ProjectMemberRepo) CheckProjectMemberExist(projectId uint, userId uint) bool {
	var projectMember models.ProjectMember
	err := p.db.First(&projectMember, "project_id = ? and user_id = ?", projectId, userId).Error
	return err == nil
}

func (p *ProjectMemberRepo) ChangeAssignees(userIds []uint, projectId uint, assignee bool) error {
	tx := p.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	var projectMember models.ProjectMember
	for _, userId := range userIds {
		if err := tx.Model(&projectMember).Where("user_id = ? and project_id = ?", userId, projectId).First(&projectMember).Update("assignee", assignee).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (p *ProjectMemberRepo) AddProjectAssignee(members []models.Member, projectId uint) error {
	tx := p.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	for _, member := range members {
		projectMember := models.ProjectMember{
			ProjectID: projectId,
			UserID:    member.UserID,
			Username:  member.Username,
			Assignee:  true,
		}
		if err := tx.Find(&projectMember).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Create(&projectMember).Error; err != nil {
			tx.Rollback()
			return err

		}
	}
	return tx.Commit().Error
}

func (p *ProjectMemberRepo) AddProjectMember(members []models.Member, projectId uint) error {
	tx := p.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	for _, member := range members {
		projectMember := models.ProjectMember{
			ProjectID: projectId,
			UserID:    member.UserID,
			Username:  member.Username,
		}
		if err := tx.Find(&projectMember).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Create(&projectMember).Error; err != nil {
			tx.Rollback()
			return err

		}
	}
	return tx.Commit().Error
}

func (p *ProjectMemberRepo) RemoveProjectMember(userIds []uint, projectId uint) error {
	tx := p.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	var projectMember models.ProjectMember
	for _, userId := range userIds {
		if err := tx.First(&projectMember, "user_id = ? and project_id = ?", userId, projectId).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Delete(&projectMember, "user_id = ? and project_id = ?", userId, projectId).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (p *ProjectMemberRepo) RemoveProjectMemberById(userId uint, projectId uint) error {
	var projectMember models.ProjectMember
	if err := p.db.First(&projectMember, "user_id = ? and project_id = ?", userId, projectId).Error; err != nil {
		return err
	}
	err := p.db.Delete(&projectMember, "user_id = ? and project_id = ?", userId, projectId).Error
	return err
}

func (p *ProjectMemberRepo) DeleteProjectMember(projectId uint) error {
	var projectMember models.ProjectMember
	if err := p.db.First(&projectMember, "project_id = ?", projectId).Error; err != nil {
		return err
	}
	err := p.db.Delete(&projectMember, "project_id = ?", projectId).Error
	return err
}

func (p *ProjectMemberRepo) GetMemberCountByProjectId(projectId uint) int64 {
	count := p.db.Model(&models.ProjectMember{}).Where("project_id = ?", projectId).RowsAffected
	return count
}

func (p *ProjectMemberRepo) GetMemberListByProjectId(projectId uint) (*[]models.ProjectMember, error) {
	var projectMembers []models.ProjectMember
	err := p.db.Order("assignee desc").Find(&projectMembers, "project_id = ?", projectId).Error
	return utils.HandleError(&projectMembers, err)
}

func (p *ProjectMemberRepo) GetProjectByUserId(userId uint) ([]models.ProjectMember, error) {
	projectMember := []models.ProjectMember{}
	err := p.db.Find(&projectMember, "user_id = ?", userId).Error
	return projectMember, err
}

func (p *ProjectMemberRepo) CheckAssignee(projectId uint, userId uint) bool {
	var projectMember models.ProjectMember
	err := p.db.First(&projectMember, "project_id = ? and user_id = ? and assignee = ?", projectId, userId, true).Error
	return err == nil
}

func (p *ProjectMemberRepo) GetAllMemberIdByProjectId(projectId uint) (*[]uint, error) {
	userIDData := []uint{}
	var projectMember models.ProjectMember
	err := p.db.Model(&projectMember).Where("project_id = ?", projectId).Pluck("user_id", &userIDData).Error
	return &userIDData, err
}

func (p *ProjectMemberRepo) GetAllAssigneeIdByProjectId(projectId uint) (*[]uint, error) {
	userIDData := []uint{}
	var projectMember models.ProjectMember
	err := p.db.Model(&projectMember).Where("project_id = ? and assignee = ?", projectId, true).Pluck("user_id", &userIDData).Error
	return &userIDData, err
}
