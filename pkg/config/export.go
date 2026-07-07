package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// ExportedFunction represents a custom function script body embedded in the manifest.
type ExportedFunction struct {
	Name string `toml:"name"`
	Code string `toml:"code"`
	Ext  string `toml:"ext"`
}

// ExportedScript represents a script library body embedded in the manifest.
type ExportedScript struct {
	Category string `toml:"category"`
	Name     string `toml:"name"`
	Code     string `toml:"code"`
}

// Manifest represents the single consolidated reshell.toml configuration schema.
type Manifest struct {
	Package struct {
		Name        string `toml:"name"`
		Description string `toml:"description"`
	} `toml:"package"`
	Aliases   []Alias            `toml:"aliases"`
	Variables []EnvVar           `toml:"variables"`
	Snippets  []Snippet          `toml:"snippets"`
	Functions []ExportedFunction `toml:"functions"`
	Scripts   []ExportedScript   `toml:"scripts"`
}

// ExportConfig compiles all config files and physical scripts into a single reshell.toml manifest.
func ExportConfig(destTomlPath string) error {
	var manifest Manifest
	manifest.Package.Name = "reshell-profile"
	manifest.Package.Description = "Consolidated reshell configuration profile package"

	// 1. Load Aliases
	if aliases, err := LoadAliases(); err == nil {
		manifest.Aliases = aliases.Aliases
	}

	// 2. Load Environment Variables
	if envs, err := LoadEnv(); err == nil {
		manifest.Variables = envs.Variables
	}

	// 3. Load Snippets
	if snippets, err := LoadSnippets(); err == nil {
		manifest.Snippets = snippets.Snippets
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// 4. Load Custom Functions
	funcsDir := filepath.Join(configDir, "functions")
	if files, err := os.ReadDir(funcsDir); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			ext := filepath.Ext(file.Name())
			if ext == ".sh" || ext == ".fish" {
				name := strings.TrimSuffix(file.Name(), ext)
				body, err := os.ReadFile(filepath.Join(funcsDir, file.Name()))
				if err == nil {
					manifest.Functions = append(manifest.Functions, ExportedFunction{
						Name: name,
						Code: string(body),
						Ext:  ext,
					})
				}
			}
		}
	}

	// 5. Load Script Library
	scriptsDir := filepath.Join(configDir, "scripts")
	if _, err := os.Stat(scriptsDir); err == nil {
		err = filepath.Walk(scriptsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			ext := filepath.Ext(path)
			if ext == ".sh" {
				rel, err := filepath.Rel(scriptsDir, path)
				if err == nil {
					dir, file := filepath.Split(rel)
					category := strings.TrimSuffix(filepath.ToSlash(dir), "/")
					if category == "" {
						category = "general"
					}
					name := strings.TrimSuffix(file, ".sh")
					body, err := os.ReadFile(path)
					if err == nil {
						manifest.Scripts = append(manifest.Scripts, ExportedScript{
							Category: category,
							Name:     name,
							Code:     string(body),
						})
					}
				}
			}
			return nil
		})
	}

	// 6. Serialize and write to target path
	data, err := toml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to serialize manifest: %w", err)
	}

	return os.WriteFile(destTomlPath, data, 0644)
}

// ImportConfig parses a single reshell.toml manifest and extracts all assets.
func ImportConfig(srcTomlPath string) error {
	// Trust verification prompt
	fmt.Print("WARNING: Importing configurations runs scripts and custom functions on your machine.\nDo you trust this configuration manifest source? (y/N): ")
	var response string
	_, errScan := fmt.Scanln(&response)
	if errScan != nil {
		return fmt.Errorf("import aborted: failed to read trust confirmation")
	}
	response = strings.ToLower(strings.TrimSpace(response))
	if response != "y" && response != "yes" {
		return fmt.Errorf("import aborted: untrusted manifest source")
	}

	data, err := os.ReadFile(srcTomlPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest file: %w", err)
	}

	var manifest Manifest
	if err := toml.Unmarshal(data, &manifest); err != nil {
		return fmt.Errorf("invalid manifest file schema: %w", err)
	}

	// 1. Import Aliases
	aliases := AliasConfig{Aliases: manifest.Aliases}
	if err := SaveAliases(&aliases); err != nil {
		return fmt.Errorf("failed to import aliases: %w", err)
	}

	// 2. Import Env variables with secret checking
	var filteredVars []EnvVar
	for _, v := range manifest.Variables {
		if IsSecret(v.Name, v.Value) {
			fmt.Printf("\n[WARNING] Potential secret detected in imported env variable %q:\n", v.Name)
			fmt.Println("  (Note: Although reshell's Git repository is purely local and never pushed to any remote, plaintext storage is discouraged.)")
			fmt.Println("  1) Skip importing (keep in original profile) [Recommended]")
			fmt.Println("  2) Import in plaintext anyway")
			var choice int
			for {
				fmt.Print("Select choice (1-2): ")
				_, scanErr := fmt.Scanln(&choice)
				if scanErr == nil && (choice == 1 || choice == 2) {
					break
				}
				if scanErr != nil {
					if scanErr.Error() == "EOF" {
						choice = 1
						break
					}
					var discard string
					_, _ = fmt.Scanln(&discard)
				}
				fmt.Println("Invalid choice. Enter 1 or 2.")
			}
			if choice == 1 {
				continue
			}
		}
		filteredVars = append(filteredVars, v)
	}

	envs := EnvConfig{Variables: filteredVars}
	if err := SaveEnv(&envs); err != nil {
		return fmt.Errorf("failed to import environment variables: %w", err)
	}

	// 3. Import Snippets
	snippets := SnippetConfig{Snippets: manifest.Snippets}
	if err := SaveSnippets(&snippets); err != nil {
		return fmt.Errorf("failed to import snippets: %w", err)
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// 4. Extract Custom Functions
	funcsDir := filepath.Join(configDir, "functions")
	_ = os.MkdirAll(funcsDir, 0700)
	for _, fn := range manifest.Functions {
		if !IsValidName(fn.Name) {
			return fmt.Errorf("security error: invalid custom function name: %q", fn.Name)
		}
		ext := fn.Ext
		if ext == "" {
			ext = ".sh" // fallback
		}
		path, err := SafeJoin(funcsDir, fn.Name+ext)
		if err != nil {
			return err
		}
		_ = os.WriteFile(path, []byte(fn.Code), 0700)
	}

	// 5. Extract Script Library
	scriptsDir := filepath.Join(configDir, "scripts")
	_ = os.MkdirAll(scriptsDir, 0700)
	for _, scr := range manifest.Scripts {
		if !IsValidName(scr.Name) {
			return fmt.Errorf("security error: invalid script name: %q", scr.Name)
		}
		// Validate that category does not try to traverse up
		catDir, err := SafeJoin(scriptsDir, scr.Category)
		if err != nil {
			return err
		}
		_ = os.MkdirAll(catDir, 0700)
		path, err := SafeJoin(catDir, scr.Name+".sh")
		if err != nil {
			return err
		}
		_ = os.WriteFile(path, []byte(scr.Code), 0700)
	}

	return nil
}
