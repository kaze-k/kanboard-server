package repositories

import (
	"server/internal/global"
	"server/internal/models"
	"server/internal/utils"

	"gorm.io/gorm"
)

type ResourceRepo struct {
	db *gorm.DB
}

var resourceRepo *ResourceRepo

func NewResourceRepo() *ResourceRepo {
	if resourceRepo == nil {
		resourceRepo = &ResourceRepo{
			db: global.DB,
		}
	}
	return resourceRepo
}

func (r *ResourceRepo) AddResource(md5 string, filetype string, filePath string, staticPath string) (*models.Resource, error) {
	var resource models.Resource
	resource.MD5 = md5
	resource.FileType = filetype
	resource.FilePath = filePath
	resource.StaticPath = staticPath
	err := r.db.Create(&resource).First(&resource).Error
	return utils.HandleError(&resource, err)
}

func (r *ResourceRepo) GetResourceByMd5(md5 string) (*models.Resource, error) {
	var resource models.Resource
	err := r.db.Find(&resource, "md5 = ?", md5).Error
	return utils.HandleError(&resource, err)
}

func (r *ResourceRepo) GetResourceById(id uint) (*models.Resource, error) {
	var resource models.Resource
	err := r.db.Find(&resource, "id = ?", id).Error
	return utils.HandleError(&resource, err)
}
