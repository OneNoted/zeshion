package connector

import (
	"errors"
	"testing"

	"github.com/OneNoted/zeshion/dir"
	"github.com/OneNoted/zeshion/home"
	"github.com/OneNoted/zeshion/lister"
	"github.com/OneNoted/zeshion/model"
	"github.com/OneNoted/zeshion/namer"
	"github.com/OneNoted/zeshion/startup"
	"github.com/OneNoted/zeshion/tmux"
	"github.com/OneNoted/zeshion/tmuxinator"
	"github.com/OneNoted/zeshion/zoxide"
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
	mockTmux.On("NewSession", "demo", "/tmp", "").Return("$1", nil)
	startupSession := connection.Session
	startupSession.ID = "$1"
	mockStartup.On("Exec", startupSession).Return("", errors.New("boom"))

	msg, err := connectToTmux(c, connection, model.ConnectOpts{})
	assert.Equal(t, "", msg)
	assert.EqualError(t, err, "boom")
	mockTmux.AssertNotCalled(t, "SwitchOrAttach")
}
