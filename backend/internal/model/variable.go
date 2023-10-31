package model

import (
	"github.com/google/uuid"
)

// Migrating table twice causes error (issue: https://github.com/go-gorm/gorm/issues/4946)
type VariableCore struct {
	Key        string        `json:"key"         gorm:"not null"`
	Value      []byte        `json:"value"       gorm:"not null"`
	ProjectID  uuid.NullUUID `json:"project_id"  gorm:"type:uuid"`
	PipelineID uuid.NullUUID `json:"pipeline_id" gorm:"type:uuid"`
	Unique     string        `json:"-"           gorm:"->;type: text GENERATED ALWAYS AS (key || '/' || COALESCE(project_id::text,'') || '/' || COALESCE(pipeline_id::text,'')) STORED;uniqueIndex:idx_variables;default:(-)"`
}

type Variable struct {
	VariableCore
	Common
}

type VariableCreate struct {
	VariableCore
}

type VariableShort struct {
	VariableCore
	Common
}
