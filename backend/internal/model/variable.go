package model

type VariableCore struct {
	Key        string `json:"key"         gorm:"uniqueIndex:sel;not null"`
	Value      []byte `json:"value"       gorm:"not null"`
	ProjectID  string `json:"project_id"  gorm:"type:uuid;uniqueIndex:sel;not null"`
	PipelineID string `json:"pipeline_id" gorm:"type:uuid;uniqueIndex:sel;not null"`
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