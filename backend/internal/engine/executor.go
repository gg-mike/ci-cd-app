package engine

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gg-mike/ci-cd-app/backend/internal/ssh"
	"github.com/gg-mike/ci-cd-app/backend/internal/vault"
	"github.com/rs/zerolog"
)

func Execute(ctx *BuildContext) error {
	privateKey, err := vault.Str(ctx.Worker.ID.String())
	if err != nil {
		return err
	}
	conn, err := ssh.CreateConnection(ctx.Worker.Username, ctx.Worker.Address, privateKey)
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	err = session.Shell()
	if err != nil {
		return err
	}

	steps := ctx.Pipeline.Config.Steps
	secretStep, err := CreateEnvSecretsBuildStep(ctx)
	if err != nil {
		return err
	}
	variableStep := CreateEnvVariablesBuildStep(ctx)

	workDirStep := model.PipelineConfigStep{
		Name: "Work dir setup",
		Commands: []string{
			"rm -rf workspace",
			"mkdir workspace",
			"cd workspace",
		},
	}

	steps = append([]model.PipelineConfigStep{secretStep, variableStep, workDirStep}, steps...)

	cmdChan := make(chan model.BuildLog)
	outChan := make(chan string)
	errChan := make(chan error)

	x := len(ctx.Build.Steps)

	// TODO: Cleanup should be done even in the case of crash (new field in config, like Cleanup)
	for i, step := range steps {
		if err := db.Get().First(&ctx.Build).Error; err != nil {
			return err
		}
		logger.Basic(zerolog.DebugLevel, "executor").Msgf("build status: %s", ctx.Build.Status)
		if ctx.Build.Status == model.BuildCanceled {
			return nil
		}
		start := time.Now()
		// TODO: Live feed support
		buildStep := model.BuildStep{Name: step.Name, BuildID: ctx.Build.ID, Logs: []model.BuildLog{}, Number: i + x}

		go runCommands(stdin, stdout, step.Commands, cmdChan, outChan, errChan)

		running := true
		var err error
		for running {
			select {
			case cmd := <-cmdChan:
				if step.Name == "Secret exports" {
					continue
				}
				fmt.Printf("\033[32m[%d/%d] $ %s\033[0m\n", cmd.Idx+1, cmd.Total, cmd.Command)
				cmd.BuildStepID = buildStep.ID
				buildStep.Logs = append(buildStep.Logs, cmd)
				// TODO: Live feed support
			case out := <-outChan:
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
			case err = <-errChan:
				running = false
			}
		}

		buildStep.Duration = time.Since(start)
		if err := db.Get().Create(&buildStep).Error; err != nil {
			return err
		}
		if err != nil {
			if err := db.Get().Model(&ctx.Build).UpdateColumn("status", model.BuildFailed).Error; err != nil {
				return err
			}
			return err
		}
	}

	return nil
}

func runCommands(
	writer io.Writer,
	reader io.Reader,
	commands []string,
	cmdChan chan model.BuildLog,
	outChan chan string,
	errChan chan error,
) {
	OUT_CMD_TERM := "Ua&&Bi9G*TjbPF62oGa4"
	ERR_CMD_TERM := "!N3o#F4SPZ&UDxybohUT"

	total := len(commands)
	in := make(chan string)
	inStatus := make(chan error)

	outReader := SyncReader{
		source:     reader,
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
			_, err := writer.Write(cmd)
			inStatus <- err
		}
	}()

	go scan(outReader, outChan)

	// Read welcome message
	if err := runCommand("echo", in, inStatus, outReader); err != nil {
		// Error during connection like 'mesg: ttyname failed: Inappropriate ioctl for device'
		if err.Error() != "command ended with error" {
			errChan <- err
		}
	}

	// Run commands
	for i, command := range commands {
		cmdChan <- model.BuildLog{Command: command, Idx: i, Total: total}
		if err := runCommand(command, in, inStatus, outReader); err != nil {
			errChan <- err
		}
	}
	errChan <- nil
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
