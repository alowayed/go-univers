// Package maven provides functionality for working with Maven package versions.
package maven

const (
	Name = "maven"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
