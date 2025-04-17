package models

import (
	"fmt"

	"server/internal/constant"
	"server/internal/event"
	"server/pkg/crypto"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string          `gorm:"<-create;size:255;not null"`
	Password   string          `gorm:"size:255;not null"`
	Avatar     uint            `gorm:"index;default:0"`
	Gender     constant.Gender `gorm:"size:1;index;default:0"`
	Email      string          `gorm:"size:255;default:null"`
	Mobile     string          `gorm:"size:255;default:null"`
	CreateFrom constant.From   `gorm:"<-create;size:1;index;not null"`
	IsAdmin    bool            `gorm:"default:0;not null"`
	Loginable  bool            `gorm:"not null"`
	Position   string          `gorm:"default:null"`
}

func (u *User) EncryptPassword() error {
	hash, err := crypto.Encrypt(u.Password)
	if err == nil {
		u.Password = hash
		return nil
	}

	return err
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	return u.EncryptPassword()
}

func (u *User) BeforeUpdate(db *gorm.DB) error {
	if db.Statement.Changed("password") {
		var newPassword string
		if pwd, ok := db.Statement.Dest.(map[string]any); ok {
			newPassword = pwd["password"].(string)
		}
		hash, err := crypto.Encrypt(newPassword)
		if err != nil {
			return err
		}
		db.Statement.SetColumn("password", hash)
	}
	return nil
}

func (u *User) AfterUpdate(db *gorm.DB) error {
	event.KanboardPublish(event.Event{UserID: &u.ID})
	return nil
}

func (u *User) AfterCreate(db *gorm.DB) error {
	content := fmt.Sprintf("『%s』新用户创建成功", u.Username)
	event.AdminPublish(event.Event{Content: &content})
	return nil
}

func (u *User) AfterDelete(db *gorm.DB) error {
	content := fmt.Sprintf("『%s』用户删除成功", u.Username)
	event.AdminPublish(event.Event{Content: &content})
	return nil
}
