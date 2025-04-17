package repositories

import (
	"time"

	"server/internal/constant"
	"server/internal/global"
	"server/internal/models"
	"server/internal/utils"

	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

var taskRepo *TaskRepo

func NewTaskRepo() *TaskRepo {
	if taskRepo == nil {
		taskRepo = &TaskRepo{
			db: global.DB,
		}
	}
	return taskRepo
}

func (t *TaskRepo) GetTaskByProjectIdLimt(projectId uint, page int, pageSize int) (*[]models.Task, error) {
	var task []models.Task
	err := t.db.Limit(pageSize).Offset((page-1)*pageSize).Find(&task, "project_id = ?", projectId).Error
	return utils.HandleError(&task, err)
}

func (t *TaskRepo) GetTaskByProjectId(projectId uint) (*[]models.Task, error) {
	var task []models.Task
	err := t.db.Find(&task, "project_id = ?", projectId).Error
	return utils.HandleError(&task, err)
}

func (t *TaskRepo) GetTaskCountByProjectId(projectId uint) (int64, error) {
	var count int64
	err := t.db.Model(&models.Task{}).Where("project_id = ?", projectId).Count(&count).Error
	return count, err
}

func (t *TaskRepo) GetTaskCount() int64 {
	var task models.Task
	var count int64
	t.db.Model(&task).Count(&count)
	return count
}

func (t *TaskRepo) GetAllUndoTaskCount() int64 {
	var task models.Task
	var count int64
	t.db.Model(&task).Where("status = ?", constant.TASK_STATUS_UNDO).Count(&count)
	return count
}

func (t *TaskRepo) GetAllDoneTaskCount() int64 {
	var task models.Task
	var count int64
	t.db.Model(&task).Where("status = ?", constant.TASK_STATUS_DONE).Count(&count)
	return count
}

func (t *TaskRepo) GetAllHighTaskCount() int64 {
	var task models.Task
	var count int64
	t.db.Model(&task).Where("priority = ?", constant.TASK_PRIORITY_HIGH).Count(&count)
	return count
}

func (t *TaskRepo) GetAllMediumTaskCount() int64 {
	var task models.Task
	var count int64
	t.db.Model(&task).Where("priority = ?", constant.TASK_PRIORITY_MEDIUM).Count(&count)
	return count
}

func (t *TaskRepo) GetAllLowTaskCount() int64 {
	var task models.Task
	var count int64
	t.db.Model(&task).Where("priority = ?", constant.TASK_PRIORITY_LOW).Count(&count)
	return count
}

func (t *TaskRepo) GetAllTaskCreatedAt() (*[]time.Time, error) {
	createdAtData := []time.Time{}
	var task models.Task
	err := t.db.Model(&task).Pluck("created_at", &createdAtData).Error
	return &createdAtData, err
}

func (t *TaskRepo) GetRecentTasks() (*[]models.Task, error) {
	var task []models.Task
	err := t.db.Order("created_at DESC").Limit(6).Find(&task).Error
	return utils.HandleError(&task, err)
}

func (t *TaskRepo) GetTaskById(id uint) (*models.Task, error) {
	var task models.Task
	err := t.db.First(&task, id).Error
	return utils.HandleError(&task, err)
}

func (t *TaskRepo) GetTaskByIdAndProjectId(id uint, projectId uint) (*models.Task, error) {
	var task models.Task
	err := t.db.First(&task, "id = ? AND project_id = ?", id, projectId).Error
	return utils.HandleError(&task, err)
}

func (t *TaskRepo) CreateTask(task models.Task) (*models.Task, error) {
	err := t.db.Create(&task).Error
	return utils.HandleError(&task, err)
}

func (t *TaskRepo) DeleteTaskById(id uint) error {
	var task models.Task
	if err := t.db.First(&task, id).Error; err != nil {
		return err
	}
	err := t.db.Delete(&task, "id = ?", id).Error
	return err
}

func (t *TaskRepo) UpdateTask(values map[string]any, id uint, projectId uint) error {
	var task models.Task
	if err := t.db.First(&task, "id = ? AND project_id = ?", id, projectId).Error; err != nil {
		return err
	}
	err := t.db.Model(&task).Where("id = ?", id).Where("project_id = ?", projectId).Updates(values).Error
	return err
}

func (t *TaskRepo) SearchTask(query map[string]any, projectId uint) (*[]models.Task, error) {
	var tasks []models.Task
	var task models.Task
	ctx := t.db.Model(&task)
	if title, ok := query["title"].(string); ok {
		ctx.Where("title LIKE ? AND project_id = ?", "%"+title+"%", projectId)
	}
	if priority, ok := query["priority"].(int); ok {
		ctx.Where("priority = ? AND project_id = ?", priority, projectId)
	}
	if creatorId, ok := query["creatorId"].(uint); ok {
		ctx.Where("creator_id = ? AND project_id = ?", creatorId, projectId)
	}
	err := ctx.Find(&tasks).Error
	return utils.HandleError(&tasks, err)
}

func (t *TaskRepo) GetTaskInProgressCountByProjectId(projectId uint) int64 {
	var task models.Task
	var count int64
	t.db.Model(&task).Where("project_id = ?", projectId).Where("status = ?", constant.TASK_STATUS_IN_PROGRESS).Count(&count)
	return count
}

func (t *TaskRepo) GetTaskDoneCountByProjectId(projectId uint) int64 {
	var task models.Task
	var count int64
	t.db.Model(&task).Where("project_id = ?", projectId).Where("status = ?", constant.TASK_STATUS_DONE).Count(&count)
	return count
}

func (t *TaskRepo) GetTaskHighPriorityCountByProjectId(projectId uint) int64 {
	var task models.Task
	var count int64
	t.db.Model(&task).Where("project_id = ?", projectId).Where("priority = ?", constant.TASK_PRIORITY_HIGH).Count(&count)
	return count
}

func (t *TaskRepo) GetTaskMediumPriorityCountByProjectId(projectId uint) int64 {
	var task models.Task
	var count int64
	t.db.Model(&task).Where("project_id = ?", projectId).Where("priority = ?", constant.TASK_PRIORITY_MEDIUM).Count(&count)
	return count
}

func (t *TaskRepo) GetTaskLowPriorityCountByProjectId(projectId uint) int64 {
	var task models.Task
	var count int64
	t.db.Model(&task).Where("project_id = ?", projectId).Where("priority = ?", constant.TASK_PRIORITY_LOW).Count(&count)
	return count
}

func (t *TaskRepo) CheckTaskDoneById(id uint) bool {
	var task models.Task
	count := t.db.Where("id = ? AND status = ?", id, constant.TASK_STATUS_DONE).Find(&task).RowsAffected
	return count > 0
}

func (t *TaskRepo) CheckTaskInProgressById(id uint) bool {
	var task models.Task
	count := t.db.Where("id = ? AND status = ?", id, constant.TASK_STATUS_IN_PROGRESS).Find(&task).RowsAffected
	return count > 0
}

func (t *TaskRepo) GetDueDateByTaskId(id uint) (models.Task, error) {
	var task models.Task
	err := t.db.Where("due_date IS NOT NULL AND id = ?", id).Find(&task).First(&task).Error
	if err != nil {
		return models.Task{}, err
	}
	return task, err
}

func (t *TaskRepo) GetCreatorIdByTaskId(id uint) (*uint, error) {
	var task models.Task
	err := t.db.Find(&task, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	creatorId := task.CreatorID
	return &creatorId, err
}
