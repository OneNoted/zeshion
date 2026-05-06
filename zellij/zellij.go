package zellij

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/notes/zesh/model"
	"github.com/notes/zesh/oswrap"
	"github.com/notes/zesh/shell"
)

var ansiEscapePattern = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

var terminalBoundaryProcesses = map[string]bool{
	"Alacritty":             true,
	"com.mitchellh.ghostty": true,
	"contour":               true,
	"foot":                  true,
	"footclient":            true,
	"ghostty":               true,
	"gnome-terminal-server": true,
	"kitty":                 true,
	"konsole":               true,
	"rio":                   true,
	"tabby":                 true,
	"tilix":                 true,
	"wezterm":               true,
	"wezterm-gui":           true,
	"xfce4-terminal":        true,
	"xterm":                 true,
}

type Zellij struct {
	os    oswrap.Os
	shell shell.Shell
	bin   string
}

func NewZellij(os oswrap.Os, shell shell.Shell, bin string) *Zellij {
	if bin == "" {
		bin = "zellij"
	}
	return &Zellij{os: os, shell: shell, bin: bin}
}

func (z *Zellij) ListSessions() ([]*model.TmuxSession, error) {
	output, err := z.shell.ListCmd(z.bin, "list-sessions")
	if err != nil {
		return []*model.TmuxSession{}, nil
	}

	sessions := make([]*model.TmuxSession, 0, len(output))
	for _, line := range output {
		name := parseSessionName(line)
		if name == "" {
			continue
		}
		sessions = append(sessions, &model.TmuxSession{Name: name})
	}
	return sessions, nil
}

func parseSessionName(line string) string {
	name := strings.TrimSpace(ansiEscapePattern.ReplaceAllString(line, ""))
	if createdIndex := strings.Index(name, " [Created "); createdIndex >= 0 {
		name = name[:createdIndex]
	}
	if statusIndex := strings.Index(name, " ("); statusIndex >= 0 {
		name = name[:statusIndex]
	}
	return strings.TrimSpace(name)
}

func (z *Zellij) ListWindows(targetSession string) ([]*model.TmuxWindow, error) {
	args := z.sessionArgs(targetSession, "action", "list-tabs", "--json")
	output, err := z.shell.Cmd(z.bin, args...)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(output) == "" {
		return []*model.TmuxWindow{}, nil
	}

	var tabs []struct {
		Position  int    `json:"position"`
		Name      string `json:"name"`
		Active    bool   `json:"active"`
		TabID     int    `json:"tab_id"`
		PaneCount int    `json:"pane_count"`
	}
	if err := json.Unmarshal([]byte(output), &tabs); err != nil {
		return nil, err
	}

	windows := make([]*model.TmuxWindow, 0, len(tabs))
	for _, tab := range tabs {
		index := tab.Position
		if index == 0 && tab.TabID != 0 {
			index = tab.TabID
		}
		windows = append(windows, &model.TmuxWindow{
			Index:  index,
			Name:   tab.Name,
			Active: tab.Active,
		})
	}
	return windows, nil
}

func (z *Zellij) NewSession(sessionName string, startDir string, shellCommand string) (string, error) {
	args := []string{"attach", "--create-background", sessionName}
	if startDir != "" {
		args = append(args, "options", "--default-cwd", startDir)
	}
	out, err := z.shell.Cmd(z.bin, args...)
	if err != nil {
		return "", err
	}
	if shellCommand != "" {
		if _, err := z.SendKeys(sessionName, shellCommand); err != nil {
			return "", err
		}
	}
	return out, nil
}

func (z *Zellij) NewWindow(startDir string, name string, shellCommand string) (string, error) {
	return z.NewWindowInSession(name, startDir, "", shellCommand)
}

func (z *Zellij) NewWindowInSession(name string, startDir string, targetSession string, shellCommand string) (string, error) {
	args := z.sessionArgs(targetSession, "action", "new-tab", "--name", name)
	if startDir != "" {
		args = append(args, "--cwd", startDir)
	}
	out, err := z.shell.Cmd(z.bin, args...)
	if err != nil {
		return "", err
	}
	if shellCommand != "" {
		if _, err := z.SendKeys(targetSession, shellCommand); err != nil {
			return "", err
		}
	}
	return out, nil
}

func (z *Zellij) IsAttached() bool {
	if z.os.Getenv("ZELLIJ") == "" && z.os.Getenv("ZELLIJ_SESSION_NAME") == "" {
		return false
	}
	return hasZellijAncestor(z.os)
}

func hasZellijAncestor(os oswrap.Os) bool {
	pid := os.Getpid()
	for range 64 {
		if pid <= 1 {
			return false
		}

		comm, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
		if err != nil {
			// Non-Linux fallback: if Zellij exported its environment but /proc is
			// unavailable, preserve the previous env-based behavior.
			return true
		}
		process := strings.TrimSpace(string(comm))
		if process == "zellij" {
			return true
		}
		if terminalBoundaryProcesses[process] {
			return false
		}

		parent, ok := parentPID(os, pid)
		if !ok || parent == pid {
			return false
		}
		pid = parent
	}
	return false
}

