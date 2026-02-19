package cmd

import (
	"fmt"
	"lumen/internal/search"
	"lumen/internal/vault"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTestVaults() vault.Vaults {
	return vault.NewVaults(map[string]string{
		"develop": "~/test/develop",
		"work":    "~/test/work",
		"private": "~/test/private",
	})
}

func resetFlags() {
	flagVault = ""
	flagTags = false
	flagHeads = false
	flagContent = false
	flagRank = false
	flagEditor = ""
}

func TestResolveVaultName(t *testing.T) {
	vaults := createTestVaults()

	tests := []struct {
		name              string
		flagVaultValue    string
		args              []string
		expectedVault     string
		expectedRemaining []string
	}{
		{
			name:              "vault from flag",
			flagVaultValue:    "work",
			args:              []string{"query", "text"},
			expectedVault:     "work",
			expectedRemaining: []string{"query", "text"},
		},
		{
			name:              "vault from flag with whitespace",
			flagVaultValue:    "  private  ",
			args:              []string{"query"},
			expectedVault:     "private",
			expectedRemaining: []string{"query"},
		},
		{
			name:              "vault from first arg",
			flagVaultValue:    "",
			args:              []string{"work", "search", "term"},
			expectedVault:     "work",
			expectedRemaining: []string{"search", "term"},
		},
		{
			name:              "default vault when no match",
			flagVaultValue:    "",
			args:              []string{"query", "text"},
			expectedVault:     "develop",
			expectedRemaining: []string{"query", "text"},
		},
		{
			name:              "default vault when no args",
			flagVaultValue:    "",
			args:              []string{},
			expectedVault:     "develop",
			expectedRemaining: []string{},
		},
		{
			name:              "flag takes precedence over args",
			flagVaultValue:    "private",
			args:              []string{"work", "query"},
			expectedVault:     "private",
			expectedRemaining: []string{"work", "query"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags()
			flagVault = tt.flagVaultValue

			vault, remaining := resolveVaultName(vaults, tt.args)

			if vault != tt.expectedVault {
				t.Errorf("expected vault %q, got %q", tt.expectedVault, vault)
			}

			if len(remaining) != len(tt.expectedRemaining) {
				t.Errorf("expected %d remaining args, got %d", len(tt.expectedRemaining), len(remaining))
			}

			for i, arg := range remaining {
				if i >= len(tt.expectedRemaining) || arg != tt.expectedRemaining[i] {
					t.Errorf("remaining args mismatch at index %d: expected %v, got %v", i, tt.expectedRemaining, remaining)
					break
				}
			}
		})
	}
}

func TestResolveVaultPath(t *testing.T) {
	vaults := createTestVaults()

	tests := []struct {
		name          string
		vaultName     string
		expectError   bool
		errorContains string
	}{
		{
			name:        "valid vault - develop",
			vaultName:   "develop",
			expectError: false,
		},
		{
			name:        "valid vault - work",
			vaultName:   "work",
			expectError: false,
		},
		{
			name:        "valid vault - private",
			vaultName:   "private",
			expectError: false,
		},
		{
			name:          "invalid vault",
			vaultName:     "nonexistent",
			expectError:   true,
			errorContains: "unknown vault: nonexistent",
		},
		{
			name:          "empty vault name",
			vaultName:     "",
			expectError:   true,
			errorContains: "unknown vault:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := resolveVaultPath(vaults, tt.vaultName)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error containing %q, got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if path == "" {
					t.Error("expected non-empty path")
				}
			}
		})
	}
}

func TestGetSearchMode(t *testing.T) {
	tests := []struct {
		name         string
		tags         bool
		heads        bool
		content      bool
		expectedMode search.Mode
	}{
		{
			name:         "mode all (default)",
			tags:         false,
			heads:        false,
			content:      false,
			expectedMode: search.ModeAll,
		},
		{
			name:         "mode tags",
			tags:         true,
			heads:        false,
			content:      false,
			expectedMode: search.ModeTags,
		},
		{
			name:         "mode headings",
			tags:         false,
			heads:        true,
			content:      false,
			expectedMode: search.ModeHeadings,
		},
		{
			name:         "mode content",
			tags:         false,
			heads:        false,
			content:      true,
			expectedMode: search.ModeContent,
		},
		{
			name:         "tags takes precedence",
			tags:         true,
			heads:        true,
			content:      true,
			expectedMode: search.ModeTags,
		},
		{
			name:         "headings over content",
			tags:         false,
			heads:        true,
			content:      true,
			expectedMode: search.ModeHeadings,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags()
			flagTags = tt.tags
			flagHeads = tt.heads
			flagContent = tt.content

			mode := getSearchMode()

			if mode != tt.expectedMode {
				t.Errorf("expected mode %v, got %v", tt.expectedMode, mode)
			}
		})
	}
}

func TestGetEditor(t *testing.T) {
	// Save original env
	originalEditor := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", originalEditor)

	tests := []struct {
		name           string
		flagValue      string
		envValue       string
		expectedEditor string
	}{
		{
			name:           "flag takes precedence",
			flagValue:      "vim",
			envValue:       "emacs",
			expectedEditor: "vim",
		},
		{
			name:           "env var when no flag",
			flagValue:      "",
			envValue:       "nano",
			expectedEditor: "nano",
		},
		{
			name:           "default nvim when neither set",
			flagValue:      "",
			envValue:       "",
			expectedEditor: "nvim",
		},
		{
			name:           "flag only",
			flagValue:      "code",
			envValue:       "",
			expectedEditor: "code",
		},
		{
			name:           "empty flag uses env",
			flagValue:      "",
			envValue:       "vi",
			expectedEditor: "vi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags()
			flagEditor = tt.flagValue

			if tt.envValue != "" {
				os.Setenv("EDITOR", tt.envValue)
			} else {
				os.Unsetenv("EDITOR")
			}

			editor := getEditor()

			if editor != tt.expectedEditor {
				t.Errorf("expected editor %q, got %q", tt.expectedEditor, editor)
			}
		})
	}
}

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tilde only",
			input:    "~",
			expected: home,
		},
		{
			name:     "tilde with path",
			input:    fmt.Sprintf("~%sDocuments%stest", string(os.PathSeparator), string(os.PathSeparator)),
			expected: filepath.Join(home, "Documents", "test"),
		},
		{
			name:     "absolute path unchanged",
			input:    "/usr/local/bin",
			expected: "/usr/local/bin",
		},
		{
			name:     "relative path unchanged",
			input:    "test/path",
			expected: "test/path",
		},
		{
			name:     "empty string unchanged",
			input:    "",
			expected: "",
		},
		{
			name:     "tilde in middle unchanged",
			input:    "/path/to/~user",
			expected: "/path/to/~user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandHome(tt.input)

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
