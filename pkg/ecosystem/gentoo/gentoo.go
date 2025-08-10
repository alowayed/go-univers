// Package gentoo provides functionality for working with Gentoo package versions.
package gentoo

const (
	Name = "gentoo"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
