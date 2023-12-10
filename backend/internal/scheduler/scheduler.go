package scheduler

import (
	"github.com/google/uuid"
)

type Scheduler interface {
	Schedule(buildID uuid.UUID)
	Finished(any)
	ChangeInWorkers()
}

var scheduler Scheduler

func Init(s Scheduler) {
	scheduler = s
}

func Get() Scheduler {
	return scheduler
}
