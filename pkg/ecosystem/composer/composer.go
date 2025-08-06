// Package composer provides functionality for working with Composer (PHP) package versions.
package composer

const (
	Name = "composer"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}