package snippets

import (
	"os"
	"testing"
)

func setupTestHome(t *testing.T) string {
	tempHome, err := os.MkdirTemp("", "reshell-snippets-test-*")
	if err != nil {
		t.Fatalf("failed to create temp home: %v", err)
	}

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)

	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
		os.RemoveAll(tempHome)
	})

	return tempHome
}

func TestAddRemoveSnippet(t *testing.T) {
	setupTestHome(t)

	// Add new snippet
	err := AddOrUpdate("test-snip", "echo 'hello'", "A test snippet", []string{"test", "hello"}, "bash", "all", false)
	if err != nil {
		t.Fatalf("AddOrUpdate failed: %v", err)
	}

	// Verify it can be searched
	results, err := Search("test")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "test-snip" {
		t.Errorf("expected snippet name 'test-snip', got %s", results[0].Name)
	}

	// Toggle favorite
	err = ToggleFavorite("test-snip")
	if err != nil {
		t.Fatalf("ToggleFavorite failed: %v", err)
	}

	results, _ = Search("test-snip")
	if !results[0].Favorite {
		t.Errorf("expected snippet to be favorite")
	}

	// Remove snippet
	err = Remove("test-snip")
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	results, _ = Search("test-snip")
	if len(results) != 0 {
		t.Errorf("expected snippet to be removed, but search returned %d results", len(results))
	}
}

func TestSnippetSearch(t *testing.T) {
	setupTestHome(t)

	_ = AddOrUpdate("golang-run", "go run .", "Run go program", []string{"go", "run"}, "go", "all", false)
	_ = AddOrUpdate("docker-build", "docker build -t test .", "Build container", []string{"docker", "build"}, "bash", "all", false)

	tests := []struct {
		query string
		count int
	}{
		{"go", 1},
		{"docker", 1},
		{"run", 1},
		{"build", 1},
		{"container", 1},
		{"nonexistent", 0},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			res, err := Search(tt.query)
			if err != nil {
				t.Fatalf("Search failed: %v", err)
			}
			if len(res) != tt.count {
				t.Errorf("Search(%q) returned %d items, expected %d", tt.query, len(res), tt.count)
			}
		})
	}
}
