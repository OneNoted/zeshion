package zeshcli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand(version string) *cobra.Command {
	base := NewBaseDeps()

	rootCmd := &cobra.Command{
		Use:              "zesh",
		Version:          version,
		Short:            "Smart session manager for tmux and Zellij",
		Long:             "Zesh is a smart terminal session manager that helps you create and manage tmux and Zellij sessions quickly and easily using zoxide.",
		TraverseChildren: true,
	}

	rootCmd.PersistentFlags().StringP("config", "C", "", "path to config file")
	rootCmd.PersistentFlags().StringP("multiplexer", "m", "auto", "multiplexer to use: auto, tmux, or zellij")

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
