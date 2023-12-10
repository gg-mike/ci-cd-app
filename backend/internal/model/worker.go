package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/gg-mike/ci-cd-app/backend/internal/scheduler"
	"github.com/gg-mike/ci-cd-app/backend/internal/ssh"
	"github.com/gg-mike/ci-cd-app/backend/internal/vault"
	"github.com/google/uuid"
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

type WorkerStrategy int

const (
	WorkerUseMin = iota
	WorkerBalanced
	WorkerUseMax
)

type Worker struct {
	ID           uuid.UUID      `json:"id"            gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string         `json:"name"          gorm:"uniqueIndex:idx_workers;not null"`
	Address      string         `json:"address"       gorm:"not null"`
	System       string         `json:"system"        gorm:"not null"`
	Username     string         `json:"username"      gorm:"not null"`
	Type         WorkerType     `json:"type"          gorm:"not null"`
	Status       WorkerStatus   `json:"status"        gorm:"not null;default:0"`
	Strategy     WorkerStrategy `json:"strategy"      gorm:"not null;default:0"`
	ActiveBuilds int            `json:"active_builds" gorm:"not null;default:0"`
	Capacity     int            `json:"capacity"      gorm:"not null;default:0"`
	Builds       []Build        `json:"builds"        gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;not null"`
	CreatedAt    time.Time      `json:"created_at"    gorm:"default:now()"`
	UpdatedAt    time.Time      `json:"updated_at"    gorm:"default:now()"`
}

type WorkerInput struct {
	Name       string     `json:"name"`
	Address    string     `json:"address"`
	System     string     `json:"system"`
	Type       WorkerType `json:"type"`
	Username   string     `json:"username"`
	PrivateKey string     `json:"private_key"`
	Capacity   int        `json:"capacity"`
}

type WorkerShort struct {
	ID           uuid.UUID      `json:"id"`
	Name         string         `json:"name"`
	Address      string         `json:"address"`
	System       string         `json:"system"`
	Username     string         `json:"username"`
	Type         WorkerType     `json:"type"`
	Status       WorkerStatus   `json:"status"`
	Strategy     WorkerStrategy `json:"strategy"`
	ActiveBuilds int            `json:"active_builds"`
	Capacity     int            `json:"capacity"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (m *Worker) BeforeCreate(tx *gorm.DB) error {
	privateKey, ok := GetFromRaw[string](tx, "private_key")
	if !ok {
		return errors.New("no private_key field given in instance")
	}
	if !testConnection(*m, privateKey) {
		m.Status = WorkerUnreachable
	}
	return nil
}

func (m *Worker) AfterCreate(tx *gorm.DB) error {
	privateKey, _ := GetFromRaw[string](tx, "private_key")
	return vault.SetStr(m.ID.String(), privateKey)
}

func (m *Worker) AfterSave(tx *gorm.DB) error {
	go scheduler.Get().ChangeInWorkers()
	return nil
}

func (m *Worker) BeforeUpdate(tx *gorm.DB) error {
	if _, ok := GetFromRaw[string](tx, "private_key"); !ok {
		return nil
	}
	return vault.Del(m.ID.String())
}

func (m *Worker) AfterUpdate(tx *gorm.DB) error {
	prev, ok := tx.InstanceGet("prev")
	if !ok {
		return errors.New("prev worker not given")
	}
	privateKey, ok := GetFromRaw[string](tx, "private_key")
	if !ok {
		pKey, err := vault.Str(m.ID.String())
		if err != nil {
			return fmt.Errorf("error during retrieving private key: %v", err)
		}
		privateKey = pKey
	}
	var status WorkerStatus
	if !testConnection(*m, privateKey) {
		status = WorkerUnreachable
	} else if prev.(Worker).Status != WorkerUnreachable {
		status = prev.(Worker).Status
	} else {
		status = WorkerIdle
	}
	if err := tx.Model(&m).UpdateColumn("status", status).Error; err != nil {
		return err
	}
	if !ok {
		return nil
	}

	return vault.SetStr(m.ID.String(), privateKey)
}

func (m *Worker) BeforeDelete(tx *gorm.DB) error {
	for _, build := range m.Builds {
		if build.Status == BuildRunning {
			return errors.New("cannot delete worker with running builds")
		}
	}
	return nil
}

func (m *Worker) AfterDelete(tx *gorm.DB) error {
	return vault.Del(m.ID.String())
}

func testConnection(worker Worker, privateKey string) bool {
	return ssh.CheckConnection(worker.Username, worker.Address, privateKey) == nil
}
