package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: rbac
type User struct {
	ID        uuid.UUID `json:"id"         gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username  string    `json:"username"   gorm:"uniqueIndex:idx_users"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

type UserInput struct {
	Username string `json:"username"`
}

func (m *User) BeforeUpdate(tx *gorm.DB) error {
	prev, ok := tx.InstanceGet("prev")
	if !ok {
		return errors.New("prev user not given")
	}
	if prev.(User).Username == "admin" {
		return errors.New("admin cannot be changed")
	}
	return nil
}

func (m *User) BeforeDelete(tx *gorm.DB) error {
	if m.Username == "admin" {
		return errors.New("admin cannot be deleted")
	}
	return nil
}
