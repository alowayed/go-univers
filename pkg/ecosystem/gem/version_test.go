package gem

import (
	"testing"
)

func TestEcosystem_NewVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid versions
		{"basic semantic version", "1.0.0", false},
		{"standard version", "1.2.3", false},
		{"version with build number", "1.2.3.4", false},
		{"minimal version", "0.0.1", false},
		{"large numbers", "10.20.30", false},

		// Prerelease versions
		{"alpha prerelease", "1.0.0-alpha", false},
		{"beta prerelease", "1.0.0-beta", false},
		{"release candidate", "1.0.0-rc1", false},
		{"pre suffix", "1.2.3.pre", false},
		{"rc with numbers", "2.0.0.rc1", false},
		{"complex prerelease", "1.0.0-beta.1", false},

		// Build metadata
		{"build metadata", "1.0.0+build.1", false},
		{"prerelease with build", "1.0.0-alpha+build", false},

		// Edge cases
		{"v prefix", "v1.0.0", false},
		{"single digit", "1", false},
		{"two segments", "1.0", false},

		// Whitespace test cases
		{"leading space", " 1.0.0", false},
		{"trailing space", "1.0.0 ", false},
		{"leading and trailing spaces", " 1.0.0 ", false},
		{"leading tab", "\t1.0.0", false},
		{"trailing tab", "1.0.0\t", false},
		{"leading and trailing tabs", "\t1.0.0\t", false},
		{"leading newline", "\n1.0.0", false},
		{"trailing newline", "1.0.0\n", false},
		{"leading and trailing newlines", "\n1.0.0\n", false},
		{"mixed whitespace", " \t\n1.0.0\n\t ", false},

		// Invalid versions
		{"empty string", "", true},
		{"non-numeric", "invalid", true},
		{"double dots", "1.2.3..4", true},
		{"only spaces", "   ", true},
		{"only tabs", "\t\t\t", true},
		{"only newlines", "\n\n\n", true},
		{"mixed whitespace only", " \t\n ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ecosystem{}
			got, err := e.NewVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ecosystem.NewVersion(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.String() != tt.input {
					t.Errorf("Version.String() = %q, want %q", got.String(), tt.input)
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
		{"equal versions", "1.0.0", "1.0.0", 0},
		{"patch difference", "1.0.0", "1.0.1", -1},
		{"patch difference reverse", "1.0.1", "1.0.0", 1},
		{"minor difference", "1.0.0", "1.1.0", -1},
		{"minor difference reverse", "1.1.0", "1.0.0", 1},
		{"major difference", "1.0.0", "2.0.0", -1},
		{"major difference reverse", "2.0.0", "1.0.0", 1},

		// Prerelease comparisons
		{"release vs prerelease", "1.0.0", "1.0.0-alpha", 1},
		{"prerelease vs release", "1.0.0-alpha", "1.0.0", -1},
		{"alpha vs beta", "1.0.0-alpha", "1.0.0-beta", -1},
		{"beta vs alpha", "1.0.0-beta", "1.0.0-alpha", 1},
		{"beta vs rc", "1.0.0-beta", "1.0.0-rc", -1},
		{"rc1 vs rc2", "1.0.0-rc1", "1.0.0-rc2", -1},
		{"pre vs release", "1.2.3.pre", "1.2.3", -1},
		{"rc vs release", "2.0.0-rc1", "2.0.0", -1},

		// Complex versions
		{"build number difference", "1.2.3.4", "1.2.3.5", -1},

		// Edge cases
		{"implicit zero", "1.0", "1.0.0", 0},
		{"single vs triple", "1", "1.0.0", 0},
		{"release vs next major prerelease", "0.9.9", "1.0.0-alpha", -1},

		// Whitespace handling in comparison
		{"whitespace vs normal", " 1.0.0 ", "1.0.0", 0},
		{"tab vs normal", "\t1.0.0\t", "1.0.0", 0},
		{"newline vs normal", "\n1.0.0\n", "1.0.0", 0},
		{"mixed whitespace vs normal", " \t\n1.0.0\n\t ", "1.0.0", 0},
		{"both with whitespace", " 1.0.0 ", "\t1.0.0\n", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := mustNewVersion(t, tt.v1)
			v2 := mustNewVersion(t, tt.v2)

			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Version.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

// mustNewVersion is a helper function to create a new Version.
func mustNewVersion(t *testing.T, version string) *Version {
	t.Helper()
	e := &Ecosystem{}
	v, err := e.NewVersion(version)
	if err != nil {
		t.Fatalf("Failed to create version %s: %v", version, err)
	}
	return v
}
