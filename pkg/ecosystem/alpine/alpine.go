// Package alpine provides functionality for working with Alpine Linux package versions.
package alpine

const (
	Name = "alpine"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
