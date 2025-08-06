package cargo

import (
	"testing"
)

func TestEcosystem_NewVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Version
		wantErr bool
	}{
		// Basic SemVer versions
		{
			name:  "basic semantic version",
			input: "1.2.3",
			want: &Version{
				major:      1,
				minor:      2,
				patch:      3,
				prerelease: "",
				build:      "",
				original:   "1.2.3",
			},
		},
		{
			name:  "zero versions",
			input: "0.0.0",
			want: &Version{
				major:      0,
				minor:      0,
				patch:      0,
				prerelease: "",
				build:      "",
				original:   "0.0.0",
			},
		},
		{
			name:  "large version numbers",
			input: "999.888.777",
			want: &Version{
				major:      999,
				minor:      888,
				patch:      777,
				prerelease: "",
				build:      "",
				original:   "999.888.777",
			},
		},

		// Prerelease versions
		{
			name:  "alpha prerelease",
			input: "1.0.0-alpha",
			want: &Version{
				major:      1,
				minor:      0,
				patch:      0,
				prerelease: "alpha",
				build:      "",
				original:   "1.0.0-alpha",
			},
		},
		{
			name:  "beta prerelease with number",
			input: "2.0.0-beta.1",
			want: &Version{
				major:      2,
				minor:      0,
				patch:      0,
				prerelease: "beta.1",
				build:      "",
				original:   "2.0.0-beta.1",
			},
		},
		{
			name:  "rc prerelease",
			input: "1.2.3-rc.1",
			want: &Version{
				major:      1,
				minor:      2,
				patch:      3,
				prerelease: "rc.1",
				build:      "",
				original:   "1.2.3-rc.1",
			},
		},
		{
			name:  "complex prerelease",
			input: "1.0.0-alpha.beta.1",
			want: &Version{
				major:      1,
				minor:      0,
				patch:      0,
				prerelease: "alpha.beta.1",
				build:      "",
				original:   "1.0.0-alpha.beta.1",
			},
		},

		// Build metadata
		{
			name:  "version with build metadata",
			input: "1.0.0+build.1",
			want: &Version{
				major:      1,
				minor:      0,
				patch:      0,
				prerelease: "",
				build:      "build.1",
				original:   "1.0.0+build.1",
			},
		},
		{
			name:  "version with prerelease and build",
			input: "1.0.0-alpha+build.1",
			want: &Version{
				major:      1,
				minor:      0,
				patch:      0,
				prerelease: "alpha",
				build:      "build.1",
				original:   "1.0.0-alpha+build.1",
			},
		},
		{
			name:  "complex version with all components",
			input: "2.1.0-rc.1+build.20230101",
			want: &Version{
				major:      2,
				minor:      1,
				patch:      0,
				prerelease: "rc.1",
				build:      "build.20230101",
				original:   "2.1.0-rc.1+build.20230101",
			},
		},

		// Error cases
		{name: "empty string", input: "", wantErr: true},
		{name: "whitespace only", input: "   ", wantErr: true},
		{name: "invalid format", input: "1.2", wantErr: true},
		{name: "invalid format with letters", input: "1.2.3a", wantErr: true},
		{name: "non-numeric major", input: "a.2.3", wantErr: true},
		{name: "non-numeric minor", input: "1.b.3", wantErr: true},
		{name: "non-numeric patch", input: "1.2.c", wantErr: true},
		{name: "negative version", input: "-1.2.3", wantErr: true},
		{name: "leading v not allowed", input: "v1.2.3", wantErr: true},
		{name: "four components", input: "1.2.3.4", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ecosystem{}
			got, err := e.NewVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !equalVersions(got, tt.want) {
					t.Errorf("NewVersion() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		// Basic comparisons
		{name: "equal versions", v1: "1.2.3", v2: "1.2.3", want: 0},
		{name: "major version difference", v1: "2.0.0", v2: "1.9.9", want: 1},
		{name: "minor version difference", v1: "1.3.0", v2: "1.2.9", want: 1},
		{name: "patch version difference", v1: "1.2.4", v2: "1.2.3", want: 1},

		// Prerelease comparisons (SemVer 2.0 rules)
		{name: "release vs prerelease", v1: "1.0.0", v2: "1.0.0-alpha", want: 1},
		{name: "prerelease vs release", v1: "1.0.0-alpha", v2: "1.0.0", want: -1},
		{name: "alpha vs beta", v1: "1.0.0-alpha", v2: "1.0.0-beta", want: -1},
		{name: "alpha.1 vs alpha.2", v1: "1.0.0-alpha.1", v2: "1.0.0-alpha.2", want: -1},
		{name: "alpha vs alpha.1", v1: "1.0.0-alpha", v2: "1.0.0-alpha.1", want: -1},
		{name: "alpha.1 vs alpha.beta", v1: "1.0.0-alpha.1", v2: "1.0.0-alpha.beta", want: -1},
		{name: "alpha.beta vs beta", v1: "1.0.0-alpha.beta", v2: "1.0.0-beta", want: -1},
		{name: "beta vs beta.2", v1: "1.0.0-beta", v2: "1.0.0-beta.2", want: -1},
		{name: "beta.2 vs beta.11", v1: "1.0.0-beta.2", v2: "1.0.0-beta.11", want: -1},
		{name: "beta.11 vs rc.1", v1: "1.0.0-beta.11", v2: "1.0.0-rc.1", want: -1},
		{name: "rc.1 vs release", v1: "1.0.0-rc.1", v2: "1.0.0", want: -1},

		// Build metadata should be ignored in comparison
		{name: "build metadata ignored - equal", v1: "1.0.0+build.1", v2: "1.0.0+build.2", want: 0},
		{name: "build metadata ignored - with prerelease", v1: "1.0.0-alpha+build.1", v2: "1.0.0-alpha+build.2", want: 0},

		// Mixed prerelease and build
		{name: "prerelease with build vs release", v1: "1.0.0-alpha+build", v2: "1.0.0", want: -1},
		{name: "release vs prerelease with build", v1: "1.0.0", v2: "1.0.0-alpha+build", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ecosystem{}
			v1, err := e.NewVersion(tt.v1)
			if err != nil {
				t.Fatalf("NewVersion(%q) error: %v", tt.v1, err)
			}

			v2, err := e.NewVersion(tt.v2)
			if err != nil {
				t.Fatalf("NewVersion(%q) error: %v", tt.v2, err)
			}

			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Compare(%q, %q) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

func TestComparePrereleaseIdentifiers(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want int
	}{
		// Numeric vs non-numeric
		{name: "numeric vs non-numeric", a: "1", b: "alpha", want: -1},
		{name: "non-numeric vs numeric", a: "alpha", b: "1", want: 1},

		// Numeric comparisons
		{name: "numeric equal", a: "1", b: "1", want: 0},
		{name: "numeric 1 vs 2", a: "1", b: "2", want: -1},
		{name: "numeric 2 vs 1", a: "2", b: "1", want: 1},
		{name: "numeric 2 vs 11", a: "2", b: "11", want: -1},

		// String comparisons
		{name: "string equal", a: "alpha", b: "alpha", want: 0},
		{name: "alpha vs beta", a: "alpha", b: "beta", want: -1},
		{name: "beta vs alpha", a: "beta", b: "alpha", want: 1},

		// Different lengths
		{name: "shorter vs longer", a: "alpha", b: "alpha.1", want: -1},
		{name: "longer vs shorter", a: "alpha.1", b: "alpha", want: 1},

		// Complex cases
		{name: "alpha.1 vs alpha.beta", a: "alpha.1", b: "alpha.beta", want: -1},
		{name: "alpha.beta vs alpha.1", a: "alpha.beta", b: "alpha.1", want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := comparePrereleaseIdentifiers(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("comparePrereleaseIdentifiers(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// Helper functions
func equalVersions(a, b *Version) bool {
	return a.major == b.major &&
		a.minor == b.minor &&
		a.patch == b.patch &&
		a.prerelease == b.prerelease &&
		a.build == b.build &&
		a.original == b.original
}
