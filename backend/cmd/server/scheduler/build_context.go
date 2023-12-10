package scheduler

import (
	"time"

	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/engine/executor/build"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/google/uuid"
)

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
		logger.Warn(module).Str("build_id", id.String()).Str("step", "context-create").Err(err).Msg("build not available - rescheduling")
		go func() {
			time.Sleep(time.Second)
			go (&Context{}).Schedule(buildID)
		}()
		return buildCtx, build.ErrInvalidBuild
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
