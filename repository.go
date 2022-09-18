package pm

type Repository struct {
	Author  string
	Name    string
	Version string
	Deps    []*Dependencies
}
