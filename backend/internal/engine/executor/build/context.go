package build

import (
	"errors"

	"github.com/gg-mike/ci-cd-app/backend/internal/db"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Context struct {
	Build           model.Build
	Pipeline        model.Pipeline
	Project         model.Project
	GlobalVariables []model.Variable
	GlobalSecrets   []model.Secret
	Worker          model.Worker
}

var (
	ErrInvalidBuild     = errors.New("invalid build")
	ErrInvalidPipeline  = errors.New("invalid project")
	ErrInvalidProject   = errors.New("invalid pipeline")
	ErrInvalidSecrets   = errors.New("invalid secrets")
	ErrInvalidVariables = errors.New("invalid variables")
)

func Init(buildID uuid.UUID) (Context, error) {
	var err error
	ctx := Context{}

	// BUILD INIT
	if err = db.Get().First(&ctx.Build, "id = ?", buildID).Error; err != nil {
		return Context{}, ErrInvalidBuild
	}
	ctx.Build.Steps = []model.BuildStep{{Name: "Build context creation", BuildID: buildID, Logs: []model.BuildLog{}, Number: 0}}
	AppendLog(&ctx, 0, "BUILD INIT", "success")

	// PIPELINE INIT
	if err = db.Get().Preload(clause.Associations).First(&ctx.Pipeline, "id = ?", ctx.Build.PipelineID).Error; err != nil {
		AppendLog(&ctx, 0, "PIPELINE INIT", "db: "+err.Error())
		return ctx, ErrInvalidPipeline
	}
	AppendLog(&ctx, 0, "PIPELINE INIT", "success")

	// PROJECT INIT
	if err = db.Get().Preload(clause.Associations).First(&ctx.Project, "id = ?", ctx.Pipeline.ProjectID).Error; err != nil {
		AppendLog(&ctx, 0, "PROJECT INIT", "db: "+err.Error())
		return ctx, ErrInvalidProject
	}
	AppendLog(&ctx, 0, "PROJECT INIT", "success")

	// GLOBAL SECRETS INIT (CAN BE EMPTY)
	err = db.Get().Find(&ctx.GlobalSecrets, "project_id IS NULL AND pipeline_id IS NULL").Error
	switch err {
	case nil:
		AppendLog(&ctx, 0, "GLOBAL SECRETS INIT", "success")
	case gorm.ErrRecordNotFound:
		AppendLog(&ctx, 0, "GLOBAL SECRETS INIT", "no records found")
	default:
		AppendLog(&ctx, 0, "GLOBAL SECRETS INIT", "db: "+err.Error())
		return ctx, ErrInvalidSecrets
	}

	// GLOBAL VARIABLES INIT (CAN BE EMPTY)
	err = db.Get().Find(&ctx.GlobalVariables, "project_id IS NULL AND pipeline_id IS NULL").Error
	switch err {
	case nil:
		AppendLog(&ctx, 0, "GLOBAL VARIABLES INIT", "success")
	case gorm.ErrRecordNotFound:
		AppendLog(&ctx, 0, "GLOBAL VARIABLES INIT", "no records found")
	default:
		AppendLog(&ctx, 0, "GLOBAL VARIABLES INIT", "db: "+err.Error())
		return ctx, ErrInvalidVariables
	}

	return ctx, nil
}
