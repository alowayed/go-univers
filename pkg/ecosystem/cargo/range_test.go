package cargo

import (
	"testing"
)

func TestEcosystem_NewVersionRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Basic ranges
		{name: "exact version", input: "1.2.3", wantErr: false},
		{name: "caret range", input: "^1.2.3", wantErr: false},
		{name: "tilde range", input: "~1.2.3", wantErr: false},
		{name: "greater than", input: ">1.2.3", wantErr: false},
		{name: "greater equal", input: ">=1.2.3", wantErr: false},
		{name: "less than", input: "<2.0.0", wantErr: false},
		{name: "less equal", input: "<=2.0.0", wantErr: false},
		{name: "wildcard", input: "1.2.*", wantErr: false},

		// Multiple constraints
		{name: "multiple constraints", input: ">=1.2.3, <2.0.0", wantErr: false},
		{name: "complex multiple", input: ">=1.0.0, <2.0.0, !=1.5.0", wantErr: false},
		
		// Partial version constraints
		{name: "tilde partial ~1.2", input: "~1.2", wantErr: false},
		{name: "tilde partial ~1", input: "~1", wantErr: false},
		{name: "caret partial ^1.2", input: "^1.2", wantErr: false},
		{name: "caret partial ^1", input: "^1", wantErr: false},
		
		// Wildcard constraints
		{name: "wildcard 1.2.*", input: "1.2.*", wantErr: false},
		{name: "wildcard 1.*", input: "1.*", wantErr: false},
		{name: "wildcard *", input: "*", wantErr: false},

		// Error cases
		{name: "empty string", input: "", wantErr: true},
		{name: "whitespace only", input: "   ", wantErr: true},
		{name: "invalid caret", input: "^invalid", wantErr: true},
		{name: "invalid tilde", input: "~invalid", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ecosystem{}
			_, err := e.NewVersionRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersionRange() error = %v, wantErr %v", err, tt.wantErr)
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
		{name: "exact match", rangeStr: "1.2.3", version: "1.2.3", want: true},
		{name: "exact no match", rangeStr: "1.2.3", version: "1.2.4", want: false},

		// Comparison operators
		{name: "greater than - true", rangeStr: ">1.2.3", version: "1.2.4", want: true},
		{name: "greater than - false", rangeStr: ">1.2.3", version: "1.2.3", want: false},
		{name: "greater equal - true", rangeStr: ">=1.2.3", version: "1.2.3", want: true},
		{name: "greater equal - true higher", rangeStr: ">=1.2.3", version: "1.2.4", want: true},
		{name: "greater equal - false", rangeStr: ">=1.2.3", version: "1.2.2", want: false},

		{name: "less than - true", rangeStr: "<2.0.0", version: "1.9.9", want: true},
		{name: "less than - false", rangeStr: "<2.0.0", version: "2.0.0", want: false},
		{name: "less equal - true", rangeStr: "<=2.0.0", version: "2.0.0", want: true},
		{name: "less equal - false", rangeStr: "<=2.0.0", version: "2.0.1", want: false},

		// Caret constraints (^)
		{name: "caret - same version", rangeStr: "^1.2.3", version: "1.2.3", want: true},
		{name: "caret - patch increment", rangeStr: "^1.2.3", version: "1.2.4", want: true},
		{name: "caret - minor increment", rangeStr: "^1.2.3", version: "1.3.0", want: true},
		{name: "caret - major increment", rangeStr: "^1.2.3", version: "2.0.0", want: false},
		{name: "caret - lower version", rangeStr: "^1.2.3", version: "1.2.2", want: false},

		// Caret with major version 0
		{name: "caret 0.x.y - minor change", rangeStr: "^0.2.3", version: "0.3.0", want: false},
		{name: "caret 0.x.y - patch change", rangeStr: "^0.2.3", version: "0.2.4", want: true},
		{name: "caret 0.0.x - no change", rangeStr: "^0.0.3", version: "0.0.4", want: false},
		{name: "caret 0.0.x - exact", rangeStr: "^0.0.3", version: "0.0.3", want: true},

		// Tilde constraints (~)
		{name: "tilde - same version", rangeStr: "~1.2.3", version: "1.2.3", want: true},
		{name: "tilde - patch increment", rangeStr: "~1.2.3", version: "1.2.4", want: true},
		{name: "tilde - minor increment", rangeStr: "~1.2.3", version: "1.3.0", want: false},
		{name: "tilde - major increment", rangeStr: "~1.2.3", version: "2.0.0", want: false},
		{name: "tilde - lower version", rangeStr: "~1.2.3", version: "1.2.2", want: false},
		
		// Partial tilde constraints
		{name: "tilde partial ~1.2", rangeStr: "~1.2", version: "1.2.5", want: true},
		{name: "tilde partial ~1.2 minor change", rangeStr: "~1.2", version: "1.3.0", want: false},
		{name: "tilde partial ~1", rangeStr: "~1", version: "1.5.0", want: true},
		{name: "tilde partial ~1 major change", rangeStr: "~1", version: "2.0.0", want: false},
		
		// Partial caret constraints  
		{name: "caret partial ^1.2", rangeStr: "^1.2", version: "1.5.0", want: true},
		{name: "caret partial ^1", rangeStr: "^1", version: "1.9.0", want: true},

		// Wildcard constraints (now converted to standard constraints)
		{name: "wildcard 1.2.* (becomes ~1.2.0)", rangeStr: "1.2.*", version: "1.2.3", want: true},
		{name: "wildcard 1.2.* patch different", rangeStr: "1.2.*", version: "1.2.9", want: true},
		{name: "wildcard 1.2.* minor different", rangeStr: "1.2.*", version: "1.3.0", want: false},
		{name: "wildcard 1.* (becomes ^1.0.0)", rangeStr: "1.*", version: "1.2.0", want: true},
		{name: "wildcard 1.* major different", rangeStr: "1.*", version: "2.0.0", want: false},
		{name: "wildcard * (becomes >=0.0.0)", rangeStr: "*", version: "1.0.0", want: true},

		// Multiple constraints (AND logic)
		{name: "multiple - both satisfied", rangeStr: ">=1.2.3, <2.0.0", version: "1.5.0", want: true},
		{name: "multiple - first not satisfied", rangeStr: ">=1.2.3, <2.0.0", version: "1.2.2", want: false},
		{name: "multiple - second not satisfied", rangeStr: ">=1.2.3, <2.0.0", version: "2.0.0", want: false},
		
		// Prerelease handling
		{name: "prerelease in range", rangeStr: ">=1.0.0-alpha", version: "1.0.0-beta", want: true},
		{name: "prerelease vs release", rangeStr: ">=1.0.0", version: "1.0.0-alpha", want: false},
		{name: "release vs prerelease range", rangeStr: ">=1.0.0-alpha", version: "1.0.0", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ecosystem{}
			vr, err := e.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("NewVersionRange(%q) error: %v", tt.rangeStr, err)
			}

			version, err := e.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("NewVersion(%q) error: %v", tt.version, err)
			}

			got := vr.Contains(version)
			if got != tt.want {
				t.Errorf("Contains(%q, %q) = %t, want %t", tt.rangeStr, tt.version, got, tt.want)
			}
		})
	}
}


