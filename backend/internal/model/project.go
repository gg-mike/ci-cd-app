package model

type ProjectCore struct {
	Name string `json:"name" gorm:"uniqueIndex"`
	Repo string `json:"repo" gorm:"not null"`
}

type Project struct {
	ProjectCore
	Variables   []Variable `json:"variables" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Secrets     []Secret   `json:"secrets"   gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Pipelines   []Pipeline `json:"pipelines" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Common
}

type ProjectCreate struct {
	ProjectCore
}

type ProjectShort struct {
	ProjectCore
	Common
}
