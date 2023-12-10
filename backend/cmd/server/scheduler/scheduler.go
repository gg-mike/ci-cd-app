package scheduler

import (
	"errors"

	"github.com/gg-mike/ci-cd-app/backend/cmd/server/queue"
	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/engine/executor/build"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/google/uuid"
)

const module = "scheduler"

var (
	ErrBuildSave       = errors.New("unable to save build to database")
	ErrBuildInitFailed = errors.New("build init ended with error")
)

type Context struct{}

var newBuild chan uuid.UUID
var finishedBuild chan build.Context

var addToQueue chan build.Context

var changeInWorkers chan any

func Init() {
	newBuild = make(chan uuid.UUID)
	finishedBuild = make(chan build.Context)

	addToQueue = make(chan build.Context)

	changeInWorkers = make(chan any)
}

func Run() {
	logger.Debug(module).Msg("scheduler is running")

	logger.Debug(module).Msg("binding any builds scheduled in previous run")
	bind()

	for {
		select {
		case buildID := <-newBuild:
			logger.Debug(module).
				Str("event", "schedule").
				Str("status", "processed").
				Msgf("scheduling build [id=%s]", buildID.String())

			schedule(buildID)

			logger.Debug(module).
				Str("event", "schedule").
				Str("status", "finished").
				Msgf("scheduling build [id=%s]", buildID.String())
		case ctx := <-finishedBuild:
			logger.Debug(module).
				Str("event", "finished").
				Str("status", "processed").
				Msgf("finishing build [id=%s]", ctx.Build.ID.String())

			finished(ctx)

			logger.Debug(module).
				Str("event", "finished").
				Str("status", "finished").
				Msgf("finishing build [id=%s]", ctx.Build.ID.String())
		case ctx := <-addToQueue:
			logger.Debug(module).
				Str("event", "add-to-queue").
				Str("status", "processed").
				Msgf("adding to queue build [id=%s]", ctx.Build.ID.String())

			if err := db.Get().Create(&queue.QueueElem{ID: ctx.Build.ID, Context: ctx}).Error; err != nil {
				logger.Debug(module).
					Str("event", "add-to-queue").
					Str("status", "failed").
					Err(err).
					Msg("error adding new build to queue")
				continue
			}
			bind()

			logger.Debug(module).
				Str("event", "add-to-queue").
				Str("status", "finished").
				Msgf("adding to queue build [id=%s]", ctx.Build.ID.String())
		case <-changeInWorkers:
			logger.Debug(module).
				Str("event", "change-in-workers").
				Str("status", "processed").
				Msg("change in worker")

			bind()

			logger.Debug(module).
				Str("event", "change-in-workers").
				Str("status", "finished").
				Msg("change in worker")
		}
	}
}

func (s *Context) Schedule(buildID uuid.UUID) {
	logger.Debug(module).
		Str("event", "schedule").
		Str("status", "received").
		Msgf("scheduling build [id=%s]", buildID.String())

	newBuild <- buildID
}

func (s *Context) Finished(ctx any) {
	logger.Debug(module).
		Str("event", "finished").
		Str("status", "received").
		Msgf("finishing build [id=%s]", ctx.(build.Context).Build.ID.String())

	finishedBuild <- ctx.(build.Context)
}

func (s *Context) AddToQueue(ctx build.Context) {
	logger.Debug(module).
		Str("event", "add-to-queue").
		Str("status", "received").
		Msgf("adding to queue build [id=%s]", ctx.Build.ID.String())

	addToQueue <- ctx
}

func (s *Context) ChangeInWorkers() {
	logger.Debug(module).
		Str("event", "change-in-workers").
		Str("status", "received").
		Msg("change in worker")

	changeInWorkers <- true
}
