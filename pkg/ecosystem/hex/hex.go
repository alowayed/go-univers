package hex

const Name = "hex"

type Ecosystem struct{}

func (e *Ecosystem) Name() string {
	return Name
}
