package cmd

import (
	"fmt"
	"net/url"
	"time"

	"github.com/marcsauter/gitlabctl/internal/auth"
	"github.com/marcsauter/gitlabctl/internal/output"
	"github.com/xanzy/go-gitlab"
	"github.com/zbindenren/king"
)

// Globals contains global arguments
type Globals struct {
	Hostname    string        `help:"Gitlab server." default:"gitlab.com"`
	AccessToken string        `help:"A gitlab API access token."`
	Timeout     time.Duration `help:"Timeout for API requests." default:"10s"`
	Debug       bool          `help:"Log debug output."`

	output.Kong

	Version    king.VersionFlag `help:"Show version information."`
	ShowConfig king.ShowConfig  `help:"Show config file used"`

	AuthTokenFilename string // set in main function
}

func (g *Globals) getAuthToken(hostname string) (string, error) {
	t, err := auth.New(g.AuthTokenFilename)
	if err != nil {
		return "", err
	}

	return t.Get(g.Hostname)
}

func (g *Globals) SetAuthToken(token string) error {
	at, err := auth.New(g.AuthTokenFilename)
	if err != nil {
		return err
	}

	return at.Set(g.Hostname, token)
}

func (g *Globals) DelAuthToken() error {
	t, err := auth.New(g.AuthTokenFilename)
	if err != nil {
		return err
	}

	return t.Remove(g.Hostname)
}

func (g *Globals) Client() (*gitlab.Client, error) {
	u, err := url.Parse(fmt.Sprintf("https://%s/api/v4", g.Hostname))
	if err != nil {
		return nil, err
	}

	token := g.AccessToken
	if token == "" {
		tok, err := g.getAuthToken(g.Hostname)
		if err != nil {
			return nil, err
		}

		token = tok
	}

	return gitlab.NewClient(token, gitlab.WithBaseURL(u.String()))
}
