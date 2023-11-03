package model

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Migrating table twice causes error (issue: https://github.com/go-gorm/gorm/issues/4946)
type Variable struct {
	ID         uuid.UUID  `json:"id"          gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Key        string     `json:"key"         gorm:"not null"`
	Value      string     `json:"value"       gorm:"not null"`
	ProjectID  *uuid.UUID `json:"project_id"  gorm:"type:uuid"`
	PipelineID *uuid.UUID `json:"pipeline_id" gorm:"type:uuid"`
	Unique     string     `json:"-"           gorm:"->;type: text GENERATED ALWAYS AS (key || '/' || COALESCE(project_id::text,'') || '/' || COALESCE(pipeline_id::text,'')) STORED;uniqueIndex:idx_variables;default:(-)"`
	CreatedAt  time.Time  `json:"created_at"  gorm:"default:now()"`
	UpdatedAt  time.Time  `json:"updated_at"  gorm:"default:now()"`
}

type VariableInput struct {
	Key        string     `json:"key"`
	Value      string     `json:"value"`
	ProjectID  *uuid.UUID `json:"project_id"`
	PipelineID *uuid.UUID `json:"pipeline_id"`
}

func (m *Variable) BeforeCreate(tx *gorm.DB) error {
	if strings.HasPrefix(m.Key, "_") {
		return errors.New("variable cannot start with '_'")
	}
	return nil
}

func (m *Variable) AfterUpdate(tx *gorm.DB) error {
	if strings.HasPrefix(m.Key, "_") {
		return errors.New("variable cannot start with '_'")
	}
	return nil
}
