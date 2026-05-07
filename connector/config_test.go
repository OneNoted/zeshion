package connector

import (
	"testing"

	"github.com/notes/zeshion/dir"
	"github.com/notes/zeshion/home"
	"github.com/notes/zeshion/lister"
	"github.com/notes/zeshion/model"
	"github.com/notes/zeshion/namer"
	"github.com/notes/zeshion/startup"
	"github.com/notes/zeshion/tmux"
	"github.com/notes/zeshion/tmuxinator"
	"github.com/notes/zeshion/zoxide"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestConfigStrategy(t *testing.T) {
	mockDir := new(dir.MockDir)
	mockHome := new(home.MockHome)
	mockLister := new(lister.MockLister)
	mockNamer := new(namer.MockNamer)
	mockStartup := new(startup.MockStartup)
	mockTmux := new(tmux.MockTmux)
	mockZoxide := new(zoxide.MockZoxide)
	mockTmuxinator := new(tmuxinator.MockTmuxinator)

	c := testConnector(mockDir, mockHome, mockLister, mockNamer, mockStartup, mockTmux, mockZoxide, mockTmuxinator)
	mockTmux.On("AttachSession", mock.Anything).Return("attaching", nil)
	mockZoxide.On("Add", mock.Anything).Return(nil)

	t.Run("should create and attach to config session", func(t *testing.T) {
		mockTmux.On("IsAttached").Return(false)
		mockLister.On("FindConfigSession", "tmux config").Return(model.SeshSession{
			Name: "tmux config",
			Path: "/Users/joshmedeski/c/dotfiles/.config/tmux",
		}, true)
		connection, err := configStrategy(c, "tmux config")
		assert.Nil(t, err)
		assert.Equal(t, "tmux config", connection.Session.Name)
	})
}
