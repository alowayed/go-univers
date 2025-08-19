package conan

import "testing"

func TestEcosystem_Name(t *testing.T) {
	e := &Ecosystem{}
	if got := e.Name(); got != "conan" {
		t.Errorf("Name() = %v, want %v", got, "conan")
	}
}
