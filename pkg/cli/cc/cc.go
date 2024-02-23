package cc

import (
	"context"
	"fmt"
	"os"

	"github.com/aserto-dev/clui"
	"github.com/aserto-dev/topaz/pkg/cli/cc/iostream"
	"github.com/aserto-dev/topaz/pkg/cli/dockerx"
	"github.com/samber/lo"
)

type CommonCtx struct {
	Context context.Context
	UI      *clui.UI
	NoCheck bool
}

type runStatus int

const (
	StatusNotRunning runStatus = iota
	StatusRunning
)

func NewCommonContext(noCheck bool) (*CommonCtx, error) {
	return &CommonCtx{
		Context: context.Background(),
		UI:      iostream.NewUI(iostream.DefaultIO()),
		NoCheck: noCheck,
	}, nil
}

func (c *CommonCtx) CheckRunStatus(containerName string, expectedStatus runStatus) bool {
	if c.NoCheck {
		return false
	}

	// set default container name if not specified.
	if containerName == "" {
		containerName = ContainerName()
	}

	running, err := dockerx.IsRunning(containerName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return false
	}

	return lo.Ternary(running, StatusRunning, StatusNotRunning) == expectedStatus
}
