package semver

import "testing"

func TestEcosystem_Name(t *testing.T) {
	e := &Ecosystem{}
	if got := e.Name(); got != "semver" {
		t.Errorf("Name() = %v, want %v", got, "semver")
	}
}
