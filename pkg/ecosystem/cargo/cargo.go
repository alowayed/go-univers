// Package cargo provides functionality for working with Cargo (Rust) package versions.
package cargo

const (
	Name = "cargo"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}

