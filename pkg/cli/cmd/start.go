package cmd

import (
	"os"
	"path"

	"github.com/aserto-dev/topaz/pkg/cc/config"
	"github.com/aserto-dev/topaz/pkg/cli/cc"
	"github.com/aserto-dev/topaz/pkg/cli/dockerx"
	"github.com/pkg/errors"

	"github.com/fatih/color"
)

type StartCmd struct {
	StartRunCmd
}

func (cmd *StartCmd) Run(c *cc.CommonCtx) error {
	if running, err := dockerx.IsRunning(dockerx.Topaz); running || err != nil {
		if !running {
			return ErrNotRunning
		}
		if err != nil {
			return err
		}
	}

	rootPath, err := dockerx.DefaultRoots()
	if err != nil {
		return err
	}

	color.Green(">>> starting topaz...")
	cmdX := cmd.StartRunCmd
	args, err := cmdX.dockerArgs(rootPath, false)
	if err != nil {
		return err
	}

	cmdArgs := []string{
		"run",
		"--config-file", "/config/config.yaml",
	}

	args = append(args, cmdArgs...)

	if _, err := os.Stat(path.Join(rootPath, "cfg", "config.yaml")); errors.Is(err, os.ErrNotExist) {
		return errors.Errorf("%s does not exist, please run 'topaz configure'", path.Join(rootPath, "cfg", "config.yaml"))
	}

	generator := config.NewGenerator("config.yaml")
	if _, err := generator.CreateCertsDir(); err != nil {
		return err
	}

	if _, err := generator.CreateDataDir(); err != nil {
		return err
	}

	return dockerx.DockerWith(cmdX.env(), args...)
}
