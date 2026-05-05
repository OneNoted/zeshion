package zellij

import (
	"testing"

	"github.com/notes/zesh/model"
	"github.com/notes/zesh/oswrap"
	"github.com/notes/zesh/shell"
	"github.com/stretchr/testify/assert"
)

func TestListSessions(t *testing.T) {
	mockOs := new(oswrap.MockOs)
	mockShell := new(shell.MockShell)
	z := NewZellij(mockOs, mockShell, "zellij")

	mockShell.EXPECT().ListCmd("zellij", "list-sessions").Return([]string{
		"\x1b[32;1mwork\x1b[m [Created \x1b[35;1m40m 22s\x1b[m ago] (current)",
		"\x1b[32;1mold\x1b[m [Created \x1b[35;1m1h\x1b[m ago] (\x1b[31;1mEXITED\x1b[m - attach to resurrect)",
		"",
	}, nil)

	sessions, err := z.ListSessions()

	assert.NoError(t, err)
	assert.Len(t, sessions, 2)
	assert.Equal(t, "work", sessions[0].Name)
	assert.Equal(t, "old", sessions[1].Name)
}

func TestNewSessionUsesBackgroundCreateAndCwd(t *testing.T) {
	mockOs := new(oswrap.MockOs)
	mockShell := new(shell.MockShell)
	z := NewZellij(mockOs, mockShell, "zellij")

	mockShell.EXPECT().
		Cmd("zellij", "attach", "--create-background", "work", "options", "--default-cwd", "/tmp/work").
		Return("", nil)

	out, err := z.NewSession("work", "/tmp/work", "")

	assert.NoError(t, err)
	assert.Equal(t, "", out)
}

func TestListWindowsUsesTargetSession(t *testing.T) {
	mockOs := new(oswrap.MockOs)
	mockShell := new(shell.MockShell)
	z := NewZellij(mockOs, mockShell, "zellij")

	mockShell.EXPECT().
		Cmd("zellij", "--session", "work", "action", "list-tabs", "--json").
		Return(`[{"position":0,"name":"editor","active":true,"tab_id":1}]`, nil)

	windows, err := z.ListWindows("work")

	assert.NoError(t, err)
	assert.Len(t, windows, 1)
	assert.Equal(t, "editor", windows[0].Name)
	assert.True(t, windows[0].Active)
}

func TestSwitchOrAttachSwitchesWhenAttached(t *testing.T) {
	mockOs := new(oswrap.MockOs)
	mockShell := new(shell.MockShell)
	z := NewZellij(mockOs, mockShell, "zellij")

	mockOs.EXPECT().Getenv("ZELLIJ").Return("1")
	mockShell.EXPECT().Cmd("zellij", "action", "switch-session", "work").Return("", nil)

	msg, err := z.SwitchOrAttach("work", model.ConnectOpts{})

	assert.NoError(t, err)
	assert.Equal(t, "switching to zellij session: work", msg)
}
