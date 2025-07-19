// Package pypi provides functionality for working with PyPI package versions following PEP 440.
package pypi

const (
	Name = "pypi"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
