// Package semver provides functionality for working with standard Semantic Versioning 2.0.0.
package semver

const (
	Name = "semver"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
