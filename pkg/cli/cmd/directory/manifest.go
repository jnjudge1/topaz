package directory

import (
	"fmt"
	"io"
	"os"

	"github.com/aserto-dev/topaz/pkg/cli/cc"
	"github.com/aserto-dev/topaz/pkg/cli/clients"
	"github.com/aserto-dev/topaz/pkg/cli/cmd/common"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

type GetManifestCmd struct {
	Path   string `arg:"path" help:"filepath to manifest file" type:"path" optional:""`
	Stdout bool   `flag:"" help:"output manifest to --stdout"`
	clients.DirectoryConfig
}

type SetManifestCmd struct {
	Path  string `arg:"" help:"filepath to manifest file" type:"path" optional:""`
	Stdin bool   `flag:"" help:"set manifest from --stdin"`
	clients.DirectoryConfig
}

type DeleteManifestCmd struct {
	Force bool `flag:"" help:"do not ask for conformation to delete manifest"`
	clients.DirectoryConfig
}

func (cmd *GetManifestCmd) Run(c *cc.CommonCtx) error {
	if !c.IsServing(cmd.Host) {
		return errors.Wrap(cc.ErrNotServing, cmd.Host)
	}
	dirClient, err := clients.NewDirectoryClient(c, &cmd.DirectoryConfig)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprint(color.Error, color.GreenString(">>> get manifest to %s", cmd.Path))

	r, err := dirClient.V3.GetManifest(c.Context)
	if err != nil {
		return err
	}

	w := os.Stdout

	if cmd.Path != "" {
		w, err = os.Create(cmd.Path)
		if err != nil {
			return err
		}
	}
	defer w.Close()

	if _, err := io.Copy(w, r); err != nil {
		return err
	}

	return nil
}

func (cmd *SetManifestCmd) Run(c *cc.CommonCtx) error {
	if !c.IsServing(cmd.Host) {
		return errors.Wrap(cc.ErrNotServing, cmd.Host)
	}
	dirClient, err := clients.NewDirectoryClient(c, &cmd.DirectoryConfig)
	if err != nil {
		return err
	}

	r := os.Stdin
	if cmd.Path != "" {
		r, err = os.Open(cmd.Path)
		if err != nil {
			return err
		}
	}

	_, _ = fmt.Fprint(color.Error, color.GreenString(">>> set manifest to %s", cmd.Path))

	return dirClient.V3.SetManifest(c.Context, r)
}

func (cmd *DeleteManifestCmd) Run(c *cc.CommonCtx) error {
	if !c.IsServing(cmd.Host) {
		return errors.Wrap(cc.ErrNotServing, cmd.Host)
	}
	dirClient, err := clients.NewDirectoryClient(c, &cmd.DirectoryConfig)
	if err != nil {
		return err
	}

	c.UI.Exclamation().Msg("WARNING: delete manifest resets all directory state, including relation and object data")
	if cmd.Force || common.PromptYesNo("Do you want to continue?", false) {
		_, _ = fmt.Fprint(color.Error, color.GreenString(">>> delete manifest"))
		return dirClient.V3.DeleteManifest(c.Context)
	}
	return nil
}
