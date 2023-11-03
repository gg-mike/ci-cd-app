package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: replace cascade constraint with trigger (issue: https://github.com/go-gorm/gorm/issues/5001)
type Project struct {
	ID        uuid.UUID  `json:"id"         gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string     `json:"name"       gorm:"uniqueIndex:idx_projects"`
	Repo      string     `json:"repo"       gorm:"not null"`
	Variables []Variable `json:"variables"  gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Secrets   []Secret   `json:"secrets"    gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Pipelines []Pipeline `json:"pipelines"  gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:now()"`
}

type ProjectInput struct {
	Name string `json:"name"`
	Repo string `json:"repo"`
}

type ProjectShort struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Repo      string    `json:"repo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Project) BeforeDelete(tx *gorm.DB) error {
	if !IsForce(tx) {
		if len(m.Pipelines) == 0 {
			tx.Model(&m).UpdateColumn("deleted", true)
			return nil
		}
		return errors.New("cannot delete project with pipelines (use 'force' query param to overwrite)")
	}

	var count int64
	tx.Joins("builds", tx.Where(&Build{Status: BuildRunning})).Count(&count).
		Find(nil, &Pipeline{ProjectID: m.ID})

	if count == 0 {
		tx.Model(&m).UpdateColumn("deleted", true)
		return nil
	}
	return errors.New("cannot delete project with running builds")
}
