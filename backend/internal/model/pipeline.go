package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: replace cascade constraint with trigger (issue: https://github.com/go-gorm/gorm/issues/5001)
type Pipeline struct {
	ID        uuid.UUID      `json:"id"         gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string         `json:"name"       gorm:"uniqueIndex:idx_pipelines;not null"`
	Branch    string         `json:"branch"     gorm:"not null"`
	LastRev   string         `json:"last_rev"   gorm:"not null"`
	Config    PipelineConfig `json:"config"     gorm:"not null;serializer:json"`
	Variables []Variable     `json:"variables"  gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Secrets   []Secret       `json:"secrets"    gorm:"not null"`
	Builds    []Build        `json:"builds"     gorm:"not null"`
	ProjectID uuid.UUID      `json:"project_id" gorm:"type:uuid;uniqueIndex:idx_pipelines;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:now()"`
}

type PipelineInput struct {
	Name      string         `json:"name"`
	Branch    string         `json:"branch"`
	Config    PipelineConfig `json:"config"`
	ProjectID uuid.UUID      `json:"project_id"`
}

type PipelineShort struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Branch    string    `json:"branch"`
	ProjectID uuid.UUID `json:"project_id"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

type PipelineConfig struct {
	System  string               `json:"system"`
	Image   string               `json:"image"`
	Steps   []PipelineConfigStep `json:"steps"`
	Cleanup []string             `json:"cleanup"`
}

type PipelineConfigStep struct {
	Name     string   `json:"name"`
	Commands []string `json:"commands"`
}

func (m *Pipeline) BeforeCreate(tx *gorm.DB) error {
	// TODO: validate config
	return nil
}

func (m *Pipeline) AfterUpdate(tx *gorm.DB) error {
	// TODO: validate config
	return nil
}

func (m *Pipeline) BeforeDelete(tx *gorm.DB) error {
	if !IsForce(tx) {
		if len(m.Builds) == 0 {
			return nil
		}
		return errors.New("cannot delete project with builds (use 'force' query param to overwrite)")
	}

	for _, build := range m.Builds {
		if build.Status == BuildRunning {
			return errors.New("cannot delete pipeline with running builds")
		}
	}
	return nil
}
