// Package rpm provides functionality for working with Red Hat Package Manager (RPM) versions.
package rpm

const (
	Name = "rpm"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