func parentPID(os oswrap.Os, pid int) (int, bool) {
	stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return 0, false
	}
	line := string(stat)
	endComm := strings.LastIndex(line, ")")
	if endComm < 0 || endComm+2 >= len(line) {
		return 0, false
	}
	fields := strings.Fields(line[endComm+1:])
	if len(fields) < 2 {
		return 0, false
	}
	parent, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, false
	}
	return parent, true
}

func (z *Zellij) AttachSession(targetSession string) (string, error) {
	return z.shell.Cmd(z.bin, "attach", targetSession)
}

func (z *Zellij) SendKeys(targetSession string, command string) (string, error) {
	args := z.sessionArgs(targetSession, "action", "paste", command)
	if _, err := z.shell.Cmd(z.bin, args...); err != nil {
		return "", err
	}
	args = z.sessionArgs(targetSession, "action", "send-keys", "Enter")
	return z.shell.Cmd(z.bin, args...)
}

func (z *Zellij) SwitchClient(targetSession string) (string, error) {
	args := []string{"action", "switch-session", targetSession}
	return z.shell.Cmd(z.bin, args...)
}

func (z *Zellij) CapturePane(targetSession string) (string, error) {
	args := z.sessionArgs(targetSession, "action", "dump-screen", "--full", "--ansi")
	return z.shell.Cmd(z.bin, args...)
}

func (z *Zellij) NextWindow() (string, error) {
	return z.shell.Cmd(z.bin, "action", "go-to-next-tab")
}

func (z *Zellij) SelectWindow(targetWindow string) (string, error) {
	targetSession, tabName := splitTarget(targetWindow)
	args := z.sessionArgs(targetSession, "action", "go-to-tab-name", tabName)
	return z.shell.Cmd(z.bin, args...)
}

func (z *Zellij) SwitchOrAttach(name string, opts model.ConnectOpts) (string, error) {
	if z.IsAttached() {
		if _, err := z.SwitchClient(name); err != nil {
			return "", fmt.Errorf("failed to switch to zellij session: %w", err)
		}
		return fmt.Sprintf("switching to zellij session: %s", name), nil
	}
	if _, err := z.AttachSession(name); err != nil {
		return "", fmt.Errorf("failed to attach to zellij session: %w", err)
	}
	return fmt.Sprintf("attaching to zellij session: %s", name), nil
}

func (z *Zellij) ListTmuxPanes() ([]*model.TmuxPane, error) {
	output, err := z.shell.Cmd(z.bin, "action", "list-panes", "--json")
	if err != nil {
		return []*model.TmuxPane{}, nil
	}
	if strings.TrimSpace(output) == "" {
		return []*model.TmuxPane{}, nil
	}

	var panes []struct {
		ID          any    `json:"id"`
		Title       string `json:"title"`
		PaneCommand string `json:"pane_command"`
		PaneCwd     string `json:"pane_cwd"`
		TabID       int    `json:"tab_id"`
		TabName     string `json:"tab_name"`
	}
	if err := json.Unmarshal([]byte(output), &panes); err != nil {
		return nil, err
	}

	result := make([]*model.TmuxPane, 0, len(panes))
	for _, pane := range panes {
		result = append(result, &model.TmuxPane{
			WindowIndex: pane.TabID,
			WindowName:  pane.TabName,
			PaneIndex:   paneNumber(pane.ID),
			PaneTitle:   pane.Title,
			PaneCommand: pane.PaneCommand,
			PanePath:    pane.PaneCwd,
			PaneID:      paneID(pane.ID),
		})
	}
	return result, nil
}

func (z *Zellij) SelectPane(windowIndex int, paneIndex int) (string, error) {
	pane := fmt.Sprintf("terminal_%d", paneIndex)
	if _, err := z.shell.Cmd(z.bin, "action", "go-to-tab", fmt.Sprintf("%d", windowIndex)); err != nil {
		return "", err
	}
	if _, err := z.shell.Cmd(z.bin, "action", "focus-pane-id", pane); err != nil {
		return "", err
	}
	return fmt.Sprintf("selected pane %d in tab %d", paneIndex, windowIndex), nil
}

func (z *Zellij) GetCurrentSession() (string, error) {
	if session := z.os.Getenv("ZELLIJ_SESSION_NAME"); session != "" {
		return session, nil
	}
	return "", fmt.Errorf("not inside a zellij session")
}

func (z *Zellij) sessionArgs(targetSession string, args ...string) []string {
	if targetSession == "" {
		return args
	}
	return append([]string{"--session", targetSession}, args...)
}

func splitTarget(target string) (string, string) {
	parts := strings.SplitN(target, ":", 2)
	if len(parts) != 2 {
		return "", target
	}
	return parts[0], parts[1]
}

func paneID(id any) string {
	switch v := id.(type) {
	case string:
		if strings.HasPrefix(v, "terminal_") || strings.HasPrefix(v, "plugin_") {
			return v
		}
		return "terminal_" + v
	case float64:
		return fmt.Sprintf("terminal_%d", int(v))
	default:
		return fmt.Sprintf("%v", v)
	}
}

func paneNumber(id any) int {
	switch v := id.(type) {
	case string:
		v = strings.TrimPrefix(v, "terminal_")
		v = strings.TrimPrefix(v, "plugin_")
		var n int
		fmt.Sscanf(v, "%d", &n)
		return n
	case float64:
		return int(v)
	default:
		return 0
	}
}
