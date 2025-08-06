package gem

import (
	"testing"
)

func TestEcosystem_NewVersionRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid ranges
		{"exact version", "1.2.3", false},
		{"explicit equals", "= 1.2.3", false},
		{"greater than or equal", ">= 1.2.3", false},
		{"greater than", "> 1.2.3", false},
		{"less than or equal", "<= 1.2.3", false},
		{"less than", "< 1.2.3", false},
		{"not equal", "!= 1.2.3", false},
		{"pessimistic constraint", "~> 1.2.3", false},
		{"pessimistic major.minor", "~> 1.2", false},

		// Multiple constraints
		{"pessimistic with minimum", "~> 1.2.3, >= 1.2.5", false},
		{"range constraint", ">= 1.0.0, < 2.0.0", false},
		{"pessimistic with exclusion", "~> 2.0, != 2.1.0", false},

		// Whitespace test cases
		{"range with leading whitespace", " >= 1.2.3 ", false},
		{"range with internal whitespace", ">= 1.0.0 , < 2.0.0", false},
		{"only whitespace", "   ", true},

		// Invalid ranges
		{"empty string", "", true},
		{"pessimistic without version", "~>", true},
		{"operator without version", ">=", true},
		{"only commas", " , ", true},
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
		// Exact matches
		{"exact match", "1.2.3", "1.2.3", true},
		{"explicit equals match", "= 1.2.3", "1.2.3", true},
		{"exact no match", "1.2.3", "1.2.4", false},

		// Comparison operators
		{"gte equal", ">= 1.2.3", "1.2.3", true},
		{"gte greater", ">= 1.2.3", "1.2.4", true},
		{"gte less", ">= 1.2.3", "1.2.2", false},
		{"gt greater", "> 1.2.3", "1.2.4", true},
		{"gt equal", "> 1.2.3", "1.2.3", false},
		{"lte equal", "<= 1.2.3", "1.2.3", true},
		{"lte less", "<= 1.2.3", "1.2.2", true},
		{"lte greater", "<= 1.2.3", "1.2.4", false},
		{"lt less", "< 1.2.3", "1.2.2", true},
		{"lt equal", "< 1.2.3", "1.2.3", false},
		{"ne different", "!= 1.2.3", "1.2.4", true},
		{"ne same", "!= 1.2.3", "1.2.3", false},

		// Pessimistic constraints (~>)
		{"pessimistic exact", "~> 1.2.3", "1.2.3", true},
		{"pessimistic patch ok", "~> 1.2.3", "1.2.4", true},
		{"pessimistic patch high", "~> 1.2.3", "1.2.9", true},
		{"pessimistic minor bump", "~> 1.2.3", "1.3.0", false},
		{"pessimistic major bump", "~> 1.2.3", "2.0.0", false},
		{"pessimistic below", "~> 1.2.3", "1.2.2", false},

		{"pessimistic minor exact", "~> 1.2", "1.2.0", true},
		{"pessimistic minor patch ok", "~> 1.2", "1.2.9", true},
		{"pessimistic minor high", "~> 1.2", "1.9.9", true},
		{"pessimistic minor major bump", "~> 1.2", "2.0.0", false},
		{"pessimistic minor below", "~> 1.2", "1.1.9", false},

		{"pessimistic major exact", "~> 1", "1.0.0", true},
		{"pessimistic major ok", "~> 1", "1.9.9", true},
		{"pessimistic major bump", "~> 1", "2.0.0", false},

		// Pessimistic constraint edge cases
		{"pessimistic zero major exact", "~> 0.1.0", "0.1.0", true},
		{"pessimistic zero major patch ok", "~> 0.1.0", "0.1.5", true},
		{"pessimistic zero major minor bump", "~> 0.1.0", "0.2.0", false},

		{"pessimistic zero minor exact", "~> 0.0.1", "0.0.1", true},
		{"pessimistic zero minor patch ok", "~> 0.0.1", "0.0.2", true},
		{"pessimistic zero minor bump", "~> 0.0.1", "0.1.0", false},

		// Prerelease handling with pessimistic constraints
		{"pessimistic prerelease below", "~> 1.0.0", "1.0.0-alpha", false},
		{"pessimistic prerelease exact", "~> 1.0.0-alpha", "1.0.0-alpha", true},
		{"pessimistic prerelease beta", "~> 1.0.0-alpha", "1.0.0-beta", true},
		{"pessimistic prerelease release", "~> 1.0.0-alpha", "1.0.0", true},
		{"pessimistic prerelease patch bump", "~> 1.0.0-alpha", "1.0.1", false},

		// Prerelease handling
		{"prerelease gte", ">= 1.0.0-alpha", "1.0.0-alpha", true},
		{"prerelease order", ">= 1.0.0-alpha", "1.0.0-beta", true},
		{"prerelease to release", ">= 1.0.0-alpha", "1.0.0", true},

		// Multiple constraints (AND logic)
		{"range contains", ">= 1.2.0, < 2.0.0", "1.5.0", true},
		{"range below", ">= 1.2.0, < 2.0.0", "1.1.9", false},
		{"range above", ">= 1.2.0, < 2.0.0", "2.0.0", false},
		{"pessimistic with min", "~> 1.2.3, >= 1.2.5", "1.2.5", true},
		{"pessimistic with min below", "~> 1.2.3, >= 1.2.5", "1.2.4", false},
		{"pessimistic with exclusion ok", "~> 2.0, != 2.1.0", "2.0.5", true},
		{"pessimistic with exclusion hit", "~> 2.0, != 2.1.0", "2.1.0", false},
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
