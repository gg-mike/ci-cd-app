package model

import (
	"errors"

	"gorm.io/gorm"
)

// TODO: rbac
type UserCore struct {
	Username string `json:"username" gorm:"uniqueIndex:idx_users"`
}

type User struct {
	UserCore
	Common
}

type UserCreate struct {
	UserCore
}

type UserShort struct {
	UserCore
	Common
}

func (m *User) BeforeUpdate(tx *gorm.DB) error {
	prev, ok := tx.InstanceGet("prev")
	if !ok {
		return errors.New("prev obj not given")
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
