// Package conan provides functionality for working with Conan C/C++ package manager versions.
package conan

const (
	Name = "conan"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
