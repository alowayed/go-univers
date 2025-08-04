package alpine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
			name:  "basic semantic version",
			input: "1.2.3",
			want: &Version{
				numeric:  []numericComponent{{value: 1, originalStr: "1"}, {value: 2, originalStr: "2"}, {value: 3, originalStr: "3"}},
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
				numeric:  []numericComponent{{value: 0, originalStr: "0"}, {value: 1, originalStr: "1"}, {value: 0, originalStr: "0"}},
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
				numeric:  []numericComponent{{value: 1, originalStr: "1"}},
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
				numeric:  []numericComponent{{value: 1, originalStr: "1"}, {value: 2, originalStr: "2"}, {value: 3, originalStr: "3"}},
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
				numeric:  []numericComponent{{value: 2, originalStr: "2"}, {value: 3, originalStr: "3"}, {value: 0, originalStr: "0"}},
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
				numeric:  []numericComponent{{value: 1, originalStr: "1"}, {value: 2, originalStr: "2"}, {value: 3, originalStr: "3"}},
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
				numeric:  []numericComponent{{value: 1, originalStr: "1"}, {value: 3, originalStr: "3"}},
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
				numeric:  []numericComponent{{value: 0, originalStr: "0"}, {value: 1, originalStr: "1"}, {value: 0, originalStr: "0"}},
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
				numeric:  []numericComponent{{value: 1, originalStr: "1"}, {value: 0, originalStr: "0"}, {value: 4, originalStr: "4"}},
				letter:   "",
				suffixes: nil,
				hash:     "",
				build:    3,
				original: "1.0.4-r3",
			},
		},
		{
			name:  "date-based version with build",
			input: "20050718-r2",
			want: &Version{
				numeric:  []numericComponent{{value: 20050718, originalStr: "20050718"}},
				letter:   "",
				suffixes: nil,
				hash:     "",
				build:    2,
				original: "20050718-r2",
			},
		},

		// Versions with hash
		{
			name:  "version with hash",
			input: "1.2.3~abc123",
			want: &Version{
				numeric:  []numericComponent{{value: 1, originalStr: "1"}, {value: 2, originalStr: "2"}, {value: 3, originalStr: "3"}},
				letter:   "",
				suffixes: nil,
				hash:     "abc123",
				build:    0,
				original: "1.2.3~abc123",
			},
		},

		// Leading zero versions
		{
			name:  "version with leading zeros",
			input: "4.09.14",
			want: &Version{
				numeric:  []numericComponent{{value: 4, originalStr: "4"}, {value: 9, originalStr: "09"}, {value: 14, originalStr: "14"}},
				letter:   "",
				suffixes: nil,
				hash:     "",
				build:    0,
				original: "4.09.14",
			},
		},

		// Unknown suffixes (should now work)
		{
			name:  "version with unknown suffix",
			input: "23_foo",
			want: &Version{
				numeric:  []numericComponent{{value: 23, originalStr: "23"}},
				letter:   "",
				suffixes: []suffix{{name: "foo", number: 0}},
				hash:     "",
				build:    0,
				original: "23_foo",
			},
		},

		// Invalid formats (should now work with string fallback)
		{
			name:  "invalid format fallback",
			input: "1.0bc",
			want: &Version{
				numeric:  nil,
				letter:   "",
				suffixes: nil,
				hash:     "",
				build:    0,
				original: "1.0bc",
			},
		},

		// Error cases
		{name: "empty string", input: "", wantErr: true},
		{name: "whitespace only", input: "   ", wantErr: true},
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
		// Basic numeric comparisons
		{name: "equal versions", v1: "1.2.3", v2: "1.2.3", want: 0},
		{name: "first version smaller", v1: "1.2.3", v2: "1.2.4", want: -1},
		{name: "first version larger", v1: "1.2.4", v2: "1.2.3", want: 1},

		// Letter comparisons
		{name: "release vs letter", v1: "1.1", v2: "1.1a", want: -1},
		{name: "letter vs release", v1: "1.1a", v2: "1.1", want: 1},

		// Suffix comparisons
		{name: "alpha vs release", v1: "1.1_alpha1", v2: "1.1", want: -1},
		{name: "release vs alpha", v1: "1.1", v2: "1.1_alpha1", want: 1},

		// Build component comparisons
		{name: "build versions", v1: "1.0.4-r3", v2: "1.0.4-r4", want: -1},
		{name: "build versions reverse", v1: "1.0.4-r4", v2: "1.0.4-r3", want: 1},

		// Leading zero handling (the key fix)
		{name: "leading zero handling", v1: "4.5.14", v2: "4.09-r1", want: 1},

		// Unknown suffix handling (the key fix)
		{name: "unknown suffix vs known suffix", v1: "23_foo", v2: "4_beta", want: 1},

		// Invalid format handling (the key fix)
		{name: "standard vs invalid format", v1: "1.0", v2: "1.0bc", want: -1},
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

func TestVersion_Compare_Fixture(t *testing.T) {
	e := &Ecosystem{}

	filename := "testdata/compare.txt"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to read fixture file %q: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := -1
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		line = removeComments(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")

		name := fmt.Sprintf("%s:%d: %s", filename, lineNumber, line)

		t.Run(name, func(t *testing.T) {
			if len(parts) != 3 {
				t.Fatalf("Invalid line format. Expected \"v1 [<|=|>] v2\", got: %q", line)
			}
			v1Str := parts[0]
			symbol := parts[1]
			v2Str := parts[2]
			symbolToCompare := map[string]int{
				"<": -1,
				"=": 0,
				">": 1,
			}
			want, ok := symbolToCompare[symbol]
			if !ok {
				t.Fatalf("Invalid comparison operator in line: %q", line)
			}

			v1, err := e.NewVersion(v1Str)
			if err != nil {
				t.Fatalf("NewVersion(%q) error: %v", v1Str, err)
			}
			v2, err := e.NewVersion(parts[2])
			if err != nil {
				t.Fatalf("NewVersion(%q) error: %v", v2Str, err)
			}

			got := v1.Compare(v2)
			if got != want {
				t.Errorf("Compare(%q, %q) = %d, want %d", v1, v2, got, want)
			}
		})
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading fixture file: %v", err)
	}
}

// Helper functions
func equalVersions(a, b *Version) bool {
	if a.original != b.original || a.letter != b.letter || a.hash != b.hash || a.build != b.build {
		return false
	}

	// Compare numeric components
	if len(a.numeric) != len(b.numeric) {
		return false
	}
	for i := range a.numeric {
		if a.numeric[i].value != b.numeric[i].value || a.numeric[i].originalStr != b.numeric[i].originalStr {
			return false
		}
	}

	// Compare suffixes
	if len(a.suffixes) != len(b.suffixes) {
		return false
	}
	for i := range a.suffixes {
		if a.suffixes[i].name != b.suffixes[i].name || a.suffixes[i].number != b.suffixes[i].number {
			return false
		}
	}

	return true
}

func removeComments(line string) string {
	if idx := strings.Index(line, "#"); idx != -1 {
		line = strings.TrimSpace(line[:idx])
	}
	return strings.TrimSpace(line)
}