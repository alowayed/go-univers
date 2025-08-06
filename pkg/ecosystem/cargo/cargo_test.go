package cargo

import (
	"testing"
)

func TestEcosystem_Name(t *testing.T) {
	e := &Ecosystem{}
	if got := e.Name(); got != "cargo" {
		t.Errorf("Name() = %v, want cargo", got)
	}
}

