package build

import "github.com/gg-mike/ci-cd-app/backend/internal/model"

func AppendLog(ctx *Context, idx int, command, output string) {
	ctx.Build.Steps[idx].Logs = append(ctx.Build.Steps[idx].Logs, model.BuildLog{Command: command, Output: output})
}
