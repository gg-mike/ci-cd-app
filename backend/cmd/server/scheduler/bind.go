package scheduler

import (
	"github.com/gg-mike/ci-cd-app/backend/cmd/server/queue"
	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"gorm.io/gorm"
)

func bind() {
	if err := db.Get().Transaction(func(tx *gorm.DB) error {
		var q []queue.QueueElem
		if err := db.Get().Find(&q).Error; err != nil {
			logger.Error(module).Str("step", "worker-bind").Err(err).Msg("error during queue loading")
			return err
		}

		if len(q) == 0 {
			return nil
		}

		for _, elem := range q {
			var workers []model.Worker
			if err := tx.Model(&model.Worker{}).Find(&workers, &model.Worker{Status: model.WorkerIdle}).Error; err != nil {
				return err
			}

			worker, err := SelectWorker(elem.Context.Pipeline.Config, workers)

			if err == ErrNoAvailableWorker {
				logger.Debug(module).Str("build_id", elem.Context.Build.ID.String()).Str("step", "worker-bind").Msg("no worker available")
				return nil
			} else if err == ErrNoAvailableWorkerForConfiguration {
				logger.Debug(module).Str("build_id", elem.Context.Build.ID.String()).Str("step", "worker-bind").Msg("no worker available for given configuration")
				continue
			}

			if err := bindWorker(&elem, worker); err != nil {
				logger.Error(module).Str("build_id", elem.Context.Build.ID.String()).Str("step", "worker-bind").Err(err).Msg("error during worker binding")
				if err := db.Get().Model(&elem.Context.Build).UpdateColumn("status", model.BuildFailed).Error; err != nil {
					logger.Error(module).Str("build_id", elem.Context.Build.ID.String()).Str("step", "worker-bind").Err(err).Msg("could not update build")
					return err
				}
				return err
			} else {
				if err := db.Get().Delete(&queue.QueueElem{}, elem.ID).Error; err != nil {
					return err
				}
				go execute(elem.Context)
			}

		}
		return nil
	}); err != nil {
		logger.Error(module).Str("step", "worker-bind").Err(err).Msg("error during worker binding")
	}
}
