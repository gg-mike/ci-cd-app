package model

type WorkerType int

const (
	WorkerStatic WorkerType = iota
	WorkerDockerHost
)

type WorkerStatus int

const (
	WorkerIdle WorkerStatus = iota
	WorkerUsed
	WorkerUnreachable
)

type WorkerCore struct {
	Name     string       `json:"name"     gorm:"uniqueIndex;not null"`
	Address  string       `json:"address"  gorm:"not null"`
	System   string       `json:"system"   gorm:"not null"`
	Status   WorkerStatus `json:"status"   gorm:"not null"`
	Type     WorkerType   `json:"type"     gorm:"not null"`
	Username string       `json:"username" gorm:"not null"`
}

type Worker struct {
	WorkerCore
	Builds     []Build `json:"builds" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	Common
}

type WorkerCreate struct {
	WorkerCore
	PrivateKey string `json:"private_key"`
}

type WorkerShort struct {
	WorkerCore
	Common
}
