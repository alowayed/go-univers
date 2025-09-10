package apache

import "testing"

func TestEcosystem_Name(t *testing.T) {
	e := &Ecosystem{}
	want := "apache"
	if got := e.Name(); got != want {
		t.Errorf("Ecosystem.Name() = %v, want %v", got, want)
	}
}
