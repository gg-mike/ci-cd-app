package model

import (
	"errors"
	"strings"
	"time"

	"github.com/gg-mike/ci-cd-app/backend/internal/vault"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Migrating table twice causes error (issue: https://github.com/go-gorm/gorm/issues/4946)
type Secret struct {
	ID         uuid.UUID  `json:"id"          gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Key        string     `json:"key"         gorm:"not null"`
	ProjectID  *uuid.UUID `json:"project_id"  gorm:"type:uuid"`
	PipelineID *uuid.UUID `json:"pipeline_id" gorm:"type:uuid"`
	Unique     string     `json:"-"           gorm:"->;type: text GENERATED ALWAYS AS (key || '/' || COALESCE(project_id::text,'') || '/' || COALESCE(pipeline_id::text,'')) STORED;uniqueIndex:idx_secrets;default:(-)"`
	CreatedAt  time.Time  `json:"created_at"  gorm:"default:now()"`
	UpdatedAt  time.Time  `json:"updated_at"  gorm:"default:now()"`
}

type SecretInput struct {
	Key        string     `json:"key"`
	Value      string     `json:"value"`
	ProjectID  *uuid.UUID `json:"project_id"`
	PipelineID *uuid.UUID `json:"pipeline_id"`
}

func (m *Secret) BeforeCreate(tx *gorm.DB) error {
	if strings.HasPrefix(m.Key, "_") {
		return errors.New("secret cannot start with '_'")
	}
	return nil
}

func (m *Secret) AfterCreate(tx *gorm.DB) error {
	value, ok := GetFromRaw[string](tx, "value")
	if !ok {
		return errors.New("no value field given in instance")
	}
	return vault.SetStr(m.ID.String(), value)
}

func (m *Secret) AfterUpdate(tx *gorm.DB) error {
	if strings.HasPrefix(m.Key, "_") {
		return errors.New("secret cannot start with '_'")
	}
	if value, ok := GetFromRaw[string](tx, "value"); ok {
		return vault.SetStr(m.ID.String(), value)
	}
	return nil
}

func (m *Secret) AfterDelete(tx *gorm.DB) error {
	return vault.Del(m.ID.String())
}
