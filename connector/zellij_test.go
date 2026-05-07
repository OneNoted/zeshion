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
)

func TestConnectToZellijUsesZellijBackend(t *testing.T) {
	mockLister := new(lister.MockLister)
	mockTmux := new(tmux.MockTmux)
	mockZellij := new(tmux.MockTmux)
	mockZoxide := new(zoxide.MockZoxide)

	c := testConnector(new(dir.MockDir), new(home.MockHome), mockLister, new(namer.MockNamer), new(startup.MockStartup), mockTmux, mockZoxide, new(tmuxinator.MockTmuxinator))
	c.zellij = mockZellij
	c.mux = mockTmux
	c.muxName = "tmux"

	session := model.SeshSession{Src: "zellij", Name: "fascinating-jellyfish"}
	mockTmux.EXPECT().IsAttached().Return(false)
	mockLister.EXPECT().FindTmuxSession("fascinating-jellyfish").Return(model.SeshSession{}, false)
	mockLister.EXPECT().FindZellijSession("fascinating-jellyfish").Return(session, true)
	mockZoxide.EXPECT().Add("").Return(nil)
	mockZellij.EXPECT().SwitchOrAttach("fascinating-jellyfish", model.ConnectOpts{Switch: true}).Return("switching to zellij session: fascinating-jellyfish", nil)

	msg, err := c.Connect("fascinating-jellyfish", model.ConnectOpts{Switch: true})

	assert.NoError(t, err)
	assert.Equal(t, "switching to zellij session: fascinating-jellyfish", msg)
	mockTmux.AssertNotCalled(t, "SwitchOrAttach", "fascinating-jellyfish", model.ConnectOpts{Switch: true})
}
