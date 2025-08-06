package maven

import (
	"testing"
)

func TestEcosystem_NewVersionRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty range", "", true},
		{"exact version", "[1.0.0]", false},
		{"lower bound inclusive", "[1.0.0,)", false},
		{"lower bound exclusive", "(1.0.0,)", false},
		{"upper bound inclusive", "(,1.0.0]", false},
		{"upper bound exclusive", "(,1.0.0)", false},
		{"inclusive range", "[1.0.0,2.0.0]", false},
		{"exclusive range", "(1.0.0,2.0.0)", false},
		{"mixed brackets", "[1.0.0,2.0.0)", false},
		{"mixed brackets reverse", "(1.0.0,2.0.0]", false},
		{"simple version", "1.0.0", false},
		{"version with qualifier", "[1.0.0-alpha,1.0.0]", false},

		// Whitespace test cases
		{"range with leading whitespace", " [1.0.0,2.0.0] ", false},
		{"range with internal whitespace", "[1.0.0, 2.0.0]", false},
		{"only whitespace", "   ", true},

		// Error cases
		{"invalid version in range", "[invalid,2.0.0]", true},
		{"malformed bracket", "[1.0.0,2.0.0", true},
		{"empty version in exact", "[]", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ecosystem{}
			got, err := e.NewVersionRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ecosystem.NewVersionRange(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.String() != tt.input {
					t.Errorf("VersionRange.String() = %q, want %q", got.String(), tt.input)
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
		// Exact version tests
		{"exact match", "[1.0.0]", "1.0.0", true},
		{"exact no match", "[1.0.0]", "1.0.1", false},
		{"exact with qualifier", "[1.0.0-alpha]", "1.0.0-alpha", true},
		{"exact qualifier no match", "[1.0.0-alpha]", "1.0.0-beta", false},

		// Lower bound inclusive tests
		{"lower bound inclusive match", "[1.0.0,)", "1.0.0", true},
		{"lower bound inclusive higher", "[1.0.0,)", "1.0.1", true},
		{"lower bound inclusive lower", "[1.0.0,)", "0.9.9", false},
		{"lower bound inclusive major", "[1.0.0,)", "2.0.0", true},

		// Lower bound exclusive tests
		{"lower bound exclusive match", "(1.0.0,)", "1.0.0", false},
		{"lower bound exclusive higher", "(1.0.0,)", "1.0.1", true},
		{"lower bound exclusive lower", "(1.0.0,)", "0.9.9", false},

		// Upper bound inclusive tests
		{"upper bound inclusive match", "(,1.0.0]", "1.0.0", true},
		{"upper bound inclusive lower", "(,1.0.0]", "0.9.9", true},
		{"upper bound inclusive higher", "(,1.0.0]", "1.0.1", false},

		// Upper bound exclusive tests
		{"upper bound exclusive match", "(,1.0.0)", "1.0.0", false},
		{"upper bound exclusive lower", "(,1.0.0)", "0.9.9", true},
		{"upper bound exclusive higher", "(,1.0.0)", "1.0.1", false},

		// Range tests (inclusive)
		{"inclusive range match lower", "[1.0.0,2.0.0]", "1.0.0", true},
		{"inclusive range match upper", "[1.0.0,2.0.0]", "2.0.0", true},
		{"inclusive range match middle", "[1.0.0,2.0.0]", "1.5.0", true},
		{"inclusive range below", "[1.0.0,2.0.0]", "0.9.9", false},
		{"inclusive range above", "[1.0.0,2.0.0]", "2.0.1", false},

		// Range tests (exclusive)
		{"exclusive range match middle", "(1.0.0,2.0.0)", "1.5.0", true},
		{"exclusive range match bounds", "(1.0.0,2.0.0)", "1.0.0", false},
		{"exclusive range upper bound", "(1.0.0,2.0.0)", "2.0.0", false},
		{"exclusive range below", "(1.0.0,2.0.0)", "0.9.9", false},
		{"exclusive range above", "(1.0.0,2.0.0)", "2.0.1", false},

		// Mixed bracket tests
		{"mixed lower inclusive", "[1.0.0,2.0.0)", "1.0.0", true},
		{"mixed upper exclusive", "[1.0.0,2.0.0)", "2.0.0", false},
		{"mixed lower exclusive", "(1.0.0,2.0.0]", "1.0.0", false},
		{"mixed upper inclusive", "(1.0.0,2.0.0]", "2.0.0", true},

		// Simple version (treated as exact)
		{"simple version match", "1.0.0", "1.0.0", true},
		{"simple version no match", "1.0.0", "1.0.1", false},

		// Qualifier precedence tests
		{"alpha in range", "[1.0.0-alpha,1.0.0]", "1.0.0-beta", true},
		{"snapshot before release", "[1.0.0-snapshot,1.0.0]", "1.0.0", true},
		{"qualifier outside range", "[1.0.0,2.0.0]", "0.9.0-alpha", false},

		// Edge cases with normalization
		{"normalized versions", "[1.0,1.0.0]", "1.0.0", true},
		{"ga equivalent", "[1.0.0-ga]", "1.0.0", true},
		{"final equivalent", "[1.0.0]", "1.0.0-final", true},
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
