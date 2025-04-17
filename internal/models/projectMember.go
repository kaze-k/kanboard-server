package models

import (
	"fmt"
	"time"

	"server/internal/constant"
	"server/internal/event"

	"gorm.io/gorm"
)

type ProjectMember struct {
	ProjectID uint      `gorm:"primary_key;index" json:"project_id"`
	UserID    uint      `gorm:"primary_key;index" json:"user_id"`
	Username  string    `gorm:"size:255;not null" json:"username"`
	Assignee  bool      `gorm:"not null;default:false" json:"assignee"`
	JoinedAt  time.Time `json:"joined_at"`
}

func (p *ProjectMember) BeforeCreate(db *gorm.DB) error {
	p.JoinedAt = time.Now()
	return nil
}

func (p *ProjectMember) AfterUpdate(db *gorm.DB) error {
	projectID := p.ProjectID
	var project Project
	err := db.Find(&project, "id = ?", projectID).Error
	if err != nil {
		return err
	}

	var content string
	if p.Assignee {
		content = fmt.Sprintf("『%s』成为项目『%s』的负责人", p.Username, project.Name)
	} else {
		content = fmt.Sprintf("『%s』不再是项目『%s』的负责人", p.Username, project.Name)
	}
	eventType := constant.PROJECT_EVENT
	event.KanboardPublish(event.Event{EventType: &eventType, Content: &content, ProjectID: &projectID, UserID: &p.UserID})
	return nil
}

func (p *ProjectMember) AfterCreate(db *gorm.DB) error {
	projectID := p.ProjectID
	var project Project
	err := db.Find(&project, "id = ?", projectID).Error
	if err != nil {
		return err
	}
	content := fmt.Sprintf("『%s』用户加入『%s』项目", p.Username, project.Name)
	eventType := constant.PROJECT_EVENT
	event.KanboardPublish(event.Event{EventType: &eventType, Content: &content, ProjectID: &projectID})
	event.KanboardPublish(event.Event{Content: &content, UserID: &p.UserID, ProjectID: &projectID})
	return nil
}

func (p *ProjectMember) AfterDelete(db *gorm.DB) error {
	projectID := p.ProjectID
	var project Project
	err := db.Find(&project, "id = ?", projectID).Error
	if err != nil {
		return err
	}
	if project.Name == "" {
		project.Name = "项目已被删除"
	}
	content := fmt.Sprintf("『%s』用户退出『%s』项目", p.Username, project.Name)
	eventType := constant.PROJECT_EVENT
	event.KanboardPublish(event.Event{EventType: &eventType, Content: &content, ProjectID: &projectID})
	event.KanboardPublish(event.Event{UserID: &p.UserID, ProjectID: &projectID})
	return nil
}
