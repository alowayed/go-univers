package gomod

import (
	"testing"
)

func TestEcosystem_NewVersionRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "exact version",
			input: "v1.2.3",
		},
		{
			name:  "greater than",
			input: ">v1.2.3",
		},
		{
			name:  "greater than or equal",
			input: ">=v1.2.3",
		},
		{
			name:  "less than",
			input: "<v2.0.0",
		},
		{
			name:  "less than or equal",
			input: "<=v1.9.9",
		},
		{
			name:  "not equal",
			input: "!=v1.3.0",
		},
		{
			name:  "multiple constraints",
			input: ">=v1.2.3 <v2.0.0",
		},
		{
			name:    "empty range",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ecosystem{}
			got, err := e.NewVersionRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersionRange(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.String() != tt.input {
					t.Errorf("VersionRange{%q}.String() = %q, want %q", tt.input, got.String(), tt.input)
				}
			}
		})
	}
}

func TestVersionRange_Contains(t *testing.T) {
	tests := []struct {
		name     string
		rangeStr string
		version  string
		want     bool
	}{
		{
			name:     "exact match",
			rangeStr: "v1.2.3",
			version:  "v1.2.3",
			want:     true,
		},
		{
			name:     "exact no match",
			rangeStr: "v1.2.3",
			version:  "v1.2.4",
			want:     false,
		},
		{
			name:     "greater than match",
			rangeStr: ">v1.2.3",
			version:  "v1.2.4",
			want:     true,
		},
		{
			name:     "greater than no match",
			rangeStr: ">v1.2.3",
			version:  "v1.2.3",
			want:     false,
		},
		{
			name:     "greater than or equal match",
			rangeStr: ">=v1.2.3",
			version:  "v1.2.3",
			want:     true,
		},
		{
			name:     "less than match",
			rangeStr: "<v2.0.0",
			version:  "v1.9.9",
			want:     true,
		},
		{
			name:     "less than no match",
			rangeStr: "<v2.0.0",
			version:  "v2.0.0",
			want:     false,
		},
		{
			name:     "not equal match",
			rangeStr: "!=v1.3.0",
			version:  "v1.2.9",
			want:     true,
		},
		{
			name:     "not equal no match",
			rangeStr: "!=v1.3.0",
			version:  "v1.3.0",
			want:     false,
		},
		{
			name:     "multiple constraints match",
			rangeStr: ">=v1.2.3 <v2.0.0",
			version:  "v1.5.0",
			want:     true,
		},
		{
			name:     "multiple constraints no match lower",
			rangeStr: ">=v1.2.3 <v2.0.0",
			version:  "v1.2.2",
			want:     false,
		},
		{
			name:     "multiple constraints no match upper",
			rangeStr: ">=v1.2.3 <v2.0.0",
			version:  "v2.0.0",
			want:     false,
		},
		{
			name:     "pseudo-version in range",
			rangeStr: ">=v1.2.3 <v2.0.0",
			version:  "v1.5.0-0.20170915032832-14c0d48ead0c",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := mustNewVersionRange(t, tt.rangeStr)
			version := mustNewVersion(t, tt.version)

			got := vr.Contains(version)
			if got != tt.want {
				t.Errorf("VersionRange{%q}.Contains(%q) = %v, want %v", tt.rangeStr, tt.version, got, tt.want)
			}
		})
	}
}

// mustNewVersionRange is a helper function to create a new VersionRange.
func mustNewVersionRange(t *testing.T, s string) *VersionRange {
	t.Helper()
	e := &Ecosystem{}
	vr, err := e.NewVersionRange(s)
	if err != nil {
		t.Fatalf("Failed to create version range %q: %v", s, err)
	}
	return vr
}
