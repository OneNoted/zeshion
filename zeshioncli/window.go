package zeshioncli

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/OneNoted/zeshion/model"
	"github.com/OneNoted/zeshion/tmux"
)

func NewWindowCommand(base *BaseDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "window",
		Aliases: []string{"w"},
		Short:   "List or switch/create windows/tabs in the selected multiplexer session",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := buildDeps(cmd, base)
			if err != nil {
				return err
			}

			targetSession, _ := cmd.Flags().GetString("session")
			jsonOutput, _ := cmd.Flags().GetBool("json")

			if targetSession == "" {
				if !deps.Mux.IsAttached() {
					return fmt.Errorf("not inside a %s session, use --session to specify one", deps.MuxName)
				}
			} else {
				sessions, err := deps.Mux.ListSessions()
				if err != nil {
					return err
				}
				if !sessionTargetExists(sessions, targetSession) {
					return fmt.Errorf("session '%s' not found", targetSession)
				}
			}

			if len(args) == 0 {
				windows, err := deps.Mux.ListWindows(targetSession)
				if err != nil {
					return err
				}
				if jsonOutput {
					out, err := json.Marshal(windows)
					if err != nil {
						return err
					}
					fmt.Println(string(out))
					return nil
				}
				for _, w := range windows {
					fmt.Println(w.Name)
				}
				return nil
			}

			name := strings.Join(args, " ")

			windows, err := deps.Mux.ListWindows(targetSession)
			if err != nil {
				return err
			}
			target, found, err := resolveWindowTarget(windows, deps.MuxName, targetSession, name)
			if err != nil {
				return err
			}
			if found {
				if _, err := deps.Mux.SelectWindow(target); err != nil {
					return fmt.Errorf("failed to select window '%s': %w", name, err)
				}
				return nil
			}

			expanded, err := base.Home.ExpandPath(name)
			if err != nil {
				return err
			}
			isDir, absPath := base.Dir.Dir(expanded)
			if !isDir {
				return fmt.Errorf("'%s' is not an existing window or valid directory", name)
			}
			windowName := filepath.Base(absPath)

			if err := createWindow(deps.Mux, deps.MuxName, windowName, absPath, targetSession); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringP("session", "s", "", "target session (default: current attached session)")
	cmd.Flags().BoolP("json", "j", false, "output as json (list mode only)")

	return cmd
}

func resolveWindowTarget(windows []*model.TmuxWindow, muxName, targetSession, name string) (string, bool, error) {
	if muxName == "herdr" {
		for _, window := range windows {
			if window.ID == name {
				return window.ID, true, nil
			}
		}
	}

	matches := make([]*model.TmuxWindow, 0, 1)
	for _, window := range windows {
		if window.Name == name {
			matches = append(matches, window)
		}
	}
	if len(matches) == 0 {
		return "", false, nil
	}
	if len(matches) > 1 && muxName == "herdr" {
		ids := make([]string, len(matches))
		for index, window := range matches {
			ids[index] = window.ID
		}
		return "", false, fmt.Errorf("window %q is ambiguous; use an ID: %s", name, strings.Join(ids, ", "))
	}

	target := matches[0].Name
	switch {
	case muxName == "herdr":
		target = matches[0].ID
	case muxName == "tmux" && targetSession != "":
		target = fmt.Sprintf("%s:%s", targetSession, matches[0].ID)
	case muxName == "zellij" && targetSession != "":
		target = fmt.Sprintf("%s:%s", targetSession, matches[0].Name)
	}
	return target, true, nil
}

func sessionTargetExists(sessions []*model.TmuxSession, target string) bool {
	for _, session := range sessions {
		if session.Name == target || session.ID == target {
			return true
		}
	}
	return false
}

func createWindow(mux tmux.Tmux, muxName, name, path, targetSession string) error {
	id, err := mux.NewWindowInSession(name, path, targetSession, "")
	if err != nil {
		return fmt.Errorf("failed to create window: %w", err)
	}
	if muxName == "herdr" {
		if _, err := mux.SelectWindow(id); err != nil {
			return fmt.Errorf("failed to select window '%s': %w", name, err)
		}
	}
	return nil
}
