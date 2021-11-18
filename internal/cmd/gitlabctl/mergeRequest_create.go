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

type createMergeRequestCmd struct {
	ProjectID   int      `help:"Project ID."`
	Title       string   `help:"Merge request title." required:"" `
	Description string   `help:"Merge request description." default:"" `
	Source      string   `help:"Source branch" required:"" `
	Target      string   `help:"Target branch" required:"" `
	Labels      []string `help:"Merge request lables"`
}

func (cmr createMergeRequestCmd) Run(app *kong.Context, g *cmd.Globals, l *zap.SugaredLogger) error {
	clnt, err := g.Client()
	if err != nil {
		return err
	}

	diff, err := cmr.hasDiffs(clnt, g.Timeout, cmr.ProjectID)
	if err != nil {
		return err
	}

	if !diff {
		return fmt.Errorf("branches are no different")
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)

	mr, resp, err := clnt.MergeRequests.CreateMergeRequest(cmr.ProjectID, &gitlab.CreateMergeRequestOptions{
		Title:        gitlab.String(cmr.Title),
		Description:  gitlab.String(cmr.Description),
		SourceBranch: gitlab.String(cmr.Source),
		TargetBranch: gitlab.String(cmr.Target),
		Labels:       gitlab.Labels(cmr.Labels),
	}, gitlab.WithContext(ctx))

	cancel()

	if err != nil {
		return fmt.Errorf("failed to create merge request: %w", err)
	}

	_ = resp.Body.Close()

	fmt.Fprintln(os.Stdout, mr.WebURL)

	return nil
}

func (cmr createMergeRequestCmd) hasDiffs(clnt *gitlab.Client, timeout time.Duration, pid int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	compare, resp, err := clnt.Repositories.Compare(pid, &gitlab.CompareOptions{
		From: gitlab.String(cmr.Target),
		To:   gitlab.String(cmr.Source),
	}, gitlab.WithContext(ctx))

	cancel()

	if err != nil {
		return false, fmt.Errorf("failed to compare %s with %s: %w", cmr.Source, cmr.Target, err)
	}

	_ = resp.Body.Close()

	return len(compare.Diffs) != 0, nil
}
