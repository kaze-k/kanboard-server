package repositories

import (
	"time"

	"server/internal/global"
	"server/internal/models"
	"server/internal/utils"

	"gorm.io/gorm"
)

type ProjectRepo struct {
	db *gorm.DB
}

var projectRepo *ProjectRepo

func NewProjectRepo() *ProjectRepo {
	if projectRepo == nil {
		projectRepo = &ProjectRepo{
			db: global.DB,
		}
	}
	return projectRepo
}

func (u *ProjectRepo) CheckProjectExistByName(name string) bool {
	var project models.Project
	count := u.db.Find(&project, "name = ?", name).RowsAffected
	return count > 0
}

func (u *ProjectRepo) CheckProjectExistById(id uint) bool {
	var project models.Project
	count := u.db.Find(&project, "id = ?", id).RowsAffected
	return count > 0
}

func (p *ProjectRepo) GetProjectById(id uint) (*models.Project, error) {
	var project models.Project
	err := p.db.Find(&project, "id = ?", id).Error
	return utils.HandleError(&project, err)
}

func (p *ProjectRepo) CreateProject(project models.Project) (*models.Project, error) {
	err := p.db.Create(&project).Error
	return utils.HandleError(&project, err)
}

func (p *ProjectRepo) UpdateProjectById(project models.Project) error {
	err := p.db.Updates(&project).Error
	return err
}

func (p *ProjectRepo) DeleteProjectById(id uint) error {
	var project models.Project
	if err := p.db.First(&project, id).Error; err != nil {
		return err
	}
	err := p.db.Delete(&project, "id = ?", id).Error
	return err
}

func (p *ProjectRepo) GetProjectCount() int64 {
	var project models.Project
	var count int64
	p.db.Model(&project).Count(&count)
	return count
}

func (p *ProjectRepo) GetProjectList(page, pageSize int) (*[]models.Project, error) {
	var projects []models.Project
	err := p.db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&projects).Error
	return utils.HandleError(&projects, err)
}

func (p *ProjectRepo) GetAllProjectCreatedAt() (*[]time.Time, error) {
	createdAtData := []time.Time{}
	var project models.Project
	err := p.db.Model(&project).Pluck("created_at", &createdAtData).Error
	return &createdAtData, err
}
