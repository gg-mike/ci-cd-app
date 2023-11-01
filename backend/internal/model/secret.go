package model

import (
	"errors"

	"github.com/gg-mike/ci-cd-app/backend/internal/vault"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Migrating table twice causes error (issue: https://github.com/go-gorm/gorm/issues/4946)
type SecretCore struct {
	Key        string        `json:"key"         gorm:"not null"`
	ProjectID  uuid.NullUUID `json:"project_id"  gorm:"type:uuid"`
	PipelineID uuid.NullUUID `json:"pipeline_id" gorm:"type:uuid"`
	Unique     string        `json:"-"           gorm:"->;type: text GENERATED ALWAYS AS (key || '/' || COALESCE(project_id::text,'') || '/' || COALESCE(pipeline_id::text,'')) STORED;uniqueIndex:idx_secrets;default:(-)"`
}

type Secret struct {
	SecretCore
	Common
}

type SecretCreate struct {
	SecretCore
	Value string `json:"value"`
}

type SecretShort struct {
	SecretCore
	Common
}

func (m *Secret) AfterCreate(tx *gorm.DB) error {
	obj, ok := tx.InstanceGet("obj")
	if !ok {
		return errors.New("no obj given in instance")
	}
	return vault.Set(m.ID.String(), map[string]any{"value": obj.(map[string]any)["value"]})
}

func (m *Secret) BeforeUpdate(tx *gorm.DB) error {
	return vault.Del(m.ID.String())
}

func (m *Secret) AfterUpdate(tx *gorm.DB) error {
	obj, ok := tx.InstanceGet("obj")
	if !ok {
		return errors.New("no obj given in instance")
	}
	return vault.Set(m.ID.String(), map[string]any{"value": obj.(map[string]any)["value"]})
}

func (m *Secret) AfterDelete(tx *gorm.DB) error {
	return vault.Del(m.ID.String())
}
