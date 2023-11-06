package model

import (
	"time"

	"github.com/google/uuid"
)

type BuildStep struct {
	ID        uuid.UUID     `json:"id"         gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Number    int           `json:"number"     gorm:"not null"`
	Name      string        `json:"name"       gorm:"uniqueIndex:idx_build_steps;not null"`
	Duration  time.Duration `json:"duration"   gorm:"not null"`
	Logs      []BuildLog    `json:"logs"       gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	BuildID   uuid.UUID     `json:"build_id"   gorm:"type:uuid;uniqueIndex:idx_build_steps;not null"`
	CreatedAt time.Time     `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time     `json:"updated_at" gorm:"default:now()"`
}

type BuildStepShort struct {
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
}

type BuildLog struct {
	Command     string    `json:"command"         gorm:"primaryKey"` // TODO: error: command can be bigger than limit for pk
	Idx         int       `json:"idx,omitempty"`
	Total       int       `json:"total,omitempty"`
	Output      string    `json:"output"          gorm:"not null"`
	BuildStepID uuid.UUID `json:"build_step_id"   gorm:"type:uuid;primaryKey"`
}
