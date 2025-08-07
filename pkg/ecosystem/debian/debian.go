// Package debian provides functionality for working with Debian package versions.
package debian

const (
	Name = "debian"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
