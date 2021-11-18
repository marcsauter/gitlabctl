package gitlabctl

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/marcsauter/gitlabctl/internal/cmd"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh/terminal"
)

type authCmd struct {
	Login  loginAuthCmd  `cmd:"" help:"Login ..."`
	Logout logoutAuthCmd `cmd:"" help:"Login ..."`
}

type loginAuthCmd struct{}

func (li loginAuthCmd) Run(app *kong.Context, g *cmd.Globals, l *zap.SugaredLogger) error {
	fmt.Print("Enter Token: ")

	tkn, err := terminal.ReadPassword(0)
	if err != nil {
		return err
	}

	return g.SetAuthToken(strings.TrimSpace(string(tkn)))
}

type logoutAuthCmd struct{}

func (lo logoutAuthCmd) Run(app *kong.Context, g *cmd.Globals, l *zap.SugaredLogger) error {
	return g.DelAuthToken()
}
