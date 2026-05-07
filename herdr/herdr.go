package herdr

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/OneNoted/zeshion/model"
	"github.com/OneNoted/zeshion/oswrap"
	"github.com/OneNoted/zeshion/shell"
)

// Herdr adapts Herdr workspaces, tabs, and panes to zeshion's session model.
type Herdr struct {
	os      oswrap.Os
	shell   shell.Shell
	bin     string
	session string
}

func NewHerdr(os oswrap.Os, shell shell.Shell, bin, session string) *Herdr {
	if bin == "" {
		bin = "herdr"
	}
	return &Herdr{os: os, shell: shell, bin: bin, session: session}
}

func (h *Herdr) IsAttached() bool {
	return h.os.Getenv("HERDR_ENV") == "1" || h.os.Getenv("HERDR_WORKSPACE_ID") != ""
}

func (h *Herdr) scopedArgs(args ...string) []string {
	if h.IsAttached() || h.session == "" || h.session == "default" {
		return args
	}
	scoped := make([]string, 0, len(args)+2)
	scoped = append(scoped, "--session", h.session)
	return append(scoped, args...)
}

func (h *Herdr) command(args ...string) (string, error) {
	return h.shell.Cmd(h.bin, h.scopedArgs(args...)...)
}

func decodeResponse[T any](output string) (T, error) {
	var response apiResponse[T]
	if err := json.Unmarshal([]byte(output), &response); err != nil {
		var zero T
		return zero, fmt.Errorf("couldn't decode Herdr response: %w", err)
	}
	return response.Result, nil
}

func (h *Herdr) snapshot() (snapshot, error) {
	output, err := h.command("api", "snapshot")
	if err != nil {
		return snapshot{}, fmt.Errorf("couldn't read Herdr state: %w", err)
	}
	result, err := decodeResponse[snapshotResult](output)
	if err != nil {
		return snapshot{}, err
	}
	return result.Snapshot, nil
}

func workspaceName(state snapshot, target workspace) string {
	if target.Label != "" {
		return target.Label
	}
	if path := workspacePath(state, target); path != "" {
		return filepath.Base(path)
	}
	return target.ID
}

func findWorkspace(state snapshot, target string) (workspace, error) {
	if target == "" {
		target = state.FocusedWorkspaceID
	}
	for _, candidate := range state.Workspaces {
		if candidate.ID == target {
			return candidate, nil
		}
	}

	matches := make([]workspace, 0, 1)
	for _, candidate := range state.Workspaces {
		if candidate.Label == target || workspaceName(state, candidate) == target {
			matches = append(matches, candidate)
		}
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) > 1 {
		ids := make([]string, len(matches))
		for i, match := range matches {
			ids[i] = match.ID
		}
		return workspace{}, fmt.Errorf("Herdr workspace label %q is ambiguous; use one of: %s", target, strings.Join(ids, ", "))
	}
	return workspace{}, fmt.Errorf("Herdr workspace %q not found", target)
}

func (h *Herdr) resolveWorkspace(state snapshot, target string) (workspace, error) {
	if target == "" {
		target = h.os.Getenv("HERDR_WORKSPACE_ID")
	}
	return findWorkspace(state, target)
}

func workspacePath(state snapshot, target workspace) string {
	if target.Worktree != nil && target.Worktree.CheckoutPath != "" {
		return target.Worktree.CheckoutPath
	}
	for _, candidate := range state.Panes {
		if candidate.WorkspaceID == target.ID && candidate.TabID == target.ActiveTabID && candidate.Cwd != "" {
			return candidate.Cwd
		}
	}
	for _, candidate := range state.Panes {
		if candidate.WorkspaceID == target.ID && candidate.Cwd != "" {
			return candidate.Cwd
		}
	}
	return ""
}

func panePath(target pane) string {
	if target.ForegroundCwd != "" {
		return target.ForegroundCwd
	}
	return target.Cwd
}

func paneForWorkspace(state snapshot, target workspace) (pane, error) {
	if target.ID == state.FocusedWorkspaceID || target.Focused {
		for _, candidate := range state.Panes {
			if candidate.WorkspaceID == target.ID &&
				(candidate.ID == state.FocusedPaneID || candidate.Focused) {
				return candidate, nil
			}
		}
	}
	for _, candidate := range state.Panes {
		if candidate.WorkspaceID == target.ID && candidate.TabID == target.ActiveTabID {
			return candidate, nil
		}
	}
	for _, candidate := range state.Panes {
		if candidate.WorkspaceID == target.ID {
			return candidate, nil
		}
	}
	return pane{}, fmt.Errorf("Herdr workspace %q has no panes", target.Label)
}

