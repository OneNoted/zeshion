package startup

import (
	"testing"

	"github.com/OneNoted/zeshion/home"
	"github.com/OneNoted/zeshion/lister"
	"github.com/OneNoted/zeshion/model"
	"github.com/OneNoted/zeshion/oswrap"
	"github.com/OneNoted/zeshion/replacer"
	"github.com/OneNoted/zeshion/tmux"
	"github.com/stretchr/testify/assert"
)

func TestExecCreatesWindowsInTargetSession(t *testing.T) {
	mockOs := new(oswrap.MockOs)
	mockLister := new(lister.MockLister)
	mockTmux := new(tmux.MockTmux)
	mockHome := new(home.MockHome)
	mockReplacer := new(replacer.MockReplacer)

	s := &RealStartup{
		os:       mockOs,
		lister:   mockLister,
		tmux:     mockTmux,
		config:   model.Config{WindowConfigs: []model.WindowConfig{{Name: "editor", StartupScript: "echo hi"}}},
		home:     mockHome,
		replacer: mockReplacer,
	}
	session := model.SeshSession{ID: "$1", Name: "demo", Path: "/tmp", WindowNames: []string{"editor"}}

	mockHome.On("ExpandPath", "/tmp").Return("/tmp", nil)
	mockOs.On("Getenv", "SHELL").Return("/bin/zsh")
	mockTmux.On("NewWindowInSession", "editor", "/tmp", "$1", `'/bin/zsh' -i -c 'echo hi'`).Return("", nil)
	mockTmux.On("SelectWindow", "$1:^").Return("", nil)
	mockLister.On("FindConfigSession", "demo").Return(model.SeshSession{}, false)
	mockLister.On("FindConfigWildcard", "/tmp").Return(model.WildcardConfig{}, false)

	msg, err := s.Exec(session)
	assert.Nil(t, err)
	assert.Equal(t, "", msg)
	mockTmux.AssertNotCalled(t, "NewWindow", "/tmp", "editor", `'/bin/zsh' -i -c 'echo hi'`)
}
