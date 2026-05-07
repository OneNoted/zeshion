package zeshioncli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand(version string) *cobra.Command {
	base := NewBaseDeps()

	rootCmd := &cobra.Command{
		Use:              "zeshion",
		Version:          version,
		Short:            "Smart session manager for tmux, Zellij, and Herdr",
		Long:             "Zeshion is a smart terminal session manager that creates and manages tmux sessions, Zellij sessions, and Herdr workspaces.",
		TraverseChildren: true,
	}

	rootCmd.PersistentFlags().StringP("config", "C", "", "path to config file")
	rootCmd.PersistentFlags().StringP("multiplexer", "m", "auto", "multiplexer to use: auto, tmux, zellij, or herdr")

	rootCmd.AddCommand(
		NewListCommand(base),
		NewLastCommand(base),
		NewConnectCommand(base),
		NewCloneCommand(base),
		NewRootSessionCommand(base),
		NewPreviewCommand(base),
		NewPickerCommand(base),
		NewWindowCommand(base),
	)

	return rootCmd
}
