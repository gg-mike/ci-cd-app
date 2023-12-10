package scheduler

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/engine/executor/build"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"gorm.io/gorm"
)

func finished(ctx build.Context) {
	logger.Debug(module).Str("worker_id", ctx.Worker.ID.String()).Msg("decrement active builds on bound worker")
	if err := db.Get().Transaction(func(tx *gorm.DB) error {
		var worker model.Worker
		if err := tx.Where(&model.Worker{ID: ctx.Worker.ID}).First(&worker).Error; err != nil {
			return err
		}
		if err := tx.Model(&worker).UpdateColumn("active_builds", worker.ActiveBuilds-1).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		logger.Error(module).Str("worker_id", ctx.Worker.ID.String()).Err(err).Msg("could not decrement active builds counter")
		return
	}
	go (&Context{}).ChangeInWorkers()
}
