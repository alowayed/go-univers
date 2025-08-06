package nuget

import (
	"testing"
)

func TestEcosystem_NewVersion(t *testing.T) {
	e := &Ecosystem{}

	tests := []struct {
		name    string
		input   string
		want    *Version
		wantErr bool
	}{
		{
			name:  "basic version",
			input: "1.2.3",
			want: &Version{
				major: 1, minor: 2, patch: 3, revision: 0,
				prerelease: "", build: "", original: "1.2.3",
			},
		},
		{
			name:  "version with v prefix",
			input: "v1.2.3",
			want: &Version{
				major: 1, minor: 2, patch: 3, revision: 0,
				prerelease: "", build: "", original: "v1.2.3",
			},
		},
		{
			name:  "version with revision",
			input: "1.2.3.4",
			want: &Version{
				major: 1, minor: 2, patch: 3, revision: 4,
				prerelease: "", build: "", original: "1.2.3.4",
			},
		},
		{
			name:  "version with prerelease",
			input: "1.2.3-alpha",
			want: &Version{
				major: 1, minor: 2, patch: 3, revision: 0,
				prerelease: "alpha", build: "", original: "1.2.3-alpha",
			},
		},
		{
			name:  "version with prerelease and build",
			input: "1.2.3-alpha+build.1",
			want: &Version{
				major: 1, minor: 2, patch: 3, revision: 0,
				prerelease: "alpha", build: "build.1", original: "1.2.3-alpha+build.1",
			},
		},
		{
			name:  "version with revision, prerelease, and build",
			input: "1.2.3.4-beta.1+build.20230101",
			want: &Version{
				major: 1, minor: 2, patch: 3, revision: 4,
				prerelease: "beta.1", build: "build.20230101", original: "1.2.3.4-beta.1+build.20230101",
			},
		},
		{
			name:  "major only",
			input: "1",
			want: &Version{
				major: 1, minor: 0, patch: 0, revision: 0,
				prerelease: "", build: "", original: "1",
			},
		},
		{
			name:  "major and minor only",
			input: "1.2",
			want: &Version{
				major: 1, minor: 2, patch: 0, revision: 0,
				prerelease: "", build: "", original: "1.2",
			},
		},
		{
			name:  "version with complex prerelease",
			input: "1.0.0-alpha.beta.1",
			want: &Version{
				major: 1, minor: 0, patch: 0, revision: 0,
				prerelease: "alpha.beta.1", build: "", original: "1.0.0-alpha.beta.1",
			},
		},
		{
			name:  "version with build metadata only",
			input: "1.0.0+build.1",
			want: &Version{
				major: 1, minor: 0, patch: 0, revision: 0,
				prerelease: "", build: "build.1", original: "1.0.0+build.1",
			},
		},
		{
			name:  "whitespace trimmed",
			input: "  1.2.3  ",
			want: &Version{
				major: 1, minor: 2, patch: 3, revision: 0,
				prerelease: "", build: "", original: "1.2.3",
			},
		},
		// Error cases
		{
			name:    "empty version",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			input:   "1.2.3.4.5",
			wantErr: true,
		},
		{
			name:    "non-numeric major",
			input:   "a.2.3",
			wantErr: true,
		},
		{
			name:    "non-numeric minor",
			input:   "1.b.3",
			wantErr: true,
		},
		{
			name:    "non-numeric patch",
			input:   "1.2.c",
			wantErr: true,
		},
		{
			name:    "non-numeric revision",
			input:   "1.2.3.d",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := e.NewVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.major != tt.want.major || got.minor != tt.want.minor ||
					got.patch != tt.want.patch || got.revision != tt.want.revision ||
					got.prerelease != tt.want.prerelease || got.build != tt.want.build ||
					got.original != tt.want.original {
					t.Errorf("NewVersion() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	e := &Ecosystem{}

	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		// Basic version comparisons
		{"equal versions", "1.2.3", "1.2.3", 0},
		{"major version difference", "2.0.0", "1.9.9", 1},
		{"minor version difference", "1.3.0", "1.2.9", 1},
		{"patch version difference", "1.2.4", "1.2.3", 1},
		{"revision version difference", "1.2.3.4", "1.2.3.3", 1},

		// Prerelease comparisons
		{"prerelease vs release", "1.0.0-alpha", "1.0.0", -1},
		{"release vs prerelease", "1.0.0", "1.0.0-alpha", 1},
		{"prerelease vs prerelease", "1.0.0-alpha", "1.0.0-beta", -1},
		{"numeric prerelease", "1.0.0-alpha.1", "1.0.0-alpha.2", -1},
		{"mixed prerelease", "1.0.0-alpha.1", "1.0.0-alpha.beta", -1},
		{"complex prerelease", "1.0.0-alpha.beta.1", "1.0.0-alpha.beta.2", -1},

		// Build metadata is ignored in comparison
		{"build metadata ignored", "1.0.0+build.1", "1.0.0+build.2", 0},
		{"prerelease with build", "1.0.0-alpha+build.1", "1.0.0-alpha+build.2", 0},

		// Edge cases
		{"zero versions", "0.0.0", "0.0.1", -1},
		{"missing components", "1.0", "1.0.0", 0},
		{"revision vs no revision", "1.2.3.0", "1.2.3", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := e.NewVersion(tt.v1)
			if err != nil {
				t.Fatalf("Failed to parse v1 %s: %v", tt.v1, err)
			}
			v2, err := e.NewVersion(tt.v2)
			if err != nil {
				t.Fatalf("Failed to parse v2 %s: %v", tt.v2, err)
			}

			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Compare(%s, %s) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}

			// Test symmetry
			expected := -tt.want
			if reverse := v2.Compare(v1); reverse != expected {
				t.Errorf("Compare(%s, %s) = %d, want %d (reverse)", tt.v2, tt.v1, reverse, expected)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	e := &Ecosystem{}

	tests := []struct {
		name  string
		input string
	}{
		{"basic version", "1.2.3"},
		{"version with v prefix", "v1.2.3"},
		{"version with prerelease", "1.2.3-alpha"},
		{"version with build", "1.0.0+build.1"},
		{"complex version", "1.2.3.4-beta.1+build.20230101"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := e.NewVersion(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse version %s: %v", tt.input, err)
			}

			if got := v.String(); got != tt.input {
				t.Errorf("String() = %s, want %s", got, tt.input)
			}
		})
	}
}
