package queue

import (
	"time"

	"github.com/gg-mike/ci-cd-app/backend/internal/engine/executor/build"
	"github.com/google/uuid"
)

type QueueElem struct {
	ID        uuid.UUID     `json:"id"         gorm:"primaryKey;type:uuid"`
	Context   build.Context `json:"context"    gorm:"not null;serializer:json"`
	CreatedAt time.Time     `json:"created_at" gorm:"default:now()"`
}
