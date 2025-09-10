package apache

const (
	Name = "apache"
)

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
