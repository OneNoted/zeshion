package herdr

import (
	"errors"
	"testing"

	"github.com/OneNoted/zeshion/oswrap"
	"github.com/OneNoted/zeshion/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const testSnapshot = `{"result":{"snapshot":{"focused_workspace_id":"w2","focused_tab_id":"w2:t1","focused_pane_id":"w2:p1","workspaces":[{"workspace_id":"w1","active_tab_id":"w1:t1","label":"repo","number":1,"pane_count":1,"tab_count":1,"focused":false,"worktree":{"checkout_path":"/work/repo-one"}},{"workspace_id":"w2","active_tab_id":"w2:t1","label":"repo","number":2,"pane_count":1,"tab_count":1,"focused":true}],"tabs":[{"tab_id":"w1:t1","workspace_id":"w1","label":"editor","number":1,"pane_count":1,"focused":false},{"tab_id":"w2:t1","workspace_id":"w2","label":"shell","number":1,"pane_count":1,"focused":true}],"panes":[{"pane_id":"w1:p1","workspace_id":"w1","tab_id":"w1:t1","cwd":"/work/repo-one","terminal_title_stripped":"nvim","focused":false},{"pane_id":"w2:p1","workspace_id":"w2","tab_id":"w2:t1","cwd":"/work/repo-two","terminal_title_stripped":"shell","focused":true}]}}}`

func newTestHerdr(t *testing.T) (*Herdr, *shell.MockShell) {
	t.Helper()
	mockOs := oswrap.NewMockOs(t)
	mockOs.On("Getenv", mock.Anything).Return("")
	mockShell := shell.NewMockShell(t)
	return NewHerdr(mockOs, mockShell, "herdr", ""), mockShell
}

func newTestHerdrInWorkspace(t *testing.T, workspaceID string) (*Herdr, *shell.MockShell) {
	t.Helper()
	mockOs := oswrap.NewMockOs(t)
	mockOs.On("Getenv", "HERDR_ENV").Return("1").Maybe()
	mockOs.On("Getenv", "HERDR_WORKSPACE_ID").Return(workspaceID).Maybe()
	mockShell := shell.NewMockShell(t)
	return NewHerdr(mockOs, mockShell, "herdr", ""), mockShell
}

func TestListSessionsPreservesWorkspaceIdentity(t *testing.T) {
	h, mockShell := newTestHerdr(t)
	mockShell.On("Cmd", "herdr", "api", "snapshot").Return(testSnapshot, nil).Once()

	sessions, err := h.ListSessions()
	require.NoError(t, err)
	require.Len(t, sessions, 2)
	assert.Equal(t, "w1", sessions[0].ID)
	assert.Equal(t, "repo", sessions[0].Name)
	assert.Equal(t, "/work/repo-one", sessions[0].Path)
	assert.Equal(t, "w2", sessions[1].ID)
	assert.Equal(t, 1, sessions[1].Attached)
	assert.Equal(t, "/work/repo-two", sessions[1].Path)
}

func TestListSessionsMarksCallerWorkspaceAttached(t *testing.T) {
	h, mockShell := newTestHerdrInWorkspace(t, "w1")
	mockShell.On("Cmd", "herdr", "api", "snapshot").Return(testSnapshot, nil).Once()

	sessions, err := h.ListSessions()

	require.NoError(t, err)
	assert.Equal(t, 1, sessions[0].Attached)
	assert.Equal(t, 0, sessions[1].Attached)
}

func TestListWindowsDefaultsToCallerWorkspace(t *testing.T) {
	h, mockShell := newTestHerdrInWorkspace(t, "w1")
	mockShell.On("Cmd", "herdr", "api", "snapshot").Return(testSnapshot, nil).Once()

	windows, err := h.ListWindows("")

	require.NoError(t, err)
	require.Len(t, windows, 1)
	assert.Equal(t, "w1:t1", windows[0].ID)
}

func TestSwitchClientRejectsAmbiguousWorkspaceLabel(t *testing.T) {
	h, mockShell := newTestHerdr(t)
	mockShell.On("Cmd", "herdr", "api", "snapshot").Return(testSnapshot, nil).Once()

	_, err := h.SwitchClient("repo")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ambiguous")
	assert.Contains(t, err.Error(), "w1, w2")
}

func TestSwitchClientUsesExactWorkspaceID(t *testing.T) {
	h, mockShell := newTestHerdr(t)
	mockShell.On("Cmd", "herdr", "api", "snapshot").Return(testSnapshot, nil).Once()
	mockShell.On("Cmd", "herdr", "workspace", "focus", "w1").Return("", nil).Once()

	result, err := h.SwitchClient("w1")
	require.NoError(t, err)
	assert.Equal(t, "focused Herdr workspace: repo", result)
}

func TestNewSessionRunsStartupCommandInCreatedPane(t *testing.T) {
	h, mockShell := newTestHerdr(t)
	created := `{"result":{"workspace":{"workspace_id":"w7","label":"repo"},"tab":{"tab_id":"w7:t1"},"root_pane":{"pane_id":"w7:p1"}}}`
	mockShell.On("Cmd", "herdr", "workspace", "create", "--label", "repo", "--no-focus", "--cwd", "/work/repo").Return(created, nil).Once()
	mockShell.On("Cmd", "herdr", "pane", "run", "w7:p1", "nvim").Return("", nil).Once()

	id, err := h.NewSession("repo", "/work/repo", "nvim")
	require.NoError(t, err)
	assert.Equal(t, "w7", id)
}

func TestNewSessionRollsBackWhenStartupDispatchFails(t *testing.T) {
	h, mockShell := newTestHerdr(t)
	created := `{"result":{"workspace":{"workspace_id":"w7","label":"repo"},"tab":{"tab_id":"w7:t1"},"root_pane":{"pane_id":"w7:p1"}}}`
	mockShell.On("Cmd", "herdr", "workspace", "create", "--label", "repo", "--no-focus", "--cwd", "/work/repo").Return(created, nil).Once()
	mockShell.On("Cmd", "herdr", "pane", "run", "w7:p1", "nvim").Return("", errors.New("dispatch failed")).Once()
	mockShell.On("Cmd", "herdr", "workspace", "close", "w7").Return("", nil).Once()

	_, err := h.NewSession("repo", "/work/repo", "nvim")
	require.Error(t, err)
	assert.ErrorContains(t, err, "couldn't run Herdr startup command")
}

func TestNewWindowRollsBackWhenStartupDispatchFails(t *testing.T) {
	h, mockShell := newTestHerdr(t)
	created := `{"result":{"tab":{"tab_id":"w1:t2"},"root_pane":{"pane_id":"w1:p2"}}}`
	mockShell.On("Cmd", "herdr", "api", "snapshot").Return(testSnapshot, nil).Once()
	mockShell.On("Cmd", "herdr", "tab", "create", "--workspace", "w1", "--label", "logs", "--no-focus", "--cwd", "/work/logs").Return(created, nil).Once()
	mockShell.On("Cmd", "herdr", "pane", "run", "w1:p2", "tail -f app.log").Return("", errors.New("dispatch failed")).Once()
	mockShell.On("Cmd", "herdr", "tab", "close", "w1:t2").Return("", nil).Once()

	_, err := h.NewWindowInSession("logs", "/work/logs", "w1", "tail -f app.log")
	require.Error(t, err)
	assert.ErrorContains(t, err, "couldn't run Herdr startup command")
}

func TestNamedSessionScopesCommandsOutsideHerdr(t *testing.T) {
	mockOs := oswrap.NewMockOs(t)
	mockOs.On("Getenv", mock.Anything).Return("")
	mockShell := shell.NewMockShell(t)
	h := NewHerdr(mockOs, mockShell, "herdr", "agents")
	mockShell.On("Cmd", "herdr", "--session", "agents", "api", "snapshot").Return(testSnapshot, nil).Once()

	_, err := h.ListWindows("w1")
	require.NoError(t, err)
}

func TestListSessionsTreatsUnavailableServerAsEmptySource(t *testing.T) {
	h, mockShell := newTestHerdr(t)
	mockShell.On("Cmd", "herdr", "api", "snapshot").Return("", errors.New("server unavailable")).Once()

	sessions, err := h.ListSessions()
	require.NoError(t, err)
	assert.Empty(t, sessions)
}

func TestPaneForWorkspacePrefersFocusedPane(t *testing.T) {
	state := snapshot{
		FocusedWorkspaceID: "w2",
		FocusedPaneID:      "w2:p2",
		Panes: []pane{
			{ID: "w2:p1", WorkspaceID: "w2", TabID: "w2:t1"},
			{ID: "w2:p2", WorkspaceID: "w2", TabID: "w2:t1", Focused: true},
		},
	}

	got, err := paneForWorkspace(state, workspace{ID: "w2", ActiveTabID: "w2:t1", Focused: true})
	require.NoError(t, err)
	assert.Equal(t, "w2:p2", got.ID)
}

func TestPaneForWorkspaceFallsBackToActiveTab(t *testing.T) {
	state := snapshot{
		Panes: []pane{
			{ID: "w2:p1", WorkspaceID: "w2", TabID: "w2:t1"},
			{ID: "w2:p2", WorkspaceID: "w2", TabID: "w2:t2"},
		},
	}

	got, err := paneForWorkspace(state, workspace{ID: "w2", ActiveTabID: "w2:t2"})
	require.NoError(t, err)
	assert.Equal(t, "w2:p2", got.ID)
}

func TestPanePathPrefersForegroundProcessDirectory(t *testing.T) {
	assert.Equal(t, "/work/current", panePath(pane{
		Cwd:           "/work",
		ForegroundCwd: "/work/current",
	}))
	assert.Equal(t, "/work", panePath(pane{Cwd: "/work"}))
}
