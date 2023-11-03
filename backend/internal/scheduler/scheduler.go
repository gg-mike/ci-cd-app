package scheduler

import (
	"github.com/google/uuid"
)

type Scheduler interface {
	Schedule(buildID uuid.UUID)
	Cancel(buildID uuid.UUID) error
}

var scheduler Scheduler

func Init(s Scheduler) {
	scheduler = s
}

func Get() Scheduler {
	return scheduler
}
