package project

type Cmd struct {
	Project projectCmd `cmd:"" help:""`
}

type projectCmd struct {
	List listProjectCmd `cmd:"" help:"List projects."`
}
