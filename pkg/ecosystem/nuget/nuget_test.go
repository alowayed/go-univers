package nuget

import "testing"

func TestEcosystem_Name(t *testing.T) {
	e := &Ecosystem{}
	want := "nuget"
	if got := e.Name(); got != want {
		t.Errorf("Ecosystem.Name() = %v, want %v", got, want)
	}
}