func (h *Herdr) ListSessions() ([]*model.TmuxSession, error) {
	state, err := h.snapshot()
	if err != nil {
		// Match tmux and Zellij listing semantics: an unavailable server is an empty source.
		return []*model.TmuxSession{}, nil
	}

	callerWorkspaceID := h.os.Getenv("HERDR_WORKSPACE_ID")
	sessions := make([]*model.TmuxSession, 0, len(state.Workspaces))
	for _, candidate := range state.Workspaces {
		attached := 0
		if callerWorkspaceID != "" {
			if candidate.ID == callerWorkspaceID {
				attached = 1
			}
		} else if candidate.ID == state.FocusedWorkspaceID || candidate.Focused {
			attached = 1
		}
		sessions = append(sessions, &model.TmuxSession{
			ID:       candidate.ID,
			Name:     workspaceName(state, candidate),
			Path:     workspacePath(state, candidate),
			Attached: attached,
			Windows:  candidate.TabCount,
		})
	}
	return sessions, nil
}

func (h *Herdr) ListWindows(targetSession string) ([]*model.TmuxWindow, error) {
	state, err := h.snapshot()
	if err != nil {
		return nil, err
	}
	target, err := h.resolveWorkspace(state, targetSession)
	if err != nil {
		return nil, err
	}

	windows := make([]*model.TmuxWindow, 0, target.TabCount)
	for _, candidate := range state.Tabs {
		if candidate.WorkspaceID != target.ID {
			continue
		}
		path := ""
		for _, p := range state.Panes {
			if p.TabID == candidate.ID {
				path = panePath(p)
				break
			}
		}
		name := candidate.Label
		if name == "" {
			name = fmt.Sprintf("tab-%d", candidate.Number)
		}
		windows = append(windows, &model.TmuxWindow{
			ID:     candidate.ID,
			Name:   name,
			Path:   path,
			Index:  candidate.Number,
			Active: candidate.ID == target.ActiveTabID || candidate.Focused,
		})
	}
	sort.SliceStable(windows, func(i, j int) bool { return windows[i].Index < windows[j].Index })
	return windows, nil
}

func (h *Herdr) runInPaneOrRollback(resource, resourceID, paneID, command string) error {
	if command == "" {
		return nil
	}
	if _, err := h.command("pane", "run", paneID, command); err != nil {
		runErr := fmt.Errorf("couldn't run Herdr startup command: %w", err)
		if _, rollbackErr := h.command(resource, "close", resourceID); rollbackErr != nil {
			return errors.Join(runErr, fmt.Errorf("couldn't roll back Herdr %s %q: %w", resource, resourceID, rollbackErr))
		}
		return runErr
	}
	return nil
}

func (h *Herdr) NewSession(sessionName, startDir, shellCommand string) (string, error) {
	args := []string{"workspace", "create", "--label", sessionName, "--no-focus"}
	if startDir != "" {
		args = append(args, "--cwd", startDir)
	}
	output, err := h.command(args...)
	if err != nil {
		return "", err
	}
	result, err := decodeResponse[createWorkspaceResult](output)
	if err != nil {
		return "", err
	}
	if err := h.runInPaneOrRollback("workspace", result.Workspace.ID, result.RootPane.ID, shellCommand); err != nil {
		return "", err
	}
	return result.Workspace.ID, nil
}

func (h *Herdr) NewWindow(startDir, name, shellCommand string) (string, error) {
	return h.NewWindowInSession(name, startDir, "", shellCommand)
}

func (h *Herdr) NewWindowInSession(name, startDir, targetSession, shellCommand string) (string, error) {
	state, err := h.snapshot()
	if err != nil {
		return "", err
	}
	target, err := h.resolveWorkspace(state, targetSession)
	if err != nil {
		return "", err
	}

	args := []string{"tab", "create", "--workspace", target.ID, "--label", name, "--no-focus"}
	if startDir != "" {
		args = append(args, "--cwd", startDir)
	}
	output, err := h.command(args...)
	if err != nil {
		return "", err
	}
	result, err := decodeResponse[createTabResult](output)
	if err != nil {
		return "", err
	}
	if err := h.runInPaneOrRollback("tab", result.Tab.ID, result.RootPane.ID, shellCommand); err != nil {
		return "", err
	}
	return result.Tab.ID, nil
}

func (h *Herdr) AttachSession(targetSession string) (string, error) {
	if _, err := h.SwitchClient(targetSession); err != nil {
		return "", err
	}
	if _, err := h.shell.CmdWithOutput(h.bin, h.scopedArgs()...); err != nil {
		return "", err
	}
	return fmt.Sprintf("attached to Herdr workspace: %s", targetSession), nil
}

func (h *Herdr) SwitchClient(targetSession string) (string, error) {
	state, err := h.snapshot()
	if err != nil {
		return "", err
	}
	target, err := h.resolveWorkspace(state, targetSession)
	if err != nil {
		return "", err
	}
	if _, err := h.command("workspace", "focus", target.ID); err != nil {
		return "", err
	}
	return fmt.Sprintf("focused Herdr workspace: %s", workspaceName(state, target)), nil
}

func (h *Herdr) SwitchOrAttach(targetSession string, opts model.ConnectOpts) (string, error) {
	if opts.Switch || h.IsAttached() {
		return h.SwitchClient(targetSession)
	}
	return h.AttachSession(targetSession)
}

