// Package gitlabctl is a general purpose CLI for Gitlab
package gitlabctl

import (
	"github.com/marcsauter/gitlabctl/internal/cmd"
	"github.com/marcsauter/gitlabctl/internal/cmd/project"
)

// CLI ist the dopatch command line interface
type CLI struct {
	Auth         authCmd         `cmd:"" help:"Authenticate with your gitlab instance."`
	MergeRequest mergeRequestCmd `cmd:"" help:""`
	Issue        issueCmd        `cmd:"" help:""`

	project.Cmd
	cmd.Globals
}
