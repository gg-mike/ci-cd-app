package executor

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/engine/executor/build"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/rs/zerolog"
)

func Execute(buildCtx *build.Context) error {
	ctx, err := Init(buildCtx)
	defer ctx.CloseSession()
	defer ctx.CloseConn()
	if err != nil {
		return err
	}

	if err = buildCtx.CreateEnvSteps(); err != nil {
		return err
	}

	steps := buildCtx.Pipeline.Config.Steps

	x := len(buildCtx.Build.Steps)

	for i, step := range steps {
		// TODO: refactor
		if err = ctx.runStep(buildCtx, step, i, x); err != nil {
			break
		}
	}

	ctx.CloseSession()
	ctx.CloseConn()
	ctx, err = Init(buildCtx)
	if err != nil {
		return err
	}

	if errCleanup := ctx.runCleanup(buildCtx, x); errCleanup != nil {
		logger.Basic(zerolog.DebugLevel, "executor").Err(errCleanup).Msg("err in cleanup")
		return err
	}

	logger.Basic(zerolog.DebugLevel, "executor").Err(err).Msg("end")

	return err
}

func (ctx *Context) runStep(buildCtx *build.Context, step model.PipelineConfigStep, i, x int) error {
	fmt.Println(step.Name)
	if err := db.Get().First(&buildCtx.Build).Error; err != nil {
		return err
	}
	logger.Basic(zerolog.DebugLevel, "executor").Msgf("build status: %s", buildCtx.Build.Status)
	if buildCtx.Build.Status == model.BuildCanceled {
		return nil
	}
	start := time.Now()
	// TODO: Live feed support
	buildStep := model.BuildStep{Name: step.Name, BuildID: buildCtx.Build.ID, Logs: []model.BuildLog{}, Number: i + x}

	go ctx.runCommands(step.Commands)

	running := true
	var err error
	for running {
		select {
		case cmd := <-ctx.CmdChan:
			if step.Name == "Secret exports" {
				continue
			}
			fmt.Printf("\033[32m[%d/%d] $ %s\033[0m\n", cmd.Idx+1, cmd.Total, cmd.Command)
			cmd.BuildStepID = buildStep.ID
			buildStep.Logs = append(buildStep.Logs, cmd)
			// TODO: Live feed support
		case out := <-ctx.OutChan:
			if step.Name == "Secret exports" {
				continue
			}
			fmt.Println(out)
			if len(buildStep.Logs) == 0 {
				continue
			}
			if buildStep.Logs[len(buildStep.Logs)-1].Output != "" {
				buildStep.Logs[len(buildStep.Logs)-1].Output += "\n"
			}
			buildStep.Logs[len(buildStep.Logs)-1].Output += out
			// TODO: Live feed support
		case err = <-ctx.ErrChan:
			running = false
		}
	}

	buildStep.Duration = time.Since(start)
	if err := db.Get().Create(&buildStep).Error; err != nil {
		return err
	}
	if err != nil {
		if err := db.Get().Model(&buildCtx.Build).UpdateColumn("status", model.BuildFailed).Error; err != nil {
			return err
		}
		return err
	}
	return nil
}

func (ctx *Context) runCleanup(buildCtx *build.Context, number int) error {
	fmt.Println("Cleanup")
	start := time.Now()
	// TODO: Live feed support
	buildStep := model.BuildStep{Name: "Cleanup", BuildID: buildCtx.Build.ID, Logs: []model.BuildLog{}, Number: number}

	go ctx.runCommands(buildCtx.Pipeline.Config.Cleanup)

	running := true
	var err error
	for running {
		select {
		case cmd := <-ctx.CmdChan:
			fmt.Printf("\033[32m[%d/%d] $ %s\033[0m\n", cmd.Idx+1, cmd.Total, cmd.Command)
			cmd.BuildStepID = buildStep.ID
			buildStep.Logs = append(buildStep.Logs, cmd)
			// TODO: Live feed support
		case out := <-ctx.OutChan:
			fmt.Println(out)
			if len(buildStep.Logs) == 0 {
				continue
			}
			if buildStep.Logs[len(buildStep.Logs)-1].Output != "" {
				buildStep.Logs[len(buildStep.Logs)-1].Output += "\n"
			}
			buildStep.Logs[len(buildStep.Logs)-1].Output += out
			// TODO: Live feed support
		case err = <-ctx.ErrChan:
			running = false
		}
	}

	buildStep.Duration = time.Since(start)
	if err := db.Get().Create(&buildStep).Error; err != nil {
		return err
	}
	if err != nil {
		if err := db.Get().Model(&buildCtx.Build).UpdateColumn("status", model.BuildFailed).Error; err != nil {
			return err
		}
		return err
	}
	return nil
}

func (ctx *Context) runCommands(commands []string) {
	OUT_CMD_TERM := "Ua&&Bi9G*TjbPF62oGa4"
	ERR_CMD_TERM := "!N3o#F4SPZ&UDxybohUT"

	total := len(commands)
	in := make(chan string)
	inStatus := make(chan error)

	outReader := SyncReader{
		source:     ctx.Reader,
		term:       OUT_CMD_TERM,
		errTerm:    ERR_CMD_TERM,
		ready:      make(chan any),
		scanStatus: make(chan error),
	}

	defer close(in)
	defer close(outReader.ready)

	go func() {
		for command := range in {
			cmd := []byte(fmt.Sprintf("%s 2>&1 && echo '%s' || echo '%s'\n", command, OUT_CMD_TERM, ERR_CMD_TERM))
			_, err := ctx.Writer.Write(cmd)
			inStatus <- err
		}
	}()

	go scan(outReader, ctx.OutChan)

	// Read welcome message
	if err := runCommand("echo", in, inStatus, outReader); err != nil {
		fmt.Printf("ok1.1 err: %v\n", err)
		// Error during connection like 'mesg: ttyname failed: Inappropriate ioctl for device'
		if err.Error() != "command ended with error" {
			ctx.ErrChan <- err
		}
	}

	// Run commands
	for i, command := range commands {
		ctx.CmdChan <- model.BuildLog{Command: command, Idx: i, Total: total}
		if err := runCommand(command, in, inStatus, outReader); err != nil {
			ctx.ErrChan <- err
		}
	}
	ctx.ErrChan <- nil
}

type SyncReader struct {
	source     io.Reader
	term       string
	errTerm    string
	ready      chan any
	scanStatus chan error
}

func scan(reader SyncReader, outChan chan string) {
	scanner := bufio.NewScanner(reader.source)
	for range reader.ready {
		for {
			if tkn := scanner.Scan(); tkn {
				text := scanner.Text()
				if strings.Contains(text, reader.term) {
					break
				} else if strings.Contains(text, reader.errTerm) {
					reader.scanStatus <- errors.New("command ended with error")
					return
				} else {
					outChan <- text
				}
			} else {
				reader.scanStatus <- scanner.Err()
				return
			}
		}
		reader.scanStatus <- nil
	}
}

func runCommand(command string, in chan string, inStatus chan error, reader SyncReader) error {
	in <- command
	if err := <-inStatus; err != nil {
		return err
	}
	reader.ready <- true
	if err := <-reader.scanStatus; err != nil {
		return err
	}
	return nil
}
