package schedulerImpl

import (
	"errors"
	"time"

	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/engine/executor"
	"github.com/gg-mike/ci-cd-app/backend/internal/engine/executor/build"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const module = "scheduler"

type Context struct {
}

var (
	ErrBuildSave       = errors.New("unable to save build to database")
	ErrBuildInitFailed = errors.New("build init ended with error")
)

func (ctx Context) Schedule(buildID uuid.UUID) {
	buildCtx, err := buildContext(buildID)
	if err != nil {
		logger.Fatal(module).Str("build_id", buildID.String()).Str("step", "context-create").Err(err).Msg("fatal error during build context creation")
	}
	id := buildCtx.Build.ID

	if err := bindWorker(&buildCtx); err != nil {
		logger.Error(module).Str("build_id", id.String()).Str("step", "worker-bind").Err(err).Msg("error during worker binding")
		if err := db.Get().Model(&buildCtx.Build).UpdateColumn("status", model.BuildFailed).Error; err != nil {
			logger.Error(module).Str("build_id", id.String()).Str("step", "worker-bind").Err(err).Msg("could not update build")
		}
		return
	}

	logger.Debug(module).Str("build_id", id.String()).Str("step", "execute").Msg("build execution started")
	if err := executor.Execute(&buildCtx); err != nil {
		logger.Warn(module).Str("build_id", id.String()).Str("step", "execute").Err(err).Msg("build execution ended with error")
		if err := db.Get().Model(&buildCtx.Build).UpdateColumn("status", model.BuildFailed).Error; err != nil {
			logger.Error(module).Str("build_id", id.String()).Str("step", "execute").Err(err).Msg("could not update build")
		}
		return
	}

	if buildCtx.Build.Status == model.BuildCanceled {
		logger.Debug(module).Str("build_id", id.String()).Str("step", "execute").Str("status", model.BuildCanceled.String()).Msg("build execution canceled")
		return
	}

	logger.Debug(module).Str("build_id", id.String()).Str("step", "execute").Str("status", model.BuildSuccessful.String()).Msg("build execution ended")
	if err := db.Get().Model(&buildCtx.Build).UpdateColumn("status", model.BuildSuccessful).Error; err != nil {
		logger.Error(module).Str("build_id", id.String()).Str("step", "execute").Err(err).Msg("could not update build")
	}
	logger.Debug(module).Str("build_id", id.String()).Str("step", "execute").Err(err).Msg("saved update to build")
}

func (ctx Context) Cancel(buildID uuid.UUID) error {
	logger.Debug(module).Str("build_id", buildID.String()).Str("step", "cancel").Msg("build cancel request")
	if err := db.Get().Model(&model.Build{ID: buildID}).UpdateColumn("status", model.BuildCanceled).Error; err != nil {
		logger.Error(module).Str("build_id", buildID.String()).Str("step", "cancel").Err(err).Msg("could not update build")
		return err
	}
	return nil
}

func buildContext(buildID uuid.UUID) (build.Context, error) {
	start := time.Now()
	logger.Debug(module).Str("build_id", buildID.String()).Str("step", "context-create").Msg("build context creation started")

	buildCtx, err := build.Init(buildID)
	id := buildCtx.Build.ID

	if err == nil {
		buildCtx.Build.Steps[0].Duration = time.Since(start)

		logger.Debug(module).Str("build_id", id.String()).Str("step", "context-create").Msg("build context creation succeeded")
		if err := db.Get().Create(&buildCtx.Build.Steps[0]).Error; err != nil {
			logger.Error(module).Str("build_id", id.String()).Str("step", "context-create").Err(err).Msg("could not write build steps")
			return buildCtx, ErrBuildSave
		}
		return buildCtx, nil
	}

	if err == build.ErrInvalidBuild {
		logger.Fatal(module).Str("build_id", id.String()).Str("step", "context-create").Err(err).Msg("fatal error during build context creation")
	}

	buildCtx.Build.Steps[0].Duration = time.Since(start)
	buildCtx.Build.Status = model.BuildFailed

	logger.Debug(module).Str("build_id", id.String()).Str("step", "context-create").Err(err).Msg("build context creation failed")
	if err := db.Get().Create(&buildCtx.Build.Steps[0]).Error; err != nil {
		logger.Error(module).Str("build_id", id.String()).Str("step", "context-create").Err(err).Msg("could not write build steps")
		return buildCtx, ErrBuildSave
	}
	return buildCtx, ErrBuildInitFailed
}

func bindWorker(buildCtx *build.Context) error {
	id := buildCtx.Build.ID
	logger.Debug(module).Str("build_id", id.String()).Str("step", "worker-bind").Msg("worker binding started")
	start := time.Now()
	buildCtx.Build.Steps = append(buildCtx.Build.Steps, model.BuildStep{Name: "Worker binding", BuildID: buildCtx.Build.ID, Logs: []model.BuildLog{}, Number: 1})

	// TODO: scheduler logic (adjusting for system & image req, balancing load, etc.) and handle cancel
	if err := db.Get().First(&buildCtx.Worker).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			build.AppendLog(buildCtx, 1, "BIND", "no worker available")
		} else {
			build.AppendLog(buildCtx, 1, "BIND", "db: "+err.Error())
		}
		buildCtx.Build.Steps[1].Duration = time.Since(start)
		buildCtx.Build.Status = model.BuildFailed
		logger.Error(module).Str("build_id", id.String()).Str("step", "worker-bind").Err(err).Msg("worker binding failed")
		if err := db.Get().Create(&buildCtx.Build.Steps[1]).Error; err != nil {
			logger.Error(module).Str("build_id", id.String()).Str("step", "worker-bind").Err(err).Msg("could not write build steps")
			return ErrBuildSave
		}
		return ErrBuildInitFailed
	}
	build.AppendLog(buildCtx, 1, "BIND", "worker ["+buildCtx.Worker.Name+"] bound")
	logger.Debug("scheduler").Msgf("worker: %+v", buildCtx.Worker)
	buildCtx.Build.Steps[1].Duration = time.Since(start)
	logger.Debug(module).Str("build_id", id.String()).Str("step", "worker-bind").Msg("worker binding succeeded")
	if err := db.Get().Create(&buildCtx.Build.Steps[1]).Error; err != nil {
		logger.Error(module).Str("build_id", id.String()).Str("step", "worker-bind").Err(err).Msg("could not write build steps")
		return ErrBuildSave
	}
	if err := db.Get().Model(&buildCtx.Build).UpdateColumns(&map[string]any{"status": model.BuildRunning, "worker_id": buildCtx.Worker.ID}).Error; err != nil {
		logger.Error(module).Str("build_id", id.String()).Str("step", "worker-bind").Err(err).Msg("could not update build")
		return ErrBuildSave
	}
	return nil
}
