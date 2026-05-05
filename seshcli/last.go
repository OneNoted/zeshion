package seshcli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewLastCommand(base *BaseDeps) *cobra.Command {
	return &cobra.Command{
		Use:     "last",
		Aliases: []string{"L"},
		Short:   "Connect to the last session in the selected multiplexer",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			lastSession, exists := deps.Lister.GetLastTmuxSession()
			if deps.MuxName == "zellij" {
				lastSession, exists = deps.Lister.GetLastZellijSession()
			}
			if !exists {
				// TODO: silently fail?
				return fmt.Errorf("No last session found")
			}
			deps.Mux.SwitchClient(lastSession.Name)
			return nil
		},
	}
}
