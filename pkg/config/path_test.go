package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Valid alphanumeric", "myVar123", true},
		{"Valid hyphen/underscore", "my-var_123", true},
		{"Empty string", "", false},
		{"Invalid dot", "my.var", false},
		{"Invalid space", "my var", false},
		{"Invalid slash", "my/var", false},
		{"Invalid backslash", "my\\var", false},
		{"Invalid backtick", "my`var", false},
		{"Invalid dollar", "my$var", false},
		{"Traversal attempt", "../var", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidName(tt.input); got != tt.want {
				t.Errorf("IsValidName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSafeJoin(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "reshell_test_safejoin")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Resolve the base directory absolute path
	baseDir, err := filepath.Abs(tempDir)
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	tests := []struct {
		name        string
		subPath     string
		wantError   bool
		expectedSub string
	}{
		{"Simple subpath", "sub", false, "sub"},
		{"Nested subpath", "sub/nested", false, "sub/nested"},
		{"Path cleanup", "sub/../sub/nested", false, "sub/nested"},
		{"Traversal escaping boundary", "../outside", true, ""},
		{"Traversal nested escape", "sub/../../outside", true, ""},
		{"Absolute path traversal", "/etc/passwd", true, ""},
		{"Base directory itself", "", false, ""},
		{"Dot path", ".", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeJoin(baseDir, tt.subPath)
			if (err != nil) != tt.wantError {
				t.Errorf("SafeJoin(%q, %q) returned error = %v, wantError = %v", baseDir, tt.subPath, err, tt.wantError)
				return
			}
			if !tt.wantError {
				expected := filepath.Join(baseDir, tt.expectedSub)
				if got != expected {
					t.Errorf("SafeJoin(%q, %q) = %q, want %q", baseDir, tt.subPath, got, expected)
				}
			}
		})
	}
}
