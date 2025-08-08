// Package cran provides functionality for working with CRAN (R package) versions.
package cran

const (
	Name = "cran"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
