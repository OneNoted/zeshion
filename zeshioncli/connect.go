package zeshioncli

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"github.com/notes/zeshion/lister"
	"github.com/notes/zeshion/model"
)

func NewConnectCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "connect",
		Aliases: []string{"cn"},
		Short:   "Connect to the given session",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("please provide a session name")
			}
			name := strings.Join(args, " ")
			if name == "" {
				return nil
			}

			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			switchFlag, _ := cmd.Flags().GetBool("switch")
			command, _ := cmd.Flags().GetString("command")
			tmuxinator, _ := cmd.Flags().GetBool("tmuxinator")
			root, _ := cmd.Flags().GetBool("root")

			if root {
				hasRootDir, rootDir := base.Dir.RootDir(name)
				if hasRootDir {
					name = rootDir
				}
			}

			opts := model.ConnectOpts{Switch: switchFlag, Command: command, Tmuxinator: tmuxinator}
			trimmedName := deps.Icon.RemoveIcon(name)
			if _, err := deps.Connector.Connect(trimmedName, opts); err != nil {
				return err
			}
			// Refresh after connecting so the next list command has fresh data.
			if deps.CachingLister != nil {
				deps.CachingLister.RefreshCache(lister.ListOptions{})
				deps.CachingLister.Wait()
			}
			return nil
		},
	}

	cmd.Flags().BoolP("switch", "s", false, "Switch the session (rather than attach). This is useful for actions triggered outside the terminal.")
	cmd.Flags().StringP("command", "c", "", "Execute a command when connecting to a new session. Will be ignored if the session exists.")
	cmd.Flags().BoolP("tmuxinator", "T", false, "Use tmuxinator to start session if it doesn't exist")
	cmd.Flags().BoolP("root", "r", false, "Switches to the root of the current session")

	return cmd
}
