package connector

import (
	"errors"
	"testing"

	"github.com/notes/zesh/dir"
	"github.com/notes/zesh/home"
	"github.com/notes/zesh/lister"
	"github.com/notes/zesh/model"
	"github.com/notes/zesh/namer"
	"github.com/notes/zesh/startup"
	"github.com/notes/zesh/tmux"
	"github.com/notes/zesh/tmuxinator"
	"github.com/notes/zesh/zoxide"
	"github.com/stretchr/testify/assert"
)

func TestConnectToTmuxReturnsStartupExecError(t *testing.T) {
	mockStartup := new(startup.MockStartup)
	mockTmux := new(tmux.MockTmux)

	c := &RealConnector{
		config:     model.Config{},
		dir:        new(dir.MockDir),
		home:       new(home.MockHome),
		lister:     new(lister.MockLister),
		namer:      new(namer.MockNamer),
		startup:    mockStartup,
		tmux:       mockTmux,
		zoxide:     new(zoxide.MockZoxide),
		tmuxinator: new(tmuxinator.MockTmuxinator),
	}

	connection := model.Connection{
		New: true,
		Session: model.SeshSession{
			Name: "demo",
			Path: "/tmp",
		},
	}

	mockStartup.On("ResolveCommand", connection.Session).Return("", nil)
	mockStartup.On("WrapForShell", "").Return("")
	mockTmux.On("NewSession", "demo", "/tmp", "").Return("", nil)
	mockStartup.On("Exec", connection.Session).Return("", errors.New("boom"))

	msg, err := connectToTmux(c, connection, model.ConnectOpts{})
	assert.Equal(t, "", msg)
	assert.EqualError(t, err, "boom")
	mockTmux.AssertNotCalled(t, "SwitchOrAttach", "demo", model.ConnectOpts{})
}
