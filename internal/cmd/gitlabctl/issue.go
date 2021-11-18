package gitlabctl

type issueCmd struct {
	Create createIssueCmd `cmd:"" help:"Create an issue"`
	List   listIssueCmd   `cmd:"" help:"List issues"`
}
