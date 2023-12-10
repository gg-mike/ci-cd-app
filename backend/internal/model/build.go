package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/gg-mike/ci-cd-app/backend/internal/scheduler"
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

func (s BuildStatus) String() string {
	switch s {
	case BuildScheduled:
		return "scheduled"
	case BuildRunning:
		return "running"
	case BuildSuccessful:
		return "successful"
	case BuildFailed:
		return "failed"
	case BuildCanceled:
		return "canceled"
	default:
		return "unknown"
	}
}

type Build struct {
	ID         uuid.UUID      `json:"id"          gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Number     uint           `json:"number"      gorm:"index:,unique,composite:idx_builds;autoIncrement"`
	RevList    pq.StringArray `json:"rev_list"    gorm:"type:text[];not null;default:'{}'"`
	Status     BuildStatus    `json:"status"      gorm:"not null;default:0"`
	Steps      []BuildStep    `json:"steps"       gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	PipelineID uuid.UUID      `json:"pipeline_id" gorm:"type:uuid;index:,unique,composite:idx_builds"`
	WorkerID   *uuid.UUID     `json:"worker_id"   gorm:"type:uuid"`
	CreatedAt  time.Time      `json:"created_at"  gorm:"default:now()"`
	UpdatedAt  time.Time      `json:"updated_at"  gorm:"default:now()"`
}

type BuildInput struct {
	PipelineID string `json:"pipeline_id"`
}

type BuildShort struct {
	ID         uuid.UUID        `json:"id"`
	Number     uint             `json:"number"`
	RevList    pq.StringArray   `json:"rev_list"    gorm:"type:text[]"`
	Status     BuildStatus      `json:"status"`
	Steps      []BuildStepShort `json:"steps"`
	PipelineID uuid.UUID        `json:"pipeline_id"`
	WorkerID   uuid.NullUUID    `json:"worker_id,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

func (m *Build) AfterCreate(tx *gorm.DB) error {
	go scheduler.Get().Schedule(m.ID)
	return nil
}

func (m *Build) BeforeUpdate(tx *gorm.DB) error {
	prev, ok := tx.InstanceGet("prev")
	if !ok {
		return errors.New("prev build not given")
	}
	switch prev.(Build).Status {
	case BuildScheduled, BuildRunning:
		tx.Statement.SetColumn("status", BuildCanceled)
		go scheduler.Get().Finished(m.ID)
		return nil
	default:
		return fmt.Errorf("cannot change status of build from [%s] to [%s]",
			prev.(Build).Status.String(), BuildCanceled.String())
	}
}
