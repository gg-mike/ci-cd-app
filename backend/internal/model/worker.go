package model

import (
	"errors"

	"github.com/gg-mike/ci-cd-app/backend/internal/vault"
	"gorm.io/gorm"
)

type WorkerType int

const (
	WorkerStatic WorkerType = iota
	WorkerDockerHost
)

type WorkerStatus int

const (
	WorkerIdle WorkerStatus = iota
	WorkerUsed
	WorkerUnreachable
)

type WorkerCore struct {
	Name     string       `json:"name"     gorm:"uniqueIndex:idx_workers;not null"`
	Address  string       `json:"address"  gorm:"not null"`
	System   string       `json:"system"   gorm:"not null"`
	Type     WorkerType   `json:"type"     gorm:"not null"`
	Username string       `json:"username" gorm:"not null"`
}

type Worker struct {
	WorkerCore
	Status     WorkerStatus `json:"status" gorm:"not null;default:0"`
	Builds     []Build      `json:"builds" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;not null"`
	Common
}

type WorkerCreate struct {
	WorkerCore
	PrivateKey string `json:"private_key"`
}

type WorkerShort struct {
	WorkerCore
	Common
}

func (m *Worker) AfterCreate(tx *gorm.DB) error {
	obj, ok := tx.InstanceGet("obj")
	if !ok {
		return errors.New("no obj given in instance")
	}
	return vault.Set(m.ID.String(), map[string]any{ "value": obj.(map[string]any)["private_key"] })
}

func (m *Worker) BeforeUpdate(tx *gorm.DB) error {
	return vault.Del(m.ID.String())
}

func (m *Worker) AfterUpdate(tx *gorm.DB) error {
	obj, ok := tx.InstanceGet("obj")
	if !ok {
		return errors.New("no obj given in instance")
	}
	return vault.Set(m.ID.String(), map[string]any{ "value": obj.(map[string]any)["private_key"] })
}

func (m *Worker) BeforeDelete(tx *gorm.DB) error {
	for _, build := range m.Builds {
		if build.Status == BuildRunning {
			return errors.New("cannot delete worker with running builds")
		}
	}
	tx.Model(&m).Update("deleted", true)
	return nil
}

func (m *Worker) AfterDelete(tx *gorm.DB) error {
	return vault.Del(m.ID.String())
}
