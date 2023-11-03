package engine

import (
	"fmt"
	"strings"

	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gg-mike/ci-cd-app/backend/internal/vault"
)

func CreateEnvSecretsBuildStep(ctx *BuildContext) (model.PipelineConfigStep, error) {
	step := model.PipelineConfigStep{Name: "Secret exports", Commands: []string{}}
	secrets := map[string]string{}

	for _, secret := range ctx.GlobalSecrets {
		value, err := vault.Str(secret.ID.String())
		if err != nil {
			return model.PipelineConfigStep{}, err
		}
		secrets[secret.Key] = value
	}

	for _, secret := range ctx.Project.Secrets {
		value, err := vault.Str(secret.ID.String())
		if err != nil {
			return model.PipelineConfigStep{}, err
		}
		secrets[secret.Key] = value
	}

	for _, secret := range ctx.Pipeline.Secrets {
		value, err := vault.Str(secret.ID.String())
		if err != nil {
			return model.PipelineConfigStep{}, err
		}
		secrets[secret.Key] = value
	}

	var template string
	// TODO: support for over OS
	if ctx.Pipeline.Config.System == "Linux" {
		template = "export %s=\"%s\""
	}
	for _, secret := range secrets {
		step.Commands = append(step.Commands, fmt.Sprintf(template, secret, secrets[secret]))
	}

	return step, nil
}

func CreateEnvVariablesBuildStep(ctx *BuildContext) model.PipelineConfigStep {
	step := model.PipelineConfigStep{Name: "Variable exports", Commands: []string{}}
	variables := map[string]string{}

	variables["__PROJECT_NAME"] = ctx.Project.Name
	variables["__REPO"] = ctx.Project.Repo
	projectUrl := strings.Split(ctx.Project.Repo, "/")
	variables["__REPO_DIR"] = projectUrl[len(projectUrl)-1]

	variables["__PIPELINE_NAME"] = ctx.Pipeline.Name
	variables["__BRANCH"] = ctx.Pipeline.Branch

	for _, variable := range ctx.GlobalVariables {
		variables[variable.Key] = variable.Value
	}

	for _, variable := range ctx.Project.Variables {
		variables[variable.Key] = variable.Value
	}

	for _, variable := range ctx.Pipeline.Variables {
		variables[variable.Key] = variable.Value
	}

	var template string
	// TODO: support for over OS
	if ctx.Pipeline.Config.System == "Linux" {
		template = "export %s=\"%s\""
	}
	for k, v := range variables {
		step.Commands = append(step.Commands, fmt.Sprintf(template, k, v))
	}

	return step
}
