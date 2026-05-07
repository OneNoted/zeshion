package connector

import (
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
