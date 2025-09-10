package github

const (
	Name = "github"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
