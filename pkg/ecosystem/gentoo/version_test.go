package gentoo

import (
	"testing"
)

func TestEcosystem_NewVersion(t *testing.T) {
	tests := []struct {
		version string
		wantErr bool
	}{
		// Valid versions
		{"1.0", false},
		{"1.2.3", false},
		{"1.0a", false},
		{"2.3.4b", false},
		{"1.0_alpha", false},
		{"1.0_alpha1", false},
		{"1.0_beta", false},
		{"1.0_beta2", false},
		{"1.0_pre", false},
		{"1.0_pre3", false},
		{"1.0_rc", false},
		{"1.0_rc1", false},
		{"1.0_p", false},
		{"1.0_p1", false},
		{"1.0-r1", false},
		{"1.2.3-r10", false},
		{"1.0_alpha1-r2", false},
		{"2.4.6c_beta3-r5", false},
		{"10.20.30", false},
		{"0.1", false},
		{"1", false},

		// Invalid versions
		{"", true},
		{"abc", true},
		{"1.2.3.4.5.6.7.8.9.10.11.12", true}, // Too many components
		{"1.2-", true},
		{"1._alpha", true},
		{"1.0_invalid", true},
		{"1.0-", true},
		{"-r1", true},
		{"1.0-r", true},
		{"1.0-r-1", true},
		{"1.0__alpha", true},
		{"1.0..", true},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			v, err := e.NewVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersion(%q) error = %v, wantErr %v", tt.version, err, tt.wantErr)
				return
			}
			if !tt.wantErr && v.String() != tt.version {
				t.Errorf("NewVersion(%q).String() = %v, want %v", tt.version, v.String(), tt.version)
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		v1   string
		v2   string
		want int
	}{
		// Basic numeric comparisons
		{"1.0", "1.0", 0},
		{"1.0", "1.1", -1},
		{"1.1", "1.0", 1},
		{"1.2.3", "1.2.3", 0},
		{"1.2.3", "1.2.4", -1},
		{"1.2.4", "1.2.3", 1},
		{"2.0.0", "1.9.9", 1},
		{"1.10", "1.2", 1}, // 10 > 2

		// Different number of components (should be equivalent per Gentoo PMS)
		{"1.0", "1.0.0", 0}, // 1.0 equivalent to 1.0.0
		{"1.0.0", "1.0", 0},
		{"1.2", "1.2.0", 0},
		{"1.2.0", "1.2", 0},

		// Letter suffixes
		{"1.0a", "1.0", 1},   // letter > no letter
		{"1.0", "1.0a", -1},  // no letter < letter
		{"1.0a", "1.0b", -1}, // a < b
		{"1.0b", "1.0a", 1},  // b > a
		{"1.0a", "1.0a", 0},  // same letter

		// Version suffixes (alpha < beta < pre < rc < release < p)
		{"1.0_alpha", "1.0", -1},      // alpha < release
		{"1.0", "1.0_alpha", 1},       // release > alpha
		{"1.0_alpha", "1.0_beta", -1}, // alpha < beta
		{"1.0_beta", "1.0_alpha", 1},  // beta > alpha
		{"1.0_beta", "1.0_pre", -1},   // beta < pre
		{"1.0_pre", "1.0_beta", 1},    // pre > beta
		{"1.0_pre", "1.0_rc", -1},     // pre < rc
		{"1.0_rc", "1.0_pre", 1},      // rc > pre
		{"1.0_rc", "1.0", -1},         // rc < release
		{"1.0", "1.0_rc", 1},          // release > rc
		{"1.0", "1.0_p", -1},          // release < p
		{"1.0_p", "1.0", 1},           // p > release
		{"1.0_alpha", "1.0_alpha", 0}, // same suffix

		// Suffix numbers
		{"1.0_alpha1", "1.0_alpha2", -1}, // alpha1 < alpha2
		{"1.0_alpha2", "1.0_alpha1", 1},  // alpha2 > alpha1
		{"1.0_alpha1", "1.0_alpha1", 0},  // same
		{"1.0_beta1", "1.0_beta10", -1},  // 1 < 10
		{"1.0_rc", "1.0_rc1", -1},        // no number < number
		{"1.0_rc1", "1.0_rc", 1},         // number > no number

		// Revisions
		{"1.0-r1", "1.0", 1},     // revision > no revision
		{"1.0", "1.0-r1", -1},    // no revision < revision
		{"1.0-r1", "1.0-r2", -1}, // r1 < r2
		{"1.0-r2", "1.0-r1", 1},  // r2 > r1
		{"1.0-r1", "1.0-r1", 0},  // same revision
		{"1.0-r10", "1.0-r2", 1}, // r10 > r2

		// Complex combinations
		{"1.2.3a_alpha1-r2", "1.2.3a_alpha1-r2", 0},  // identical complex versions
		{"1.2.3a_alpha1-r1", "1.2.3a_alpha1-r2", -1}, // different revision
		{"1.2.3a_alpha1-r2", "1.2.3a_beta1-r1", -1},  // alpha < beta
		{"1.2.3b_rc2-r3", "1.2.3a_p1-r1", 1},         // higher base version
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.v1+"_vs_"+tt.v2, func(t *testing.T) {
			v1, err := e.NewVersion(tt.v1)
			if err != nil {
				t.Fatalf("NewVersion(%q) failed: %v", tt.v1, err)
			}
			v2, err := e.NewVersion(tt.v2)
			if err != nil {
				t.Fatalf("NewVersion(%q) failed: %v", tt.v2, err)
			}

			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Compare(%q, %q) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}
