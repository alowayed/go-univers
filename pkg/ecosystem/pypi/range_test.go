package pypi

import "testing"

func TestEcosystem_NewVersionRange(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantErr      bool
		wantOriginal string
	}{
		{
			name:         "exact version",
			input:        "1.2.3",
			wantErr:      false,
			wantOriginal: "1.2.3",
		},
		{
			name:         "greater than",
			input:        ">1.2.3",
			wantErr:      false,
			wantOriginal: ">1.2.3",
		},
		{
			name:         "compatible release",
			input:        "~=1.2.3",
			wantErr:      false,
			wantOriginal: "~=1.2.3",
		},
		{
			name:         "wildcard constraint",
			input:        "==1.2.*",
			wantErr:      false,
			wantOriginal: "==1.2.*",
		},
		{
			name:         "multiple constraints",
			input:        ">=1.0.0, <2.0.0",
			wantErr:      false,
			wantOriginal: ">=1.0.0, <2.0.0",
		},
		{
			name:         "arbitrary equality",
			input:        "===1.2.3",
			wantErr:      false,
			wantOriginal: "===1.2.3",
		},
		{
			name:    "empty specifier",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ecosystem{}
			got, err := e.NewVersionRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersionRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.wantOriginal {
				t.Errorf("NewVersionRange().String() = %v, want %v", got.String(), tt.wantOriginal)
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
			rangeStr: "==1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "exact no match",
			rangeStr: "==1.2.3",
			version:  "1.2.4",
			want:     false,
		},
		{
			name:     "greater than match",
			rangeStr: ">1.2.3",
			version:  "1.2.4",
			want:     true,
		},
		{
			name:     "greater than no match",
			rangeStr: ">1.2.3",
			version:  "1.2.3",
			want:     false,
		},
		{
			name:     "greater than or equal match",
			rangeStr: ">=1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "less than match",
			rangeStr: "<1.2.3",
			version:  "1.2.2",
			want:     true,
		},
		{
			name:     "less than or equal match",
			rangeStr: "<=1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "not equal match",
			rangeStr: "!=1.2.3",
			version:  "1.2.4",
			want:     true,
		},
		{
			name:     "not equal no match",
			rangeStr: "!=1.2.3",
			version:  "1.2.3",
			want:     false,
		},
		{
			name:     "compatible release match",
			rangeStr: "~=1.2.3",
			version:  "1.2.5",
			want:     true,
		},
		{
			name:     "compatible release no match",
			rangeStr: "~=1.2.3",
			version:  "1.3.0",
			want:     false,
		},
		{
			name:     "wildcard constraint match",
			rangeStr: "==1.2.*",
			version:  "1.2.5",
			want:     true,
		},
		{
			name:     "wildcard constraint no match",
			rangeStr: "==1.2.*",
			version:  "1.3.0",
			want:     false,
		},
		{
			name:     "multiple constraints match",
			rangeStr: ">=1.0.0, <2.0.0",
			version:  "1.5.0",
			want:     true,
		},
		{
			name:     "multiple constraints no match",
			rangeStr: ">=1.0.0, <2.0.0",
			version:  "2.0.0",
			want:     false,
		},
		{
			name:     "arbitrary equality match",
			rangeStr: "===1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "arbitrary equality no match",
			rangeStr: "===1.2.3",
			version:  "1.2.3.post1",
			want:     false,
		},
		{
			name:     "prerelease handling",
			rangeStr: ">=1.0.0a1",
			version:  "1.0.0b1",
			want:     true,
		},
		{
			name:     "epoch handling",
			rangeStr: ">=1!1.0.0",
			version:  "1!1.0.1",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := mustNewVersionRange(t, tt.rangeStr)
			v := mustNewVersion(t, tt.version)

			got := vr.Contains(v)
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
