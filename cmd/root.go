package cmd

import (
	"fmt"
	"lumen/internal/index"
	"lumen/internal/search"
	"lumen/internal/tui"
	"lumen/internal/vault"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	flagVault   string
	flagTags    bool
	flagHeads   bool
	flagContent bool
	flagRank    bool
	flagEditor  string
)

func run(cmd *cobra.Command, args []string) error {
	vaults := vault.DefaultVaults()

	vname, remainingArgs := resolveVaultName(vaults, args)
	vpath, err := resolveVaultPath(vaults, vname)
	if err != nil {
		return err
	}

	idx, err := index.Build(vpath)
	if err != nil {
		return err
	}

	cfg := tui.Config{
		VaultName:  vname,
		VaultPath:  vpath,
		Vaults:     vaults,
		Mode:       getSearchMode(),
		Rank:       flagRank,
		EditorCmd:  getEditor(),
		InitialQ:   strings.Join(remainingArgs, " "),
		ShowHidden: false,
	}

	m := tui.NewModel(cfg, idx)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func resolveVaultName(vaults vault.Vaults, args []string) (string, []string) {
	if vname := strings.TrimSpace(flagVault); vname != "" {
		return vname, args
	}

	if len(args) > 0 && vaults.Has(args[0]) {
		return args[0], args[1:]
	}

	return "develop", args
}

func resolveVaultPath(vaults vault.Vaults, vname string) (string, error) {
	vpath, ok := vaults.Path(vname)
	if !ok {
		return "", fmt.Errorf("unknown vault: %s (known: %s)", vname, strings.Join(vaults.Names(), ", "))
	}
	return expandHome(vpath), nil
}

func getSearchMode() search.Mode {
	switch {
	case flagTags:
		return search.ModeTags
	case flagHeads:
		return search.ModeHeadings
	case flagContent:
		return search.ModeContent
	default:
		return search.ModeAll
	}
}

func getEditor() string {
	if flagEditor != "" {
		return flagEditor
	}
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	return "nvim"
}

var rootCmd = &cobra.Command{
	Use:   "lumen [vault?] [query...]",
	Short: "Search and open Obsidian notes in a TUI (Cobra + Bubble Tea)",
	Args:  cobra.ArbitraryArgs,
	RunE:  run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&flagVault, "vault", "v", "", "vault name (develop/work/private)")
	rootCmd.Flags().BoolVarP(&flagTags, "tags", "t", false, "search only tags")
	rootCmd.Flags().BoolVarP(&flagHeads, "headings", "H", false, "search only headings/title/filename")
	rootCmd.Flags().BoolVarP(&flagContent, "content", "c", false, "search only content")
	rootCmd.Flags().BoolVarP(&flagRank, "rank", "r", false, "enable ranking (title/tags > content)")
	rootCmd.Flags().StringVar(&flagEditor, "editor", "", "editor command (default $EDITOR or nvim)")
}

func expandHome(p string) string {
	if p == "" {
		return p
	}
	if strings.HasPrefix(p, "~"+string(os.PathSeparator)) || p == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return p
		}
		if p == "~" {
			return home
		}
		return filepath.Join(home, p[2:])
	}
	return p
}
