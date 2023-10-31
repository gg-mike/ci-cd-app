package model

import (
	"github.com/google/uuid"
)

type PipelineCore struct {
	Name      string         `json:"name"       gorm:"uniqueIndex:sel;not null"`
	Branch    string         `json:"branch"     gorm:"not null"`
	Config    PipelineConfig `json:"config"     gorm:"not null;serializer:json"`
	ProjectID uuid.UUID      `json:"project_id" gorm:"type:uuid;uniqueIndex:sel;not null"`
}

type Pipeline struct {
	PipelineCore
	Variables    []Variable `json:"variables"  gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Secrets      []Secret   `json:"secrets"    gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Builds       []Build    `json:"builds"     gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Common
}

type PipelineCreate struct {
	PipelineCore
}

type PipelineShort struct {
	PipelineCore
	Common
}

type PipelineConfig struct{
	System   string   `json:"system"`
	Image    string   `json:"image"`
	Steps    []struct {
		Name     string   `json:"name"`
		Commands []string `json:"commands"`
	}                 `json:"steps"`
}
