package model

import (
	"time"

	"github.com/google/uuid"
)

type Common struct {
	ID        uuid.UUID `json:"id"         gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
	Deleted   bool      `json:"-"          gorm:"default:false"`
}
