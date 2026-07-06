package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStripQuotes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`'world'`, "world"},
		{`"nested 'quotes'"`, "nested 'quotes'"},
		{`simple`, "simple"},
		{`""`, ""},
	}

	for _, tc := range tests {
		res := stripQuotes(tc.input)
		if res != tc.expected {
			t.Errorf("stripQuotes(%q) = %q; expected %q", tc.input, res, tc.expected)
		}
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Clean Docker Containers", "clean-docker-containers"},
		{"git_status!", "git_status"},
		{"my...custom-snippet", "my...custom-snippet"},
		{"  spaces   everywhere  ", "spaces-everywhere"},
	}

	for _, tc := range tests {
		res := slugify(tc.input)
		if res != tc.expected {
			t.Errorf("slugify(%q) = %q; expected %q", tc.input, res, tc.expected)
		}
	}
}

func TestIsSecret(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"API_KEY", "12345", true},
		{"aws_secret_key", "xyz", true},
		{"GITHUB_TOKEN", "ghp_abc", true},
		{"PATH", "/usr/bin", false},
		{"MY_VAR", "ghp_1234567890", true},
	}

	for _, tc := range tests {
		res := IsSecret(tc.name, tc.value)
		if res != tc.expected {
			t.Errorf("IsSecret(%q, %q) = %v; expected %v", tc.name, tc.value, res, tc.expected)
		}
	}
}

func TestParseBashFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "reshell-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	content := `
alias ll="ls -l"
alias gs='git status'
export EDITOR=vim
export GITHUB_TOKEN="ghp_secure"

function myfunc() {
    echo "hello"
}

myfunc2() {
    echo "world"
}
`
	filePath := filepath.Join(tmpDir, ".bashrc")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	var results DiscoveryResults
	if err := parseBashFile(filePath, &results); err != nil {
		t.Fatal(err)
	}

	if len(results.Aliases) != 2 {
		t.Errorf("expected 2 aliases, got %d", len(results.Aliases))
	}
	if results.Aliases[0].Name != "ll" || results.Aliases[0].Value != "ls -l" {
		t.Errorf("unexpected alias: %+v", results.Aliases[0])
	}

	if len(results.EnvVars) != 2 {
		t.Errorf("expected 2 env vars, got %d", len(results.EnvVars))
	}
	if results.EnvVars[0].Name != "EDITOR" || results.EnvVars[0].Value != "vim" {
		t.Errorf("unexpected env var: %+v", results.EnvVars[0])
	}

	if len(results.Functions) != 2 {
		t.Errorf("expected 2 functions, got %d", len(results.Functions))
	}
	if results.Functions[0].Name != "myfunc" || !testingContains(results.Functions[0].Code, `echo "hello"`) {
		t.Errorf("unexpected function myfunc: %+v", results.Functions[0])
	}
	if results.Functions[1].Name != "myfunc2" || !testingContains(results.Functions[1].Code, `echo "world"`) {
		t.Errorf("unexpected function myfunc2: %+v", results.Functions[1])
	}
}

func TestParseFishFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "reshell-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	content := `
alias gp "git push"
set -gx MY_ENV "value"
setenv ANOTHER_ENV "value2"

function my_fish_func
    if true
        echo "fish"
    end
end
`
	filePath := filepath.Join(tmpDir, "config.fish")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	var results DiscoveryResults
	if err := parseFishFile(filePath, &results); err != nil {
		t.Fatal(err)
	}

	if len(results.Aliases) != 1 {
		t.Errorf("expected 1 alias, got %d", len(results.Aliases))
	}
	if results.Aliases[0].Name != "gp" || results.Aliases[0].Value != "git push" {
		t.Errorf("unexpected fish alias: %+v", results.Aliases[0])
	}

	if len(results.EnvVars) != 2 {
		t.Errorf("expected 2 env vars, got %d", len(results.EnvVars))
	}
	if results.EnvVars[0].Name != "MY_ENV" || results.EnvVars[0].Value != "value" {
		t.Errorf("unexpected env var: %+v", results.EnvVars[0])
	}

	if len(results.Functions) != 1 {
		t.Errorf("expected 1 function, got %d", len(results.Functions))
	}
	if results.Functions[0].Name != "my_fish_func" || !testingContains(results.Functions[0].Code, `echo "fish"`) {
		t.Errorf("unexpected function: %+v", results.Functions[0])
	}
}

func TestParseVSCodeSnippetFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "reshell-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	content := `{
  "Print to console": {
    "prefix": "log",
    "body": [
      "console.log('$1');",
      "$2"
    ],
    "description": "Log output to console"
  }
}`
	filePath := filepath.Join(tmpDir, "js.json")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	var results DiscoveryResults
	if err := parseVSCodeSnippetFile(filePath, &results); err != nil {
		t.Fatal(err)
	}

	if len(results.Snippets) != 1 {
		t.Errorf("expected 1 snippet, got %d", len(results.Snippets))
	}
	snip := results.Snippets[0]
	if snip.Name != "print-to-console" || !testingContains(snip.Code, "console.log") {
		t.Errorf("unexpected snippet: %+v", snip)
	}
}

func TestParsePetSnippetFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "reshell-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	content := `
[[snippets]]
  description = "Docker Clean"
  command = "docker rm -f $(docker ps -a -q)"
  tag = ["docker", "clean"]
`
	filePath := filepath.Join(tmpDir, "snippet.toml")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	var results DiscoveryResults
	if err := parsePetSnippetFile(filePath, &results); err != nil {
		t.Fatal(err)
	}

	if len(results.Snippets) != 1 {
		t.Errorf("expected 1 snippet, got %d", len(results.Snippets))
	}
	snip := results.Snippets[0]
	if snip.Name != "docker-clean" || snip.Code != "docker rm -f $(docker ps -a -q)" {
		t.Errorf("unexpected snippet: %+v", snip)
	}
}

func testingContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s[0:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr))
}

func containsMiddle(s, substr string) bool {
	for i := 1; i < len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
