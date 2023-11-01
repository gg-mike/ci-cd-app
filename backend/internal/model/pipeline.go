package model

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PipelineCore struct {
	Name      string    `json:"name"       gorm:"uniqueIndex:idx_pipelines;not null"`
	Branch    string    `json:"branch"     gorm:"not null"`
	ProjectID uuid.UUID `json:"project_id" gorm:"type:uuid;uniqueIndex:idx_pipelines;not null"`
}

// TODO: replace cascade constraint with trigger (issue: https://github.com/go-gorm/gorm/issues/5001)
type Pipeline struct {
	PipelineCore
	LastRev   string         `json:"last_rev"  gorm:"not null"`
	Config    PipelineConfig `json:"config"    gorm:"not null;serializer:json"`
	Variables []Variable     `json:"variables" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Secrets   []Secret       `json:"secrets"   gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Builds    []Build        `json:"builds"    gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Common
}

type PipelineCreate struct {
	PipelineCore
	Config PipelineConfig `json:"config"`
}

type PipelineShort struct {
	PipelineCore
	Common
}

type PipelineConfig struct {
	System string `json:"system"`
	Image  string `json:"image"`
	Steps  []struct {
		Name     string   `json:"name"`
		Commands []string `json:"commands"`
	} `json:"steps"`
}

func (m *Pipeline) BeforeDelete(tx *gorm.DB) error {
	force, ok := tx.InstanceGet("force")
	if !force.(bool) || !ok {
		if len(m.Builds) == 0 {
			tx.Model(&m).Update("deleted", true)
			return nil
		}
		return errors.New("cannot delete project with builds (use 'force' query param to overwrite)")
	}

	for _, build := range m.Builds {
		if build.Status == BuildRunning {
			return errors.New("cannot delete pipeline with running builds")
		}
	}
	tx.Model(&m).Update("deleted", true)
	return nil
}
