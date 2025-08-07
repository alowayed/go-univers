// Package npm provides functionality for working with NPM package versions.
package npm

const (
	Name = "npm"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
