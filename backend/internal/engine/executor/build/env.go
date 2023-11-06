package build

import (
	"encoding/base64"
	"fmt"
	"maps"
	"regexp"
	"strings"

	"github.com/gg-mike/ci-cd-app/backend/internal/logger"
	"github.com/gg-mike/ci-cd-app/backend/internal/model"
	"github.com/gg-mike/ci-cd-app/backend/internal/vault"
	"github.com/rs/zerolog"
)

type envInstance struct {
	value string
	path  string
}

func (ctx *Context) CreateEnvSteps() error {
	workdirSteps, workdirCleanup := ctx.createWorkdirStep()
	secretsSteps, secretsCleanup, err := ctx.createSecretsStep()
	logger.Basic(zerolog.DebugLevel, "build").Msgf("%+v", secretsSteps)
	if err != nil {
		return err
	}
	variablesSteps, variablesCleanup, err := ctx.createVariablesStep()
	if err != nil {
		return err
	}

	ctx.Pipeline.Config.Steps = append([]model.PipelineConfigStep{
		workdirSteps, secretsSteps, variablesSteps,
	}, ctx.Pipeline.Config.Steps...)

	ctx.Pipeline.Config.Cleanup = append(ctx.Pipeline.Config.Cleanup, workdirCleanup...)
	ctx.Pipeline.Config.Cleanup = append(ctx.Pipeline.Config.Cleanup, secretsCleanup...)
	ctx.Pipeline.Config.Cleanup = append(ctx.Pipeline.Config.Cleanup, variablesCleanup...)

	return nil
}

func (ctx Context) createWorkdirStep() (model.PipelineConfigStep, []string) {
	return model.PipelineConfigStep{
		Name:     "Work dir setup",
		Commands: []string{"mkdir -p workdir", "cd workdir"},
	}, []string{"cd ~", "rm -rf workdir"}
}

func (ctx Context) createSecretsStep() (model.PipelineConfigStep, []string, error) {
	secrets := map[string]envInstance{}
	groups := [][]model.Secret{ctx.GlobalSecrets, ctx.Project.Secrets, ctx.Pipeline.Secrets}

	for _, group := range groups {
		for _, secret := range group {
			value, err := vault.Str(secret.ID.String())
			if err != nil {
				return model.PipelineConfigStep{}, []string{}, err
			}
			secrets[secret.Key] = envInstance{value, secret.Path}
		}
	}

	commands, cleanUpCommands, err := prepareStepCommands(ctx.Pipeline.Config.System, secrets, "_")
	if err != nil {
		return model.PipelineConfigStep{}, []string{}, err
	}

	return model.PipelineConfigStep{Name: "Secret exports", Commands: commands}, cleanUpCommands, nil
}

func (ctx Context) createVariablesStep() (model.PipelineConfigStep, []string, error) {
	variables := map[string]envInstance{}

	variables["__PROJECT_NAME"] = envInstance{ctx.Project.Name, ""}
	variables["__REPO"] = envInstance{ctx.Project.Repo, ""}

	variables["__PIPELINE_NAME"] = envInstance{ctx.Pipeline.Name, ""}
	variables["__BRANCH"] = envInstance{ctx.Pipeline.Branch, ""}
	maps.Copy(variables, setRepoVariables(ctx.Project.Repo))

	groups := [][]model.Variable{ctx.GlobalVariables, ctx.Project.Variables, ctx.Pipeline.Variables}

	for _, group := range groups {
		for _, variable := range group {
			variables[variable.Key] = envInstance{variable.Value, variable.Path}
		}
	}

	commands, cleanUpCommands, err := prepareStepCommands(ctx.Pipeline.Config.System, variables, "")
	if err != nil {
		return model.PipelineConfigStep{}, []string{}, err
	}

	return model.PipelineConfigStep{Name: "Variable exports", Commands: commands}, cleanUpCommands, nil
}

func prepareStepCommands(system string, env map[string]envInstance, prefix string) ([]string, []string, error) {
	var templateEnv, templateFile, templateFileDelete string
	// TODO: support for over OS
	if system == "Linux" {
		templateEnv = "export %s%s=\"%s\""
		templateFile = "export %s%s=\"%s\" && echo \"%s\" > %s"
		templateFileDelete = "rm -f %s"
	}
	commands := []string{}
	cleanUpCommands := []string{}
	for k, v := range env {
		if v.path != "" {
			value, err := base64.StdEncoding.DecodeString(v.value)
			if err != nil {
				return []string{}, []string{}, err
			}
			commands = append(commands, fmt.Sprintf(templateFile, prefix, k, v.path, value, v.path))
			cleanUpCommands = append(cleanUpCommands, fmt.Sprintf(templateFileDelete, v.path))
		} else {
			commands = append(commands, fmt.Sprintf(templateEnv, prefix, k, v.value))
		}
	}
	return commands, cleanUpCommands, nil
}

func setRepoVariables(repo string) map[string]envInstance {
	// TODO: support for over providers
	if repo == "" {
		return map[string]envInstance{}
	} else if strings.Contains(repo, "git@github.com") {
		re := regexp.MustCompile(`git@github\.com:(?P<owner>[\w\d\.\-_]+)\/(?P<name>[\w\d\.\-_]+)\.git`)
		matches := re.FindStringSubmatch(repo)
		return map[string]envInstance{
			"__GITHUB_OWNER": {matches[re.SubexpIndex("owner")], ""},
			"__GITHUB_NAME":  {matches[re.SubexpIndex("name")], ""},
		}
	} else if strings.Contains(repo, "https://github.com") {
		re := regexp.MustCompile(`https:\/\/github\.com\/(?P<owner>[\w\d\.\-_]+)\/(?P<name>[\w\d\.\-_]+)`)
		matches := re.FindStringSubmatch(repo)
		return map[string]envInstance{
			"__GITHUB_OWNER": {matches[re.SubexpIndex("owner")], ""},
			"__GITHUB_NAME":  {matches[re.SubexpIndex("name")], ""},
		}
	}
	return map[string]envInstance{}
}
