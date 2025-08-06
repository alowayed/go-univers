// Package gomod provides functionality for working with Go module versions following semantic versioning with Go-specific extensions.
package gomod

const (
	Name = "go"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}

