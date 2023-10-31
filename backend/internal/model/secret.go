package model

type SecretCore struct {
	Key        string `json:"key"         gorm:"uniqueIndex:sel;not null"`
	ProjectID  string `json:"project_id"  gorm:"type:uuid;uniqueIndex:sel;not null"`
	PipelineID string `json:"pipeline_id" gorm:"type:uuid;uniqueIndex:sel;not null"`
}

type Secret struct {
	SecretCore
	Common
}

type SecretCreate struct {
	SecretCore
	Value      string `json:"value"`
}

type SecretShort struct {
	SecretCore
	Common
}
