package cran

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
		// Basic versions
		{
			name:  "basic two component version",
			input: "1.2",
			want: &Version{
				components: []int{1, 2},
				original:   "1.2",
			},
		},
		{
			name:  "basic three component version",
			input: "1.2.3",
			want: &Version{
				components: []int{1, 2, 3},
				original:   "1.2.3",
			},
		},
		{
			name:  "four component version",
			input: "1.2.3.4",
			want: &Version{
				components: []int{1, 2, 3, 4},
				original:   "1.2.3.4",
			},
		},
		{
			name:  "version with dashes",
			input: "1-2-3",
			want: &Version{
				components: []int{1, 2, 3},
				original:   "1-2-3",
			},
		},
		{
			name:  "mixed separators",
			input: "1.2-3",
			want: &Version{
				components: []int{1, 2, 3},
				original:   "1.2-3",
			},
		},
		{
			name:  "version with zeros",
			input: "0.0.1",
			want: &Version{
				components: []int{0, 0, 1},
				original:   "0.0.1",
			},
		},
		{
			name:  "large version numbers",
			input: "10.20.30",
			want: &Version{
				components: []int{10, 20, 30},
				original:   "10.20.30",
			},
		},
		{
			name:  "version with whitespace",
			input: "  1.2.3  ",
			want: &Version{
				components: []int{1, 2, 3},
				original:   "  1.2.3  ",
			},
		},
		// Error cases
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: true,
		},
		{
			name:    "single component",
			input:   "1",
			wantErr: true,
		},
		{
			name:    "non-numeric component",
			input:   "1.2.abc",
			wantErr: true,
		},
		{
			name:    "empty component",
			input:   "1..3",
			wantErr: true,
		},
		{
			name:    "negative component",
			input:   "1.-2.3",
			wantErr: true,
		},
		{
			name:    "version starting with separator",
			input:   ".1.2",
			wantErr: true,
		},
		{
			name:    "version ending with separator",
			input:   "1.2.",
			wantErr: true,
		},
		{
			name:    "only separators",
			input:   "...",
			wantErr: true,
		},
		{
			name:    "mixed invalid characters",
			input:   "1.2.3a",
			wantErr: true,
		},
		{
			name:    "letters only",
			input:   "abc.def",
			wantErr: true,
		},
	}

	ecosystem := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ecosystem.NewVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ecosystem.NewVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if got.String() != tt.want.String() {
				t.Errorf("Version.String() = %v, want %v", got.String(), tt.want.String())
			}

			if len(got.components) != len(tt.want.components) {
				t.Errorf("Version components length = %v, want %v", len(got.components), len(tt.want.components))
				return
			}

			for i := range got.components {
				if got.components[i] != tt.want.components[i] {
					t.Errorf("Version component[%d] = %v, want %v", i, got.components[i], tt.want.components[i])
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
		// Equal versions
		{
			name: "equal basic versions",
			v1:   "1.2",
			v2:   "1.2",
			want: 0,
		},
		{
			name: "equal three component versions",
			v1:   "1.2.3",
			v2:   "1.2.3",
			want: 0,
		},
		{
			name: "equal with dashes",
			v1:   "1-2-3",
			v2:   "1.2.3",
			want: 0,
		},
		// First version greater
		{
			name: "major version greater",
			v1:   "2.0",
			v2:   "1.0",
			want: 1,
		},
		{
			name: "minor version greater",
			v1:   "1.2",
			v2:   "1.1",
			want: 1,
		},
		{
			name: "patch version greater",
			v1:   "1.2.3",
			v2:   "1.2.2",
			want: 1,
		},
		{
			name: "more components greater",
			v1:   "1.2.3",
			v2:   "1.2",
			want: 1,
		},
		{
			name: "more components greater with same prefix",
			v1:   "1.2.3.4",
			v2:   "1.2.3",
			want: 1,
		},
		// First version less
		{
			name: "major version less",
			v1:   "1.0",
			v2:   "2.0",
			want: -1,
		},
		{
			name: "minor version less",
			v1:   "1.1",
			v2:   "1.2",
			want: -1,
		},
		{
			name: "patch version less",
			v1:   "1.2.2",
			v2:   "1.2.3",
			want: -1,
		},
		{
			name: "fewer components less",
			v1:   "1.2",
			v2:   "1.2.3",
			want: -1,
		},
		{
			name: "fewer components less with same prefix",
			v1:   "1.2.3",
			v2:   "1.2.3.4",
			want: -1,
		},
		// Complex cases
		{
			name: "different major overrides minor",
			v1:   "2.0",
			v2:   "1.9",
			want: 1,
		},
		{
			name: "different minor overrides patch",
			v1:   "1.2.0",
			v2:   "1.1.9",
			want: 1,
		},
		{
			name: "zero components",
			v1:   "0.1",
			v2:   "0.0",
			want: 1,
		},
		{
			name: "large numbers",
			v1:   "100.200",
			v2:   "100.199",
			want: 1,
		},
	}

	ecosystem := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := ecosystem.NewVersion(tt.v1)
			if err != nil {
				t.Fatalf("Failed to parse v1 %s: %v", tt.v1, err)
			}

			v2, err := ecosystem.NewVersion(tt.v2)
			if err != nil {
				t.Fatalf("Failed to parse v2 %s: %v", tt.v2, err)
			}

			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Version.Compare() = %v, want %v (v1=%s, v2=%s)", got, tt.want, tt.v1, tt.v2)
			}

			// Test symmetry - reverse comparison should give opposite result
			reverse := v2.Compare(v1)
			expected := -tt.want
			if reverse != expected {
				t.Errorf("Reverse comparison failed: v2.Compare(v1) = %v, want %v", reverse, expected)
			}
		})
	}
}
