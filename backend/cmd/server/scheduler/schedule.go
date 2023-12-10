package scheduler

import (
	"github.com/gg-mike/ci-cd-app/backend/internal/engine/executor/build"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/google/uuid"
)

func schedule(buildID uuid.UUID) {
	ctx, err := buildContext(buildID)
	if err != nil {
		if err != build.ErrInvalidBuild {
			logger.Fatal(module).Str("build_id", buildID.String()).Str("step", "context-create").Err(err).Msg("fatal error during build context creation")
		}
		return
	}

	go (&Context{}).AddToQueue(ctx)
}
