package main

import (
	"os/user"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/marcsauter/gitlabctl/internal/cmd/gitlabctl"
	"github.com/postfinance/flash"
	"github.com/zbindenren/king"
)

const (
	TokenFile = ".gitlabctl"
)

//nolint:gochecknoglobals // there is no other way to initialize these values
var (
	version = "0.0.0"
	commit  = "12345678"
	date    string
)

func main() {
	l := flash.New(flash.WithColor())

	u, err := user.Current()
	if err != nil {
		l.Fatal(err)
	}

	b, err := king.NewBuildInfo(version,
		king.WithDateString(date),
		king.WithRevision(commit),
	)
	if err != nil {
		l.Fatal(err)
	}

	cli := gitlabctl.CLI{}

	app := kong.Parse(&cli, king.DefaultOptions(
		king.Config{
			Name:        "gitlabctl",
			Description: "Command line tool for Gitlab.",
			BuildInfo:   b,
		},
	)...)

	cli.Globals.AuthTokenFilename = filepath.Join(u.HomeDir, TokenFile)

	l.SetDebug(cli.Debug)

	if err := app.Run(&cli.Globals, l.Get()); err != nil {
		l.Fatal(err)
	}
}
