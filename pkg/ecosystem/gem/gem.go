// Package gem provides functionality for working with Ruby Gem package versions.
package gem

const (
	Name = "gem"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
