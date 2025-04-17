package models

import (
	"fmt"
	"time"

	"server/internal/constant"
	"server/internal/event"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title     string    `gorm:"size:255;not null"`
	Desc      string    `gorm:"not null"`
	Status    uint      `gorm:"size:1;index;default:0;not null"`
	DueDate   time.Time `gorm:"size:1;index;default:null"`
	Priority  int       `gorm:"size:1;index;default:0;not null"`
	ProjectID uint      `gorm:"index;not null"`
	CreatorID uint      `gorm:"index;not null"`
}

func (t *Task) AfterCreate(db *gorm.DB) error {
	eventType := constant.TASK_EVENT
	content := fmt.Sprintf("新增任务『%s』", t.Title)
	event.KanboardPublish(event.Event{EventType: &eventType, Content: &content, ProjectID: &t.ProjectID, TaskID: &t.ID})
	return nil
}

func (t *Task) AfterUpdate(db *gorm.DB) error {
	var content string
	if t.Status == constant.TASK_STATUS_UNDO {
		content = fmt.Sprintf("任务『%s』标记为未完成", t.Title)
	} else if t.Status == constant.TASK_STATUS_DONE {
		content = fmt.Sprintf("任务『%s』标记为完成", t.Title)
	} else if t.Status == constant.TASK_STATUS_IN_PROGRESS {
		content = fmt.Sprintf("任务『%s』标记为进行中", t.Title)
	} else {
		return nil
	}

	eventType := constant.TASK_EVENT
	event.KanboardPublish(event.Event{EventType: &eventType, Content: &content, ProjectID: &t.ProjectID, TaskID: &t.ID})
	return nil
}
