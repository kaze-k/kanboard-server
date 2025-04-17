package models

import (
	"fmt"

	"server/internal/event"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name string  `gorm:"size:255;not null"`
	Desc *string `gorm:"type:text;default:null"`
}

func (p *Project) AfterCreate(db *gorm.DB) error {
	content := fmt.Sprintf("『%s』新项目创建成功", p.Name)
	event.AdminPublish(event.Event{Content: &content})
	return nil
}

func (p *Project) AfterDelete(db *gorm.DB) error {
	content := fmt.Sprintf("『%s』项目删除成功", p.Name)
	event.AdminPublish(event.Event{Content: &content})
	return nil
}
