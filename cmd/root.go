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

var rootCmd = &cobra.Command{
	Use:   "lumen [vault?] [query...]",
	Short: "Search and open Obsidian notes in a TUI (Cobra + Bubble Tea)",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		vaults := vault.DefaultVaults()

		vname := strings.TrimSpace(flagVault)

		if vname == "" {
			if len(args) > 0 && vaults.Has(args[0]) {
				vname = args[0]
				args = args[1:]
			} else {
				vname = "develop"
			}
		}

		vpath, ok := vaults.Path(vname)

		if !ok {
			return fmt.Errorf("unknown vault: %s (known: %s)", vname, strings.Join(vaults.Names(), ", "))
		}
		vpath = expandHome(vpath)

		mode := search.ModeAll
		if flagTags {
			mode = search.ModeTags
		} else if flagHeads {
			mode = search.ModeHeadings
		} else if flagContent {
			mode = search.ModeContent
		}

		editor := flagEditor
		if editor == "" {
			editor = os.Getenv("EDITOR")
		}
		if editor == "" {
			editor = "nvim"
		}

		q := strings.Join(args, " ")

		idx, err := index.Build(vpath)
		if err != nil {
			return err
		}

		cfg := tui.Config{
			VaultName:  vname,
			VaultPath:  vpath,
			Vaults:     vaults,
			Mode:       mode,
			Rank:       flagRank,
			EditorCmd:  editor,
			InitialQ:   q,
			ShowHidden: false,
		}

		m := tui.NewModel(cfg, idx)
		p := tea.NewProgram(m, tea.WithAltScreen())
		_, err = p.Run()
		return err
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
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
