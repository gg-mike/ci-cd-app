package executor

import (
	"io"

	"github.com/gg-mike/ci-cd-app/backend/internal/engine/executor/build"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gg-mike/ci-cd-app/backend/internal/ssh"
	"github.com/gg-mike/ci-cd-app/backend/internal/vault"
)

type Context struct {
	Writer       io.Writer
	Reader       io.Reader
	CmdChan      chan model.BuildLog
	OutChan      chan string
	ErrChan      chan error
	CloseConn    func() error
	CloseSession func() error
}

func Init(buildCtx *build.Context) (Context, error) {
	var err error
	ctx := Context{
		CloseConn:    func() error { return nil },
		CloseSession: func() error { return nil },
	}
	privateKey, err := vault.Str(buildCtx.Worker.ID.String())
	if err != nil {
		return ctx, err
	}
	conn, err := ssh.CreateConnection(buildCtx.Worker.Username, buildCtx.Worker.Address, privateKey)
	if err != nil {
		return ctx, err
	}
	ctx.CloseConn = conn.Close

	session, err := conn.NewSession()
	if err != nil {
		return ctx, err
	}
	ctx.CloseSession = session.Close

	ctx.Writer, err = session.StdinPipe()
	if err != nil {
		return ctx, err
	}
	ctx.Reader, err = session.StdoutPipe()
	if err != nil {
		return ctx, err
	}

	if err = session.Shell(); err != nil {
		return ctx, err
	}

	ctx.CmdChan = make(chan model.BuildLog)
	ctx.OutChan = make(chan string)
	ctx.ErrChan = make(chan error)

	return ctx, nil
}
