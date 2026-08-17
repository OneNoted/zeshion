package connector

import (
	"testing"

	"github.com/OneNoted/zeshion/lister"
	"github.com/OneNoted/zeshion/model"
	"github.com/OneNoted/zeshion/tmux"
	"github.com/OneNoted/zeshion/tmuxinator"
	"github.com/OneNoted/zeshion/zoxide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectSessionRoutesHerdrByWorkspaceID(t *testing.T) {
	backend := tmux.NewMockTmux(t)
	backend.On("SwitchOrAttach", "w2", model.ConnectOpts{}).Return("focused", nil).Once()
	sessions := lister.NewMockLister(t)
	selected := model.SeshSession{ID: "w2", Src: "herdr", Name: "repo"}
	sessions.On("FindHerdrSession", "w2").Return(selected, true, nil).Once()
	c := &RealConnector{lister: sessions, herdr: backend}

	result, err := c.ConnectSession(selected, model.ConnectOpts{})
	require.NoError(t, err)
	assert.Equal(t, "focused", result)
}

func TestConnectSessionRefreshesStaleTmuxID(t *testing.T) {
	backend := tmux.NewMockTmux(t)
	backend.On("SwitchOrAttach", "$7", model.ConnectOpts{}).Return("attached", nil).Once()
	sessions := lister.NewMockLister(t)
	sessions.On("FindTmuxSession", "repo").Return(model.SeshSession{
		ID:   "$7",
		Src:  "tmux",
		Name: "repo",
	}, true).Once()
	c := &RealConnector{lister: sessions, tmux: backend}

	result, err := c.ConnectSession(
		model.SeshSession{ID: "$1", Src: "tmux", Name: "repo"},
		model.ConnectOpts{},
	)

	require.NoError(t, err)
	assert.Equal(t, "attached", result)
}

func TestConnectSessionRecoversFromReassignedHerdrID(t *testing.T) {
	backend := tmux.NewMockTmux(t)
	backend.On("SwitchOrAttach", "w7", model.ConnectOpts{}).Return("focused", nil).Once()
	sessions := lister.NewMockLister(t)
	sessions.On("FindHerdrSession", "w1").Return(model.SeshSession{
		ID:   "w1",
		Src:  "herdr",
		Name: "other",
	}, true, nil).Once()
	sessions.On("FindHerdrSession", "repo").Return(model.SeshSession{
		ID:   "w7",
		Src:  "herdr",
		Name: "repo",
	}, true, nil).Once()
	c := &RealConnector{lister: sessions, herdr: backend}

	result, err := c.ConnectSession(
		model.SeshSession{ID: "w1", Src: "herdr", Name: "repo"},
		model.ConnectOpts{},
	)

	require.NoError(t, err)
	assert.Equal(t, "focused", result)
}
func TestConnectSessionKeepsHerdrIDWhenPanePathChanges(t *testing.T) {
	backend := tmux.NewMockTmux(t)
	backend.On("SwitchOrAttach", "w2", model.ConnectOpts{}).Return("focused", nil).Once()
	sessions := lister.NewMockLister(t)
	live := model.SeshSession{ID: "w2", Src: "herdr", Name: "repo", Path: "/work/new"}
	sessions.On("FindHerdrSession", "w2").Return(live, true, nil).Once()
	zoxideStore := zoxide.NewMockZoxide(t)
	zoxideStore.On("Add", "/work/new").Return(nil).Once()
	c := &RealConnector{lister: sessions, herdr: backend, zoxide: zoxideStore}

	result, err := c.ConnectSession(
		model.SeshSession{ID: "w2", Src: "herdr", Name: "repo", Path: "/work/old"},
		model.ConnectOpts{},
	)

	require.NoError(t, err)
	assert.Equal(t, "focused", result)
}

func TestConnectSessionDoesNotAddPanePathToZoxide(t *testing.T) {
	backend := tmux.NewMockTmux(t)
	backend.On("ListTmuxPanes").Return([]*model.TmuxPane{{
		WindowIndex: 2,
		PaneIndex:   3,
		PaneID:      "w1:p3",
	}}, nil).Once()
	backend.On("SelectPane", 2, 3).Return("selected", nil).Once()
	zoxideStore := zoxide.NewMockZoxide(t)
	c := &RealConnector{mux: backend, zoxide: zoxideStore}

	result, err := c.ConnectSession(model.SeshSession{
		ID:   "w1:p3",
		Src:  "herdr-pane",
		Name: "shell/zsh",
		Path: "/work/repo",
	}, model.ConnectOpts{})

	require.NoError(t, err)
	assert.Equal(t, "selected", result)
	zoxideStore.AssertNotCalled(t, "Add")
}

func TestConnectSessionResolvesCreationRecipeAgainstLiveSessions(t *testing.T) {
	backend := tmux.NewMockTmux(t)
	backend.On("IsAttached").Return(false).Once()
	backend.On("SwitchOrAttach", "$1", model.ConnectOpts{}).Return("attached", nil).Once()
	sessions := lister.NewMockLister(t)
	sessions.On("FindTmuxSession", "repo").Return(model.SeshSession{
		ID:   "$1",
		Src:  "tmux",
		Name: "repo",
	}, true).Once()

	c := &RealConnector{
		lister:  sessions,
		tmux:    backend,
		zellij:  backend,
		mux:     backend,
		muxName: "tmux",
	}
	result, err := c.ConnectSession(model.SeshSession{Src: "config", Name: "repo"}, model.ConnectOpts{})

	require.NoError(t, err)
	assert.Equal(t, "attached", result)
}
func TestConnectSessionResolvesTmuxinatorRecipeAgainstLiveSession(t *testing.T) {
	backend := tmux.NewMockTmux(t)
	backend.On("IsAttached").Return(false).Once()
	backend.On("SwitchOrAttach", "$1", model.ConnectOpts{}).Return("attached", nil).Once()
	sessions := lister.NewMockLister(t)
	sessions.On("FindTmuxSession", "repo").Return(model.SeshSession{
		ID:   "$1",
		Src:  "tmux",
		Name: "repo",
	}, true).Once()
	tmuxinatorRunner := tmuxinator.NewMockTmuxinator(t)

	c := &RealConnector{
		lister:     sessions,
		tmux:       backend,
		zellij:     backend,
		mux:        backend,
		muxName:    "tmux",
		tmuxinator: tmuxinatorRunner,
	}
	result, err := c.ConnectSession(model.SeshSession{Src: "tmuxinator", Name: "repo"}, model.ConnectOpts{})

	require.NoError(t, err)
	assert.Equal(t, "attached", result)
	tmuxinatorRunner.AssertNotCalled(t, "Start")
}

func TestConnectPrioritizesSelectedHerdrBackend(t *testing.T) {
	backend := tmux.NewMockTmux(t)
	backend.On("IsAttached").Return(false).Once()
	backend.On("SwitchOrAttach", "w2", model.ConnectOpts{}).Return("focused", nil).Once()
	sessions := lister.NewMockLister(t)
	sessions.On("FindHerdrSession", "w2").Return(model.SeshSession{
		ID:   "w2",
		Src:  "herdr",
		Name: "repo",
	}, true, nil).Once()

	c := &RealConnector{
		lister:  sessions,
		herdr:   backend,
		mux:     backend,
		muxName: "herdr",
	}
	result, err := c.Connect("w2", model.ConnectOpts{})

	require.NoError(t, err)
	assert.Equal(t, "focused", result)
}
