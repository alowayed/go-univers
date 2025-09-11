// Package golang provides functionality for working with Go module versions following semantic versioning with Go-specific extensions.
package golang

const (
	Name = "golang"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
