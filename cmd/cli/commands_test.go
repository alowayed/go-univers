package cli

import (
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantCode int
	}{
		{
			name:     "no arguments",
			args:     []string{},
			wantCode: 1,
		},
		{
			name:     "unknown ecosystem",
			args:     []string{"unknown"},
			wantCode: 1,
		},
		{
			name:     "npm without command",
			args:     []string{"npm"},
			wantCode: 1,
		},
		{
			name:     "npm unknown command",
			args:     []string{"npm", "unknown"},
			wantCode: 1,
		},
		{
			name:     "npm compare equal versions",
			args:     []string{"npm", "compare", "1.2.3", "1.2.3"},
			wantCode: 0,
		},
		{
			name:     "npm compare less than",
			args:     []string{"npm", "compare", "1.2.3", "1.2.4"},
			wantCode: 0,
		},
		{
			name:     "npm compare greater than",
			args:     []string{"npm", "compare", "1.2.4", "1.2.3"},
			wantCode: 0,
		},
		{
			name:     "npm compare invalid version",
			args:     []string{"npm", "compare", "invalid", "1.2.3"},
			wantCode: 1,
		},
		{
			name:     "npm compare wrong number of args",
			args:     []string{"npm", "compare", "1.2.3"},
			wantCode: 1,
		},
		{
			name:     "npm sort single version",
			args:     []string{"npm", "sort", "1.2.3"},
			wantCode: 0,
		},
		{
			name:     "npm sort multiple versions",
			args:     []string{"npm", "sort", "2.0.0", "1.2.3", "1.10.0"},
			wantCode: 0,
		},
		{
			name:     "npm sort no args",
			args:     []string{"npm", "sort"},
			wantCode: 1,
		},
		{
			name:     "npm sort invalid version",
			args:     []string{"npm", "sort", "invalid", "1.2.3"},
			wantCode: 1,
		},
		{
			name:     "npm satisfies true",
			args:     []string{"npm", "satisfies", "1.2.5", "^1.2.0"},
			wantCode: 0,
		},
		{
			name:     "npm satisfies false",
			args:     []string{"npm", "satisfies", "2.0.0", "^1.2.0"},
			wantCode: 1,
		},
		{
			name:     "npm satisfies invalid version",
			args:     []string{"npm", "satisfies", "invalid", "^1.2.0"},
			wantCode: 1,
		},
		{
			name:     "npm satisfies invalid range",
			args:     []string{"npm", "satisfies", "1.2.3", "invalid"},
			wantCode: 1,
		},
		{
			name:     "npm satisfies wrong number of args",
			args:     []string{"npm", "satisfies", "1.2.3"},
			wantCode: 1,
		},
		{
			name:     "pypi compare equal versions",
			args:     []string{"pypi", "compare", "1.2.3", "1.2.3"},
			wantCode: 0,
		},
		{
			name:     "pypi compare with prerelease",
			args:     []string{"pypi", "compare", "1.0.0a1", "1.0.0"},
			wantCode: 0,
		},
		{
			name:     "pypi sort versions",
			args:     []string{"pypi", "sort", "1.0.0a1", "1.0.0", "0.9.0"},
			wantCode: 0,
		},
		{
			name:     "pypi satisfies compatible release",
			args:     []string{"pypi", "satisfies", "1.2.5", "~=1.2.0"},
			wantCode: 0,
		},
		{
			name:     "pypi satisfies wildcard",
			args:     []string{"pypi", "satisfies", "1.2.5", "==1.2.*"},
			wantCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output for tests that should succeed
			if tt.wantCode == 0 {
				// For successful tests, we just verify the exit code
				got := Run(tt.args)
				if got != tt.wantCode {
					t.Errorf("Run() = %v, want %v", got, tt.wantCode)
				}
			} else {
				// For error tests, we also just verify the exit code
				got := Run(tt.args)
				if got != tt.wantCode {
					t.Errorf("Run() = %v, want %v", got, tt.wantCode)
				}
			}
		})
	}
}

func TestRunEcosystem(t *testing.T) {
	tests := []struct {
		name      string
		ecosystem string
		args      []string
		wantCode  int
	}{
		{
			name:      "empty args",
			ecosystem: "npm",
			args:      []string{},
			wantCode:  1,
		},
		{
			name:      "unknown command",
			ecosystem: "npm",
			args:      []string{"unknown"},
			wantCode:  1,
		},
		{
			name:      "valid compare command",
			ecosystem: "npm",
			args:      []string{"compare", "1.2.3", "1.2.4"},
			wantCode:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runEcosystem(tt.ecosystem, tt.args)
			if got != tt.wantCode {
				t.Errorf("runEcosystem() = %v, want %v", got, tt.wantCode)
			}
		})
	}
}

// Integration test to verify the CLI works end-to-end
func TestCLIIntegration(t *testing.T) {
	// Test that we can build and run the CLI
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test compare command
	result := Run([]string{"npm", "compare", "1.2.3", "1.2.4"})
	if result != 0 {
		t.Errorf("Expected successful comparison, got exit code %d", result)
	}

	// Test sort command  
	result = Run([]string{"npm", "sort", "2.0.0", "1.2.3", "1.10.0"})
	if result != 0 {
		t.Errorf("Expected successful sort, got exit code %d", result)
	}

	// Test satisfies command (true case)
	result = Run([]string{"npm", "satisfies", "1.2.5", "^1.2.0"})
	if result != 0 {
		t.Errorf("Expected version to satisfy range, got exit code %d", result)
	}

	// Test satisfies command (false case)
	result = Run([]string{"npm", "satisfies", "2.0.0", "^1.2.0"})
	if result != 1 {
		t.Errorf("Expected version to not satisfy range, got exit code %d", result)
	}
}