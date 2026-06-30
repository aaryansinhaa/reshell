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

	// 2. Import Env variables
	envs := EnvConfig{Variables: manifest.Variables}
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
	_ = os.MkdirAll(funcsDir, 0755)
	for _, fn := range manifest.Functions {
		ext := fn.Ext
		if ext == "" {
			ext = ".sh" // fallback
		}
		path := filepath.Join(funcsDir, fn.Name+ext)
		_ = os.WriteFile(path, []byte(fn.Code), 0755)
	}

	// 5. Extract Script Library
	scriptsDir := filepath.Join(configDir, "scripts")
	_ = os.MkdirAll(scriptsDir, 0755)
	for _, scr := range manifest.Scripts {
		catDir := filepath.Join(scriptsDir, scr.Category)
		_ = os.MkdirAll(catDir, 0755)
		path := filepath.Join(catDir, scr.Name+".sh")
		_ = os.WriteFile(path, []byte(scr.Code), 0755)
	}

	return nil
}
