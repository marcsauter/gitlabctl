package gitlabctl

type mergeRequestCmd struct {
	Create createMergeRequestCmd `cmd:"" help:"Create a merge request."`
	List   listMergeRequestCmd   `cmd:"" help:"List merge requests"`
}
