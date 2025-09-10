package alpm

// Ecosystem represents the ALPM (Arch Linux Package Manager) versioning ecosystem
type Ecosystem struct{}

// Name returns the name of this ecosystem
const Name = "alpm"

func (e *Ecosystem) Name() string {
	return Name
}
