package model

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type BuildStatus uint8

const (
	BuildScheduled BuildStatus = iota
	BuildRunning
	BuildSuccessful
	BuildFailed
	BuildCanceled
)

type BuildCore struct {
	WorkerID   uuid.UUID `json:"worker_id"   gorm:"type:uuid;not null"`
	PipelineID uuid.UUID `json:"pipeline_id" gorm:"type:uuid;index:,unique,composite:idx_builds"`
}

type Build struct {
	BuildCore
	Number  uint           `json:"number"   gorm:"index:,unique,composite:idx_builds;autoIncrement"`
	RevList pq.StringArray `json:"rev_list" gorm:"type:text[];not null"`
	Status  BuildStatus    `json:"status"   gorm:"not null;default:0"`
	Steps   []BuildStep    `json:"steps"    gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Common
}

type BuildCreate struct {
	BuildCore
}

type BuildShort struct {
	BuildCore
	Number    uint             `json:"number"`
	RevList   pq.StringArray   `json:"rev_list"`
	Steps     []BuildStepShort `json:"steps"`
	Common
}

func (m *Build) AfterCreate(tx *gorm.DB) error {
	// TODO: Schedule build
	return nil
}

func (m *Build) BeforeUpdate(tx *gorm.DB) error {
	prev, ok := tx.InstanceGet("prev")
	if !ok {
		return errors.New("prev obj not given")
	}
	if prev.(Build).Status == BuildCanceled {
		return errors.New("cannot change status of canceled build")
	}
	if m.Status != BuildCanceled {
		return errors.New("operation not allowed")
	}
	// TODO: Cancel job
	return nil
}
