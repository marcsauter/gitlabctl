package gitlabctl

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/marcsauter/gitlabctl/internal/cmd"
	"github.com/xanzy/go-gitlab"
	"go.uber.org/zap"
)

type createIssueCmd struct {
	ProjectID   int       `help:"Project ID."`
	Title       string    `help:"Issue title" required:""`
	Description string    `help:"Issue description" default:""`
	Labels      []string  `help:"Issue lables"`
	DueDate     time.Time `help:"Issue due date"`
}

func (ci createIssueCmd) Run(app *kong.Context, g *cmd.Globals, l *zap.SugaredLogger) error {
	clnt, err := g.Client()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)

	issueCreateOpts := &gitlab.CreateIssueOptions{
		Title:       gitlab.String(ci.Title),
		Description: gitlab.String(ci.Description),
		Labels:      gitlab.Labels(ci.Labels),
	}

	if !ci.DueDate.IsZero() {
		due := gitlab.ISOTime(ci.DueDate)
		issueCreateOpts.DueDate = &due
	}

	issue, resp, err := clnt.Issues.CreateIssue(ci.ProjectID, issueCreateOpts, gitlab.WithContext(ctx))

	cancel()

	if err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}

	_ = resp.Body.Close()

	fmt.Fprintln(os.Stdout, issue.WebURL)

	return nil
}
