package gitlabctl

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/marcsauter/gitlabctl/internal/cmd"
	"github.com/marcsauter/gitlabctl/internal/output"
	"github.com/xanzy/go-gitlab"
	"go.uber.org/zap"
)

type listIssueCmd struct {
	ProjectID int      `help:"Project ID."`
	State     string   `help:"Issue state." enum:"all,opened,closed" default:"opened"`
	Labels    []string `help:"Issue lables"`
}

func (li listIssueCmd) Run(app *kong.Context, g *cmd.Globals, l *zap.SugaredLogger) error {
	clnt, err := g.Client()
	if err != nil {
		return err
	}

	issues := []*gitlab.Issue{}

	opt := &gitlab.ListProjectIssuesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
		State:  &li.State,
		Labels: li.Labels,
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)

		i, resp, err := clnt.Issues.ListProjectIssues(li.ProjectID, opt, gitlab.WithContext(ctx))

		cancel()

		if err != nil {
			return fmt.Errorf("failed to list issues: %w", err)
		}

		_ = resp.Body.Close()

		issues = append(issues, i...)

		// Exit the loop when we've seen all pages.
		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		// Update the page number to get the next page.
		opt.Page = resp.NextPage
	}

	/*
		for _, i := range issues {
			assignee := "NONE"

			if i.Assignee != nil {
				assignee = i.Assignee.Username
			}

			fmt.Fprintf(w, "%d\t%q\t%s\t%s\t%s\n", i.IID, i.Title, assignee, i.State, i.WebURL)
		}
	*/

	out, err := output.New(g.Format, issues, []string{"IID", "Title", "Assignee", "State", "WebURL"})
	if err != nil {
		return err
	}

	return out.Print(os.Stdout)
}
