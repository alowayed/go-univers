package alpine

import (
	"testing"
)

func TestEcosystem_Name(t *testing.T) {
	ecosystem := &Ecosystem{}
	want := "alpine"
	
	got := ecosystem.Name()
	if got != want {
		t.Errorf("Ecosystem.Name() = %q, want %q", got, want)
	}
}