func (h *Herdr) SendKeys(targetSession, command string) (string, error) {
	state, err := h.snapshot()
	if err != nil {
		return "", err
	}
	target, err := h.resolveWorkspace(state, targetSession)
	if err != nil {
		return "", err
	}
	p, err := paneForWorkspace(state, target)
	if err != nil {
		return "", err
	}
	return h.command("pane", "run", p.ID, command)
}

func (h *Herdr) CapturePane(targetSession string) (string, error) {
	state, err := h.snapshot()
	if err != nil {
		return "", err
	}
	target, err := h.resolveWorkspace(state, targetSession)
	if err != nil {
		return "", err
	}
	p, err := paneForWorkspace(state, target)
	if err != nil {
		return "", err
	}
	return h.command("pane", "read", p.ID, "--source", "visible", "--format", "ansi")
}

func (h *Herdr) NextWindow() (string, error) {
	state, err := h.snapshot()
	if err != nil {
		return "", err
	}
	target, err := h.resolveWorkspace(state, "")
	if err != nil {
		return "", err
	}
	windows, err := h.ListWindows(target.ID)
	if err != nil {
		return "", err
	}
	if len(windows) == 0 {
		return "", fmt.Errorf("Herdr workspace %q has no tabs", workspaceName(state, target))
	}
	next := 0
	for i, candidate := range windows {
		if candidate.Active {
			next = (i + 1) % len(windows)
			break
		}
	}
	return h.SelectWindow(windows[next].ID)
}

func (h *Herdr) SelectWindow(targetWindow string) (string, error) {
	state, err := h.snapshot()
	if err != nil {
		return "", err
	}
	for _, candidate := range state.Tabs {
		if candidate.ID == targetWindow {
			_, err := h.command("tab", "focus", candidate.ID)
			return candidate.ID, err
		}
	}

	if strings.HasSuffix(targetWindow, ":^") {
		workspaceTarget := strings.TrimSuffix(targetWindow, ":^")
		target, err := h.resolveWorkspace(state, workspaceTarget)
		if err != nil {
			return "", err
		}
		for _, candidate := range state.Tabs {
			if candidate.WorkspaceID == target.ID && candidate.Number == 1 {
				_, err := h.command("tab", "focus", candidate.ID)
				return candidate.ID, err
			}
		}
	}

	target, err := h.resolveWorkspace(state, "")
	if err != nil {
		return "", err
	}
	matches := make([]tab, 0, 1)
	for _, candidate := range state.Tabs {
		if candidate.WorkspaceID == target.ID && candidate.Label == targetWindow {
			matches = append(matches, candidate)
		}
	}
	if len(matches) != 1 {
		return "", fmt.Errorf("Herdr tab %q is not uniquely selectable", targetWindow)
	}
	_, err = h.command("tab", "focus", matches[0].ID)
	return matches[0].ID, err
}

func (h *Herdr) ListTmuxPanes() ([]*model.TmuxPane, error) {
	state, err := h.snapshot()
	if err != nil {
		return nil, err
	}
	target, err := h.resolveWorkspace(state, "")
	if err != nil {
		return nil, err
	}

	tabByID := make(map[string]tab, len(state.Tabs))
	for _, candidate := range state.Tabs {
		if candidate.WorkspaceID == target.ID {
			tabByID[candidate.ID] = candidate
		}
	}
	paneIndexes := make(map[string]int)
	panes := make([]*model.TmuxPane, 0, target.PaneCount)
	for _, candidate := range state.Panes {
		tabInfo, ok := tabByID[candidate.TabID]
		if !ok {
			continue
		}
		paneIndex := paneIndexes[candidate.TabID]
		paneIndexes[candidate.TabID] = paneIndex + 1
		command := candidate.Agent
		if command == "" {
			command = candidate.Title
		}
		tabName := tabInfo.Label
		if tabName == "" {
			tabName = fmt.Sprintf("tab-%d", tabInfo.Number)
		}
		panes = append(panes, &model.TmuxPane{
			WindowIndex: tabInfo.Number,
			WindowName:  tabName,
			PaneIndex:   paneIndex,
			PaneTitle:   candidate.Title,
			PaneCommand: command,
			PanePath:    panePath(candidate),
			PaneID:      candidate.ID,
		})
	}
	return panes, nil
}

func (h *Herdr) SelectPane(windowIndex, paneIndex int) (string, error) {
	return "", fmt.Errorf("Herdr %s does not expose direct pane focus through its CLI", strconv.Itoa(windowIndex)+"."+strconv.Itoa(paneIndex))
}

func (h *Herdr) GetCurrentSession() (string, error) {
	state, err := h.snapshot()
	if err != nil {
		return "", err
	}
	target, err := h.resolveWorkspace(state, "")
	if err != nil {
		return "", err
	}
	return workspaceName(state, target), nil
}
