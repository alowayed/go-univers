package rpm

import (
	"testing"
)

func TestEcosystem_Name(t *testing.T) {
	e := &Ecosystem{}
	if got := e.Name(); got != Name {
		t.Errorf("Ecosystem.Name() = %v, want %v", got, Name)
	}
}
