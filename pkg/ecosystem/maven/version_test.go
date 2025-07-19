package maven

import (
	"testing"
)

func TestEcosystem_NewVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty version", "", true},
		{"simple version", "1.0.0", false},
		{"version with qualifier", "1.0.0-alpha", false},
		{"snapshot version", "1.0.0-SNAPSHOT", false},

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
		// Basic numeric comparison
		{"equal versions", "1.0.0", "1.0.0", 0},
		{"major version", "1.0.0", "2.0.0", -1},
		{"minor version", "1.1.0", "1.0.0", 1},
		{"patch version", "1.0.1", "1.0.0", 1},

		// Qualifier precedence
		{"alpha vs beta", "1.0.0-alpha", "1.0.0-beta", -1},
		{"beta vs milestone", "1.0.0-beta", "1.0.0-milestone", -1},
		{"milestone vs rc", "1.0.0-milestone", "1.0.0-rc", -1},
		{"rc vs snapshot", "1.0.0-rc", "1.0.0-snapshot", -1},
		{"snapshot vs release", "1.0.0-snapshot", "1.0.0", -1},
		{"release vs sp", "1.0.0", "1.0.0-sp", -1},

		// Qualifier shortcuts
		{"a vs alpha", "1.0.0-a", "1.0.0-alpha", 0},
		{"b vs beta", "1.0.0-b", "1.0.0-beta", 0},
		{"m vs milestone", "1.0.0-m", "1.0.0-milestone", 0},
		{"cr vs rc", "1.0.0-cr", "1.0.0-rc", 0},

		// Normalization
		{"release equivalents", "1.0.0", "1.0.0-ga", 0},
		{"final equivalent", "1.0.0-final", "1.0.0", 0},
		{"trailing zeros", "1.0", "1.0.0", 0},

		// Complex qualifiers
		{"qualified numbers", "1.0.0-alpha-1", "1.0.0-alpha-2", -1},
		{"beta with number", "1.0.0-beta-1", "1.0.0-beta-10", -1},

		// Mixed types
		{"number vs qualifier", "1.0.1", "1.0.0-alpha", 1},
		{"different lengths", "1.0", "1.0.0.1", -1},

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