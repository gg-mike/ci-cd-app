package scheduler

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/engine/executor"
	"github.com/gg-mike/ci-cd-app/backend/internal/engine/executor/build"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
)

func execute(ctx build.Context) {
	id := ctx.Build.ID.String()

	logger.Debug(module).Str("build_id", id).Str("step", "execute").Msg("build execution started")
	if err := executor.Execute(&ctx); err != nil {
		go (&Context{}).Finished(ctx)
		logger.Warn(module).Str("build_id", id).Str("step", "execute").Err(err).Msg("build execution ended with error")
		if err := db.Get().Model(&ctx.Build).UpdateColumn("status", model.BuildFailed).Error; err != nil {
			logger.Error(module).Str("build_id", id).Str("step", "execute").Err(err).Msg("could not update build")
			return
		}
		return
	}

	go (&Context{}).Finished(ctx)

	logger.Debug(module).Str("build_id", id).Str("step", "execute").Str("status", model.BuildSuccessful.String()).Msg("build execution ended")
	if err := db.Get().Model(&ctx.Build).UpdateColumn("status", model.BuildSuccessful).Error; err != nil {
		logger.Error(module).Str("build_id", id).Str("step", "execute").Err(err).Msg("could not update build")
		return
	}
	logger.Debug(module).Str("build_id", id).Str("step", "execute").Msg("saved update to build")
}
