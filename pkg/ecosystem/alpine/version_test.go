package alpine

import (
	"testing"
)

func TestEcosystem_NewVersion(t *testing.T) {
	ecosystem := &Ecosystem{}
	
	tests := []struct {
		name    string
		input   string
		want    *Version
		wantErr bool
	}{
		// Basic versions
		{
			name:  "basic semantic version",
			input: "1.2.3",
			want: &Version{
				numeric:  []int{1, 2, 3},
				letter:   "",
				suffixes: nil,
				hash:     "",
				build:    0,
				original: "1.2.3",
			},
		},
		{
			name:  "two component version",
			input: "0.1.0",
			want: &Version{
				numeric:  []int{0, 1, 0},
				letter:   "",
				suffixes: nil,
				hash:     "",
				build:    0,
				original: "0.1.0",
			},
		},
		{
			name:  "single component version",
			input: "1",
			want: &Version{
				numeric:  []int{1},
				letter:   "",
				suffixes: nil,
				hash:     "",
				build:    0,
				original: "1",
			},
		},
		
		// Versions with letters
		{
			name:  "version with letter suffix",
			input: "1.2.3a",
			want: &Version{
				numeric:  []int{1, 2, 3},
				letter:   "a",
				suffixes: nil,
				hash:     "",
				build:    0,
				original: "1.2.3a",
			},
		},
		{
			name:  "version with letter b",
			input: "2.3.0b",
			want: &Version{
				numeric:  []int{2, 3, 0},
				letter:   "b",
				suffixes: nil,
				hash:     "",
				build:    0,
				original: "2.3.0b",
			},
		},
		
		// Versions with suffixes
		{
			name:  "version with alpha suffix",
			input: "1.2.3_alpha",
			want: &Version{
				numeric:  []int{1, 2, 3},
				letter:   "",
				suffixes: []suffix{{name: "alpha", number: 0}},
				hash:     "",
				build:    0,
				original: "1.2.3_alpha",
			},
		},
		{
			name:  "version with numbered alpha suffix",
			input: "1.3_alpha2",
			want: &Version{
				numeric:  []int{1, 3},
				letter:   "",
				suffixes: []suffix{{name: "alpha", number: 2}},
				hash:     "",
				build:    0,
				original: "1.3_alpha2",
			},
		},
		{
			name:  "version with multiple suffixes",
			input: "0.1.0_alpha_pre2",
			want: &Version{
				numeric:  []int{0, 1, 0},
				letter:   "",
				suffixes: []suffix{{name: "alpha", number: 0}, {name: "pre", number: 2}},
				hash:     "",
				build:    0,
				original: "0.1.0_alpha_pre2",
			},
		},
		
		// Versions with build components
		{
			name:  "version with build component",
			input: "1.0.4-r3",
			want: &Version{
				numeric:  []int{1, 0, 4},
				letter:   "",
				suffixes: nil,
				hash:     "",
				build:    3,
				original: "1.0.4-r3",
			},
		},
		{
			name:  "version with build r1",
			input: "2.2.39-r1",
			want: &Version{
				numeric:  []int{2, 2, 39},
				letter:   "",
				suffixes: nil,
				hash:     "",
				build:    1,
				original: "2.2.39-r1",
			},
		},
		{
			name:  "date-based version with build",
			input: "20050718-r2",
			want: &Version{
				numeric:  []int{20050718},
				letter:   "",
				suffixes: nil,
				hash:     "",
				build:    2,
				original: "20050718-r2",
			},
		},
		
		// Complex versions
		{
			name:  "version with letter and build",
			input: "2.3.0b-r1",
			want: &Version{
				numeric:  []int{2, 3, 0},
				letter:   "b",
				suffixes: nil,
				hash:     "",
				build:    1,
				original: "2.3.0b-r1",
			},
		},
		{
			name:  "version with letter, suffix, and build",
			input: "1.2.3a_beta2-r5",
			want: &Version{
				numeric:  []int{1, 2, 3},
				letter:   "a",
				suffixes: []suffix{{name: "beta", number: 2}},
				hash:     "",
				build:    5,
				original: "1.2.3a_beta2-r5",
			},
		},
		
		// Versions with hash
		{
			name:  "version with hash",
			input: "1.2.3~abc123",
			want: &Version{
				numeric:  []int{1, 2, 3},
				letter:   "",
				suffixes: nil,
				hash:     "abc123",
				build:    0,
				original: "1.2.3~abc123",
			},
		},
		{
			name:  "version with hash and build",
			input: "1.2.3~abc123-r1",
			want: &Version{
				numeric:  []int{1, 2, 3},
				letter:   "",
				suffixes: nil,
				hash:     "abc123",
				build:    1,
				original: "1.2.3~abc123-r1",
			},
		},
		
		// Error cases
		{name: "empty string", input: "", wantErr: true},
		{name: "whitespace only", input: "   ", wantErr: true},
		{name: "invalid format", input: "abc", wantErr: true},
		{name: "incomplete build", input: "1.2.3-", wantErr: true},
		{name: "unknown suffix", input: "1.2.3_unknown", wantErr: true},
		{name: "invalid build format", input: "1.2.3-abc", wantErr: true},
		{name: "incomplete hash", input: "1.2.3~", wantErr: true},
		{name: "invalid characters", input: "1.2.3!", wantErr: true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ecosystem.NewVersion(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewVersion() expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("NewVersion() unexpected error: %v", err)
				return
			}
			
			if got.String() != tt.want.original {
				t.Errorf("NewVersion().String() = %q, want %q", got.String(), tt.want.original)
			}
			
			// Check components
			if !equalIntSlices(got.numeric, tt.want.numeric) {
				t.Errorf("NewVersion().numeric = %v, want %v", got.numeric, tt.want.numeric)
			}
			
			if got.letter != tt.want.letter {
				t.Errorf("NewVersion().letter = %q, want %q", got.letter, tt.want.letter)
			}
			
			if !equalSuffixSlices(got.suffixes, tt.want.suffixes) {
				t.Errorf("NewVersion().suffixes = %v, want %v", got.suffixes, tt.want.suffixes)
			}
			
			if got.hash != tt.want.hash {
				t.Errorf("NewVersion().hash = %q, want %q", got.hash, tt.want.hash)
			}
			
			if got.build != tt.want.build {
				t.Errorf("NewVersion().build = %d, want %d", got.build, tt.want.build)
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	ecosystem := &Ecosystem{}
	
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		// Basic numeric comparisons
		{name: "equal versions", v1: "1.2.3", v2: "1.2.3", want: 0},
		{name: "first version smaller", v1: "1.2.3", v2: "1.2.4", want: -1},
		{name: "first version larger", v1: "1.2.4", v2: "1.2.3", want: 1},
		{name: "missing patch component", v1: "1.2", v2: "1.2.0", want: 0},
		{name: "missing patch component reverse", v1: "1.2.0", v2: "1.2", want: 0},
		
		// Multi-component comparisons
		{name: "major version difference", v1: "1.0.4", v2: "2.2.39", want: -1},
		{name: "major version difference reverse", v1: "2.2.39", v2: "1.0.4", want: 1},
		{name: "minor version difference", v1: "4.8", v2: "3.10.18", want: 1},
		{name: "minor version difference reverse", v1: "3.10.18", v2: "4.8", want: -1},
		
		// Letter comparisons
		{name: "release vs letter", v1: "1.1", v2: "1.1a", want: -1},
		{name: "letter vs release", v1: "1.1a", v2: "1.1", want: 1},
		{name: "letter comparison", v1: "1.1a", v2: "1.1b", want: -1},
		{name: "letter comparison reverse", v1: "1.1b", v2: "1.1a", want: 1},
		
		// Suffix comparisons (based on test data from alpine_test.txt)
		{name: "alpha versions", v1: "0.1.0_alpha", v2: "0.1.3_alpha", want: -1},
		{name: "alpha versions reverse", v1: "0.1.3_alpha", v2: "0.1.0_alpha", want: 1},
		{name: "release vs alpha", v1: "1.1", v2: "1.1_alpha1", want: 1},
		{name: "alpha vs release", v1: "1.1_alpha1", v2: "1.1", want: -1},
		{name: "numbered alpha", v1: "1.3_alpha", v2: "1.3_alpha2", want: -1},
		{name: "numbered alpha reverse", v1: "1.3_alpha2", v2: "1.3_alpha", want: 1},
		{name: "alpha with pre vs alpha", v1: "0.1.0_alpha_pre2", v2: "0.1.0_alpha", want: -1},
		{name: "alpha vs alpha with pre", v1: "0.1.0_alpha", v2: "0.1.0_alpha_pre2", want: 1},
		
		// Build component comparisons
		{name: "build versions", v1: "1.0.4-r3", v2: "1.0.4-r4", want: -1},
		{name: "build versions reverse", v1: "1.0.4-r4", v2: "1.0.4-r3", want: 1},
		{name: "letter with build", v1: "2.3.0b-r1", v2: "2.3.0b-r2", want: -1},
		{name: "letter with build reverse", v1: "2.3.0b-r2", v2: "2.3.0b-r1", want: 1},
		{name: "different base with build", v1: "2.2.39-r1", v2: "1.0.4-r3", want: 1},
		{name: "different base with build reverse", v1: "1.0.4-r3", v2: "2.2.39-r1", want: -1},
		
		// Date-based versions
		{name: "date-based build", v1: "20050718-r1", v2: "20050718-r2", want: -1},
		{name: "date-based build reverse", v1: "20050718-r2", v2: "20050718-r1", want: 1},
		
		// Mixed suffix types (suffix precedence: alpha < beta < pre < rc < (none) < cvs < svn < git < hg < p)
		{name: "alpha vs beta", v1: "1.0.0_alpha", v2: "1.0.0_beta", want: -1},
		{name: "beta vs alpha", v1: "1.0.0_beta", v2: "1.0.0_alpha", want: 1},
		{name: "beta vs pre", v1: "1.0.0_beta", v2: "1.0.0_pre", want: -1},
		{name: "pre vs beta", v1: "1.0.0_pre", v2: "1.0.0_beta", want: 1},
		{name: "pre vs rc", v1: "1.0.0_pre", v2: "1.0.0_rc", want: -1},
		{name: "rc vs pre", v1: "1.0.0_rc", v2: "1.0.0_pre", want: 1},
		{name: "rc vs release", v1: "1.0.0_rc", v2: "1.0.0", want: -1},
		{name: "release vs rc", v1: "1.0.0", v2: "1.0.0_rc", want: 1},
		{name: "release vs cvs", v1: "1.0.0", v2: "1.0.0_cvs", want: -1},
		{name: "cvs vs release", v1: "1.0.0_cvs", v2: "1.0.0", want: 1},
		
		// Hash comparisons (lexicographic)
		{name: "hash comparison", v1: "1.0.0~abc", v2: "1.0.0~abd", want: -1},
		{name: "hash comparison reverse", v1: "1.0.0~abd", v2: "1.0.0~abc", want: 1},
		{name: "hash equality", v1: "1.0.0~abc", v2: "1.0.0~abc", want: 0},
		{name: "no hash vs hash", v1: "1.0.0", v2: "1.0.0~abc", want: -1},
		{name: "hash vs no hash", v1: "1.0.0~abc", v2: "1.0.0", want: 1},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := ecosystem.NewVersion(tt.v1)
			if err != nil {
				t.Fatalf("NewVersion(%q) error: %v", tt.v1, err)
			}
			
			v2, err := ecosystem.NewVersion(tt.v2)
			if err != nil {
				t.Fatalf("NewVersion(%q) error: %v", tt.v2, err)
			}
			
			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Compare(%q, %q) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}
			
			// Test symmetry
			reverse := v2.Compare(v1)
			if reverse != -tt.want {
				t.Errorf("Compare symmetry failed: Compare(%q, %q) = %d, but Compare(%q, %q) = %d, want %d", 
					tt.v1, tt.v2, got, tt.v2, tt.v1, reverse, -tt.want)
			}
		})
	}
}

// Helper functions
func equalIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalSuffixSlices(a, b []suffix) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].name != b[i].name || a[i].number != b[i].number {
			return false
		}
	}
	return true
}