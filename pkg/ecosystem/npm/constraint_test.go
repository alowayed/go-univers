package npm

import (
	"testing"
)

func TestConstraint_Matches(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		constVer string
		testVer  string
		want     bool
	}{
		{
			name:     "equals match",
			operator: "=",
			constVer: "1.2.3",
			testVer:  "1.2.3",
			want:     true,
		},
		{
			name:     "equals no match",
			operator: "=",
			constVer: "1.2.3",
			testVer:  "1.2.4",
		},
		{
			name:     "not equals match",
			operator: "!=",
			constVer: "1.2.3",
			testVer:  "1.2.4",
			want:     true,
		},
		{
			name:     "greater than match",
			operator: ">",
			constVer: "1.2.3",
			testVer:  "1.2.4",
			want:     true,
		},
		{
			name:     "greater than no match",
			operator: ">",
			constVer: "1.2.3",
			testVer:  "1.2.3",
		},
		{
			name:     "greater than or equal match",
			operator: ">=",
			constVer: "1.2.3",
			testVer:  "1.2.3",
			want:     true,
		},
		{
			name:     "less than match",
			operator: "<",
			constVer: "1.2.3",
			testVer:  "1.2.2",
			want:     true,
		},
		{
			name:     "less than or equal match",
			operator: "<=",
			constVer: "1.2.3",
			testVer:  "1.2.3",
			want:     true,
		},
		{
			name:     "wildcard match",
			operator: "*",
			constVer: "*",
			testVer:  "1.2.3",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &constraint{
				operator: tt.operator,
				version:  tt.constVer,
			}

			v, err := NewVersion(tt.testVer)
			if err != nil {
				t.Fatalf("Failed to parse test version: %v", err)
			}

			got := c.matches(v)
			if got != tt.want {
				t.Errorf("constraint.matches() = %v, want %v", got, tt.want)
			}
		})
	}
}