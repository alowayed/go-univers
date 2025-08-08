package cran

import (
	"testing"
)

func TestEcosystem_Name(t *testing.T) {
	ecosystem := &Ecosystem{}
	want := "cran"

	got := ecosystem.Name()
	if got != want {
		t.Errorf("Ecosystem.Name() = %q, want %q", got, want)
	}
}
