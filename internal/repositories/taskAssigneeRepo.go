package repositories

import (
	"server/internal/global"
	"server/internal/models"
	"server/internal/utils"

	"gorm.io/gorm"
)

type TaskAssigneeRepo struct {
	db *gorm.DB
}

var taskAssigneeRepo *TaskAssigneeRepo

func NewTaskAssigneeRepo() *TaskAssigneeRepo {
	if taskAssigneeRepo == nil {
		taskAssigneeRepo = &TaskAssigneeRepo{
			db: global.DB,
		}
	}
	return taskAssigneeRepo
}

func (t *TaskAssigneeRepo) GetTaskAssigneesByProjectIdAndTankId(projectId, taskId uint) (*[]models.TaskAssignee, error) {
	var taskAssignee []models.TaskAssignee
	err := t.db.Find(&taskAssignee, "project_id = ? and task_id = ?", projectId, taskId).Error
	return utils.HandleError(&taskAssignee, err)
}

func (t *TaskAssigneeRepo) GetTaskByUserIdLimt(uerId uint, page int, pageSize int) (*[]models.TaskAssignee, error) {
	var taskAssignee []models.TaskAssignee
	err := t.db.Limit(pageSize).Offset((page-1)*pageSize).Find(&taskAssignee, "user_id = ?", uerId).Error
	return utils.HandleError(&taskAssignee, err)
}

func (t *TaskAssigneeRepo) GetTaskByUserId(userId uint) (*[]models.TaskAssignee, error) {
	var taskAssignees []models.TaskAssignee
	err := t.db.Find(&taskAssignees, "user_id = ?", userId).Error
	return utils.HandleError(&taskAssignees, err)
}

func (t *TaskAssigneeRepo) GetTaskCountByUserId(userId uint) (int64, error) {
	var count int64
	err := t.db.Model(&models.TaskAssignee{}).Where("user_id = ?", userId).Count(&count).Error
	return count, err
}

func (t *TaskAssigneeRepo) CreateTaskAssignee(taskAssignee models.TaskAssignee) (*models.TaskAssignee, error) {
	err := t.db.Create(&taskAssignee).Error
	return utils.HandleError(&taskAssignee, err)
}

func (t *TaskAssigneeRepo) DeleteTaskAssigneeById(id uint) error {
	var taskAssignee models.TaskAssignee
	err := t.db.Where("task_id = ?", id).Delete(&taskAssignee).Error
	return err
}

func (t *TaskAssigneeRepo) AddTaskAssignee(taskAssignees []models.TaskAssignee, projectId uint, taskId uint) error {
	tx := t.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	for _, assignee := range taskAssignees {
		var count int64
		var modelTaskAssignee models.TaskAssignee
		if err := tx.Model(&modelTaskAssignee).
			Where("task_id = ? AND user_id = ?", taskId, assignee.UserID).
			Count(&count).Error; err != nil {
			tx.Rollback()
			return err
		}

		if count > 0 {
			continue
		}

		taskAssignee := models.TaskAssignee{
			ProjectID: projectId,
			TaskID:    taskId,
			UserID:    assignee.UserID,
			Username:  assignee.Username,
		}
		if err := tx.Create(&taskAssignee).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (t *TaskAssigneeRepo) RemoveAssignee(id uint, projectId uint, userId uint) error {
	var taskAssignee models.TaskAssignee
	if err := t.db.First(&taskAssignee, "task_id = ? and project_id = ? and user_id = ?", id, projectId, userId).Error; err != nil {
		return err
	}
	err := t.db.Delete(&taskAssignee, "task_id = ? and project_id = ? and user_id = ?", id, projectId, userId).Error
	return err
}

func (t *TaskAssigneeRepo) SearchTask(userId uint, projectId uint) (*[]models.TaskAssignee, error) {
	var taskAssignees []models.TaskAssignee
	var taskAssignee models.TaskAssignee
	err := t.db.Model(&taskAssignee).Where("user_id = ? AND project_id = ?", userId, projectId).Find(&taskAssignees).Error
	return utils.HandleError(&taskAssignees, err)
}

func (t *TaskAssigneeRepo) GetAllAssigneeIdByTaskId(taskID uint) (*[]uint, error) {
	assigneeIds := []uint{}
	var taskAssignee models.TaskAssignee
	err := t.db.Model(&taskAssignee).Where("task_id = ?", taskID).Pluck("user_id", &assigneeIds).Error
	return &assigneeIds, err
}
