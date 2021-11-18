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

type listMergeRequestCmd struct {
	ProjectID int      `help:"Project ID."`
	State     string   `help:"Issue state." enum:"opened,closed,locked,merged" default:"opened"`
	Labels    []string `help:"Merge request lables"`
}

func (lmr listMergeRequestCmd) Run(app *kong.Context, g *cmd.Globals, l *zap.SugaredLogger) error {
	clnt, err := g.Client()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)

	mergeRequests, resp, err := clnt.MergeRequests.ListProjectMergeRequests(lmr.ProjectID, &gitlab.ListProjectMergeRequestsOptions{
		State:   gitlab.String(lmr.State),
		Labels:  gitlab.Labels(lmr.Labels),
		OrderBy: gitlab.String("created_at"),
		Sort:    gitlab.String("asc"),
		View:    gitlab.String("simple"),
	}, gitlab.WithContext(ctx))

	cancel()

	if err != nil {
		return fmt.Errorf("failed to list merge requests: %w", err)
	}

	_ = resp.Body.Close()

	/*
		for _, m := range mergeRequests {
			assignee := "none"
			if m.Assignee != nil {
				assignee = m.Assignee.Username
			}

			fmt.Fprintf(w, "%d\t%q\t%s\t%s\t%s\n", m.IID, m.Title, assignee, m.State, m.WebURL)
		}
	*/

	out, err := output.New(g.Format, mergeRequests, []string{"IID", "Title", "Assignee", "State", "WebURL"})
	if err != nil {
		return err
	}

	return out.Print(os.Stdout)
}
