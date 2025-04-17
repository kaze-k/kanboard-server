package models

type TaskAssignee struct {
	ProjectID uint   `gorm:"primary_key" json:"project_id"`
	TaskID    uint   `gorm:"primary_key" json:"task_id"`
	UserID    uint   `gorm:"primary_key" json:"user_id"`
	Username  string `gorm:"size:255;not null" json:"username"`
}
