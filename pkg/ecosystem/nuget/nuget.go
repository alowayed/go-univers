// Package nuget provides functionality for working with NuGet (.NET) package versions.
package nuget

const (
	Name = "nuget"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}

