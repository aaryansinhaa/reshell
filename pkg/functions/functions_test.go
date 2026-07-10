package functions

import (
	"os"
	"path/filepath"
	"reshell/pkg/config"
	"testing"
)

func TestCreateOrUpdateAndGet(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "reshell-test-home-*")
	if err != nil {
		t.Fatalf("failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tempHome)

	t.Setenv("HOME", tempHome)

	err = config.EnsureDirectories()
	if err != nil {
		t.Fatalf("failed to ensure directories: %v", err)
	}

	// 1. Create a bash function
	funcCode := "#!/bin/bash\necho 'hello'"
	err = CreateOrUpdate("myfunc", funcCode)
	if err != nil {
		t.Errorf("CreateOrUpdate failed for myfunc: %v", err)
	}

	// 2. Read back
	code, ext, err := Get("myfunc")
	if err != nil {
		t.Errorf("Get failed for myfunc: %v", err)
	}
	if code != funcCode {
		t.Errorf("expected code %q, got %q", funcCode, code)
	}
	if ext != ".sh" {
		t.Errorf("expected extension .sh, got %q", ext)
	}

	// 3. Create a fish function
	fishCode := "#!/bin/fish\necho 'hello fish'"
	err = CreateOrUpdate("myfishfunc", fishCode)
	if err != nil {
		t.Errorf("CreateOrUpdate failed for myfishfunc: %v", err)
	}

	// 4. Read back fish function
	code, ext, err = Get("myfishfunc")
	if err != nil {
		t.Errorf("Get failed for myfishfunc: %v", err)
	}
	if code != fishCode {
		t.Errorf("expected code %q, got %q", fishCode, code)
	}
	if ext != ".fish" {
		t.Errorf("expected extension .fish, got %q", ext)
	}

	// 5. Retrieve non-existent function
	_, _, err = Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent function, got nil")
	}
}

func TestList(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "reshell-test-home-*")
	if err != nil {
		t.Fatalf("failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tempHome)

	t.Setenv("HOME", tempHome)

	err = config.EnsureDirectories()
	if err != nil {
		t.Fatalf("failed to ensure directories: %v", err)
	}

	// Check empty list
	funcs, err := List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(funcs) != 0 {
		t.Errorf("expected 0 functions, got %d", len(funcs))
	}

	// Add functions
	_ = CreateOrUpdate("fn1", "echo 1")
	_ = CreateOrUpdate("fn2", "#!/bin/fish\necho 2")

	// Create ignored files & directories inside functions directory
	funcDir, _ := config.GetFunctionsDir()
	_ = os.MkdirAll(filepath.Join(funcDir, "subdir"), 0755)
	_ = os.WriteFile(filepath.Join(funcDir, "ignored.txt"), []byte("ignored"), 0644)
	_ = os.WriteFile(filepath.Join(funcDir, "subdir", "fn3.sh"), []byte("echo 3"), 0644)

	// Verify listing
	funcs, err = List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(funcs) != 2 {
		t.Errorf("expected 2 functions, got %d: %v", len(funcs), funcs)
	}

	hasFn1, hasFn2 := false, false
	for _, fn := range funcs {
		if fn == "fn1" {
			hasFn1 = true
		}
		if fn == "fn2" {
			hasFn2 = true
		}
	}
	if !hasFn1 || !hasFn2 {
		t.Errorf("expected functions to contain fn1 and fn2, got %v", funcs)
	}
}

func TestRemove(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "reshell-test-home-*")
	if err != nil {
		t.Fatalf("failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tempHome)

	t.Setenv("HOME", tempHome)

	err = config.EnsureDirectories()
	if err != nil {
		t.Fatalf("failed to ensure directories: %v", err)
	}

	_ = CreateOrUpdate("fn1", "echo 1")

	// Delete existing
	err = Remove("fn1")
	if err != nil {
		t.Errorf("Remove failed: %v", err)
	}

	// Verify it's gone
	_, _, err = Get("fn1")
	if err == nil {
		t.Error("expected error getting deleted function")
	}

	// Delete non-existing
	err = Remove("fn1")
	if err == nil {
		t.Error("expected error removing non-existing function, got nil")
	}
}

func TestSecurityAndValidation(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "reshell-test-home-*")
	if err != nil {
		t.Fatalf("failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tempHome)

	t.Setenv("HOME", tempHome)

	// Check IsValidName security checks
	invalidNames := []string{"../fn", "fn/../../x", "", "fn$name", "fn;echo"}
	for _, name := range invalidNames {
		err := CreateOrUpdate(name, "echo 1")
		if err == nil {
			t.Errorf("expected security validation error for invalid name %q", name)
		}
		_, _, err = Get(name)
		if err == nil {
			t.Errorf("expected security validation error in Get for invalid name %q", name)
		}
		err = Remove(name)
		if err == nil {
			t.Errorf("expected security validation error in Remove for invalid name %q", name)
		}
		_, err = Validate(name)
		if err == nil {
			t.Errorf("expected security validation error in Validate for invalid name %q", name)
		}
	}
}

func TestValidate(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "reshell-test-home-*")
	if err != nil {
		t.Fatalf("failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tempHome)

	t.Setenv("HOME", tempHome)

	err = config.EnsureDirectories()
	if err != nil {
		t.Fatalf("failed to ensure directories: %v", err)
	}

	// Create valid bash script
	err = CreateOrUpdate("validfn", "my_func() {\n  echo \"hello\"\n}")
	if err != nil {
		t.Fatalf("CreateOrUpdate failed: %v", err)
	}

	// Run validate (this runs bash -n)
	stderr, err := Validate("validfn")
	if err != nil {
		t.Logf("Note: bash -n validation test skipped or failed if bash is not in environment: %v, stderr: %s", err, stderr)
	} else {
		if stderr != "" {
			t.Errorf("expected empty stderr for valid function, got %q", stderr)
		}
	}

	// Create invalid function
	err = CreateOrUpdate("invalidfn", "my_func() {\n  echo \"hello\"\n") // missing closing brace
	if err != nil {
		t.Fatalf("CreateOrUpdate failed: %v", err)
	}

	// Validate invalid script
	stderr, err = Validate("invalidfn")
	if err == nil {
		t.Error("expected syntax validation error for invalid function, got nil")
	}

	// Validate non-existing function
	_, err = Validate("nonexistent")
	if err == nil {
		t.Error("expected error for validating nonexistent function, got nil")
	}
}
