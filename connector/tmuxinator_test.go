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
	"github.com/stretchr/testify/mock"
)

func newConnectorWithMocks() (*RealConnector, *tmuxinator.MockTmuxinator, *tmux.MockTmux) {
	mockTmuxinator := new(tmuxinator.MockTmuxinator)
	mockTmux := new(tmux.MockTmux)
	c := &RealConnector{
		model.Config{},
		new(dir.MockDir),
		new(home.MockHome),
		new(lister.MockLister),
		new(namer.MockNamer),
		new(startup.MockStartup),
		mockTmux,
		new(zoxide.MockZoxide),
		mockTmuxinator,
	}
	return c, mockTmuxinator, mockTmux
}

func TestConnectToTmuxinator(t *testing.T) {
	t.Run("propagates error from Start and does not attach", func(t *testing.T) {
		c, mockTmuxinator, mockTmux := newConnectorWithMocks()
		mockTmuxinator.EXPECT().
			Start("sys").
			Return("", errors.New("boom"))

		connection := model.Connection{
			Found:   true,
			Session: model.SeshSession{Src: "tmuxinator", Name: "sys"},
			New:     true,
		}
		_, err := connectToTmuxinator(c, connection, model.ConnectOpts{})

		assert.ErrorContains(t, err, "failed to start tmuxinator session")
		assert.ErrorContains(t, err, "boom")
		mockTmux.AssertNotCalled(t, "AttachSession", mock.Anything)
		mockTmux.AssertNotCalled(t, "SwitchClient", mock.Anything)
	})

	t.Run("attaches via SwitchOrAttach after successful Start", func(t *testing.T) {
		c, mockTmuxinator, mockTmux := newConnectorWithMocks()
		opts := model.ConnectOpts{}
		mockTmuxinator.EXPECT().Start("sys").Return("", nil)
		mockTmux.EXPECT().SwitchOrAttach("sys", opts).Return("attaching to tmux session: sys", nil)

		connection := model.Connection{
			Found:   true,
			Session: model.SeshSession{Src: "tmuxinator", Name: "sys"},
			New:     true,
		}
		msg, err := connectToTmuxinator(c, connection, opts)

		assert.Nil(t, err)
		assert.Equal(t, "attaching to tmux session: sys", msg)
	})
}
