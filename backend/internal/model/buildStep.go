package model

import (
	"time"

	"github.com/google/uuid"
)

type BuildStepCore struct {
	Name     string        `json:"name"     gorm:"uniqueIndex:sel;not null"`
	Duration time.Duration `json:"duration" gorm:"not null"`
	BuildID  uuid.UUID     `json:"build_id" gorm:"type:uuid;uniqueIndex:sel;not null"`
}

type BuildStep struct {
	BuildStepCore
	Logs          []BuildLog `json:"logs" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Common
}

type BuildStepCreate struct {
	BuildStepCore
}

type BuildStepShort struct {
	BuildStepCore
	Common
}

type BuildLog struct {
	Command     string    `json:"command"       gorm:"primaryKey"`
	Output      string    `json:"output"        gorm:"not null"`
	BuildStepID uuid.UUID `json:"build_step_id" gorm:"type:uuid;primaryKey"`
}
