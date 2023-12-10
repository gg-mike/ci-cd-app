package scheduler

import (
	"time"

	"github.com/gg-mike/ci-cd-app/backend/cmd/server/queue"
	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/engine/executor/build"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
)

func bindWorker(queueElem *queue.QueueElem, worker model.Worker) error {
	id := queueElem.Context.Build.ID

	logger.Debug(module).Str("build_id", id.String()).Str("step", "worker-bind").Msg("worker binding started")
	queueElem.Context.Build.Steps = append(queueElem.Context.Build.Steps, model.BuildStep{Name: "Worker binding", BuildID: queueElem.Context.Build.ID, Logs: []model.BuildLog{}, Number: 1})
	queueElem.Context.Worker = worker

	if err := db.Get().Model(&queueElem.Context.Worker).UpdateColumn("active_builds", queueElem.Context.Worker.ActiveBuilds+1).Error; err != nil {
		logger.Error(module).Str("build_id", id.String()).Str("step", "worker-bind").Err(err).Msg("could not update worker")
		return ErrBuildSave
	}

	build.AppendLog(&queueElem.Context, 1, "BIND", "worker ["+queueElem.Context.Worker.Name+"] bound")
	logger.Debug("scheduler").Msgf("worker: %+v", queueElem.Context.Worker)
	queueElem.Context.Build.Steps[1].Duration = time.Since(queueElem.CreatedAt)
	logger.Debug(module).Str("build_id", id.String()).Str("step", "worker-bind").Msg("worker binding succeeded")

	if err := db.Get().Create(&queueElem.Context.Build.Steps[1]).Error; err != nil {
		logger.Error(module).Str("build_id", id.String()).Str("step", "worker-bind").Err(err).Msg("could not write build steps")
		return ErrBuildSave
	}

	if err := db.Get().Model(&queueElem.Context.Build).UpdateColumns(&map[string]any{"status": model.BuildRunning, "worker_id": queueElem.Context.Worker.ID}).Error; err != nil {
		logger.Error(module).Str("build_id", id.String()).Str("step", "worker-bind").Err(err).Msg("could not update build")
		return ErrBuildSave
	}

	return nil
}
