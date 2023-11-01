package model

import (
	"errors"

	"gorm.io/gorm"
)

type ProjectCore struct {
	Name string `json:"name" gorm:"uniqueIndex:idx_projects"`
	Repo string `json:"repo" gorm:"not null"`
}

// TODO: replace cascade constraint with trigger (issue: https://github.com/go-gorm/gorm/issues/5001)
type Project struct {
	ProjectCore
	Variables []Variable `json:"variables" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Secrets   []Secret   `json:"secrets"   gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Pipelines []Pipeline `json:"pipelines" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Common
}

type ProjectCreate struct {
	ProjectCore
}

type ProjectShort struct {
	ProjectCore
	Common
}

func (m *Project) BeforeDelete(tx *gorm.DB) error {
	force, ok := tx.InstanceGet("force")
	if !force.(bool) || !ok {
		if len(m.Pipelines) == 0 {
			tx.Model(&m).Update("deleted", true)
			return nil
		}
		return errors.New("cannot delete project with pipelines (use 'force' query param to overwrite)")
	}

	var count int64
	tx.Joins("builds", tx.Where(&Build{Status: BuildRunning})).Count(&count).
		Find(nil, &Pipeline{PipelineCore: PipelineCore{ProjectID: m.ID}})

	if count == 0 {
		tx.Model(&m).Update("deleted", true)
		return nil
	}
	return errors.New("cannot delete project with running builds")
}
