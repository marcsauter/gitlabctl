package project

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/marcsauter/gitlabctl/internal/cmd"
	"github.com/marcsauter/gitlabctl/internal/options"
	"github.com/marcsauter/gitlabctl/internal/output"
	"github.com/xanzy/go-gitlab"
	"go.uber.org/zap"
)

type listProjectCmd struct {
	Kind               string    `help:"Kind of project." enum:"all,user,group" default:"all"`
	LastActivityAfter  time.Time `help:"Limit results to projects with last_activity after specified time. Format: ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)"`
	LastActivityBefore time.Time `help:"Limit results to projects with last_activity before specified time. Format: ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)"`
	Membership         bool      `help:"Limit by projects that the current user is a member of."`
	Owned              bool      `help:"Limit by projects explicitly owned by the current user."`
	Search             string    `help:"Return list of projects matching the search criteria."`
	Starred            bool      `help:"Limit by projects starred by the current user."`
}

func (lpc listProjectCmd) Run(app *kong.Context, g *cmd.Globals, l *zap.SugaredLogger) error {
	clnt, err := g.Client()
	if err != nil {
		return err
	}

	projects := []*gitlab.Project{}

	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
	}

	if err := options.Transfer(&lpc, opt); err != nil {
		return err // TODO: panic?
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)

		proj, resp, err := clnt.Projects.ListProjects(opt, gitlab.WithContext(ctx))
		cancel()

		if err != nil {
			return fmt.Errorf("failed to list projects: %w", err)
		}

		_ = resp.Body.Close()

		switch lpc.Kind {
		case "all":
			projects = append(projects, proj...)
		default:
			for _, p := range proj {
				if p.Namespace.Kind == lpc.Kind {
					projects = append(projects, p)
				}
			}
		}

		// Exit the loop when we've seen all pages.
		// https://docs.gitlab.com/ee/user/gitlab_com/index.html#pagination-response-headers
		if len(proj) == 0 || (resp.TotalPages > 0 && resp.CurrentPage >= resp.TotalPages) {
			break
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}

	out, err := output.New(g.Format, projects, []string{"ID", "Name", "Visibility", "Namespace.Kind", "HTTPURLToRepo"})
	if err != nil {
		return err
	}

	return out.Print(os.Stdout)
}
