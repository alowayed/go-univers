package gentoo

import "testing"

func TestEcosystem_Name(t *testing.T) {
	e := &Ecosystem{}
	expected := "gentoo"
	if got := e.Name(); got != expected {
		t.Errorf("Ecosystem.Name() = %v, want %v", got, expected)
	}
}
