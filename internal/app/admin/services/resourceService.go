package services

import (
	"server/internal/repositories"
)

type ResourceService struct {
	resourceRepo *repositories.ResourceRepo
}

var resourceService *ResourceService

func NewResourceService() *ResourceService {
	if resourceService == nil {
		resourceService = &ResourceService{
			resourceRepo: repositories.NewResourceRepo(),
		}
	}
	return resourceService
}
