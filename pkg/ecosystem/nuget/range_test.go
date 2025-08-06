package nuget

import (
	"testing"
)

func TestEcosystem_NewVersionRange(t *testing.T) {
	e := &Ecosystem{}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid ranges
		{"exact version", "[1.2.3]", false},
		{"inclusive range", "[1.0.0,2.0.0]", false},
		{"exclusive range", "(1.0.0,2.0.0)", false},
		{"mixed range 1", "[1.0.0,2.0.0)", false},
		{"mixed range 2", "(1.0.0,2.0.0]", false},
		{"unbounded minimum", "[1.0.0,)", false},
		{"unbounded maximum", "(,2.0.0]", false},
		{"minimum version", "1.0.0", false},
		{"comma separated", ">=1.0.0,<2.0.0", false},

		// Error cases
		{"empty range", "", true},
		{"whitespace only", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := e.NewVersionRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersionRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVersionRange_Contains(t *testing.T) {
	e := &Ecosystem{}

	tests := []struct {
		name     string
		rangeStr string
		version  string
		want     bool
	}{
		// Exact version matches
		{"exact match", "[1.2.3]", "1.2.3", true},
		{"exact no match", "[1.2.3]", "1.2.4", false},

		// Inclusive ranges
		{"inclusive range start", "[1.0.0,2.0.0]", "1.0.0", true},
		{"inclusive range end", "[1.0.0,2.0.0]", "2.0.0", true},
		{"inclusive range middle", "[1.0.0,2.0.0]", "1.5.0", true},
		{"inclusive range outside low", "[1.0.0,2.0.0]", "0.9.0", false},
		{"inclusive range outside high", "[1.0.0,2.0.0]", "2.1.0", false},

		// Exclusive ranges
		{"exclusive range start", "(1.0.0,2.0.0)", "1.0.0", false},
		{"exclusive range end", "(1.0.0,2.0.0)", "2.0.0", false},
		{"exclusive range middle", "(1.0.0,2.0.0)", "1.5.0", true},
		{"exclusive range outside low", "(1.0.0,2.0.0)", "0.9.0", false},
		{"exclusive range outside high", "(1.0.0,2.0.0)", "2.1.0", false},

		// Mixed ranges
		{"mixed range [,)", "[1.0.0,2.0.0)", "1.0.0", true},
		{"mixed range [,)", "[1.0.0,2.0.0)", "2.0.0", false},
		{"mixed range (,]", "(1.0.0,2.0.0]", "1.0.0", false},
		{"mixed range (,]", "(1.0.0,2.0.0]", "2.0.0", true},

		// Unbounded ranges
		{"unbounded min", "[1.0.0,)", "1.0.0", true},
		{"unbounded min", "[1.0.0,)", "2.0.0", true},
		{"unbounded min", "[1.0.0,)", "0.9.0", false},
		{"unbounded max", "(,2.0.0]", "2.0.0", true},
		{"unbounded max", "(,2.0.0]", "1.0.0", true},
		{"unbounded max", "(,2.0.0]", "2.1.0", false},

		// Minimum version (default behavior)
		{"minimum version", "1.0.0", "1.0.0", true},
		{"minimum version", "1.0.0", "1.5.0", true},
		{"minimum version", "1.0.0", "0.9.0", false},

		// Prerelease versions
		{"prerelease in range", "[1.0.0,2.0.0]", "1.5.0-alpha", true},
		{"prerelease exact", "[1.0.0-alpha]", "1.0.0-alpha", true},
		{"prerelease vs release", "[1.0.0-alpha]", "1.0.0", false},

		// Revision versions
		{"revision in range", "[1.0.0,2.0.0]", "1.5.0.1", true},
		{"revision exact", "[1.0.0.1]", "1.0.0.1", true},
		{"revision vs no revision", "[1.0.0.0]", "1.0.0", true},

		// Comma-separated constraints
		{"comma constraints", ">=1.0.0,<2.0.0", "1.5.0", true},
		{"comma constraints", ">=1.0.0,<2.0.0", "0.9.0", false},
		{"comma constraints", ">=1.0.0,<2.0.0", "2.0.0", false},
		{"comma constraints with exclusion", ">=1.0.0,<2.0.0,!=1.5.0", "1.4.0", true},
		{"comma constraints with exclusion", ">=1.0.0,<2.0.0,!=1.5.0", "1.5.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr, err := e.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("Failed to parse range %s: %v", tt.rangeStr, err)
			}

			v, err := e.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("Failed to parse version %s: %v", tt.version, err)
			}

			got := vr.Contains(v)
			if got != tt.want {
				t.Errorf("Contains(%s, %s) = %v, want %v", tt.rangeStr, tt.version, got, tt.want)
			}
		})
	}
}

func TestVersionRange_String(t *testing.T) {
	e := &Ecosystem{}

	tests := []string{
		"[1.2.3]",
		"[1.0.0,2.0.0]",
		"(1.0.0,2.0.0)",
		"[1.0.0,)",
		"(,2.0.0]",
		"1.0.0",
		">=1.0.0,<2.0.0",
	}

	for _, rangeStr := range tests {
		t.Run(rangeStr, func(t *testing.T) {
			vr, err := e.NewVersionRange(rangeStr)
			if err != nil {
				t.Fatalf("Failed to parse range %s: %v", rangeStr, err)
			}

			if got := vr.String(); got != rangeStr {
				t.Errorf("String() = %s, want %s", got, rangeStr)
			}
		})
	}
}

func TestVersionRange_ParseEdgeCases(t *testing.T) {
	e := &Ecosystem{}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Malformed ranges
		{"unclosed bracket", "[1.0.0,2.0.0", true},
		{"unclosed paren", "(1.0.0,2.0.0", true},
		{"mismatched brackets", "[1.0.0,2.0.0)", false}, // This should work (mixed range)
		{"empty brackets", "[]", true},
		{"empty parens", "()", true},
		{"single comma", ",", true},                // Should fail - no valid constraints
		{"multiple commas", "1.0.0,,2.0.0", false}, // Should handle gracefully
		{"whitespace in brackets", "[ 1.0.0 , 2.0.0 ]", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := e.NewVersionRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersionRange(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}


