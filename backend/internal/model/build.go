package model

import (
	"github.com/google/uuid"
)

type BuildStatus uint8;

const (
	BuildRunning BuildStatus = iota
	BuildSuccessful
	BuildFailed
	BuildCanceled
)

type BuildCore struct {
	Number     uint        `json:"number"      gorm:"uniqueIndex:sel;autoIncrement"`
	Tags       string      `json:"tags"        gorm:"not null"`
	Status     BuildStatus `json:"status"      gorm:"not null"`
	WorkerID   uuid.UUID   `json:"worker_id"   gorm:"type:uuid;not null"`
	PipelineID uuid.UUID   `json:"pipeline_id" gorm:"type:uuid;uniqueIndex:sel"`
}

type Build struct {
	BuildCore
	Steps     []BuildStep `json:"steps" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Common
}

type BuildCreate struct {
	BuildCore
}

type BuildShort struct {
	BuildCore
	Steps     []BuildStepShort `json:"steps"`
	Common
}
