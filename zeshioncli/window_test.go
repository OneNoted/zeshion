package zeshioncli

import (
	"errors"
	"testing"

	"github.com/OneNoted/zeshion/model"
	"github.com/OneNoted/zeshion/tmux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionTargetExistsAcceptsNameOrID(t *testing.T) {
	sessions := []*model.TmuxSession{{ID: "w2", Name: "repo"}}

	assert.True(t, sessionTargetExists(sessions, "repo"))
	assert.True(t, sessionTargetExists(sessions, "w2"))
	assert.False(t, sessionTargetExists(sessions, "missing"))
}

func TestResolveWindowTargetAcceptsHerdrTabID(t *testing.T) {
	windows := []*model.TmuxWindow{{ID: "w2:t2", Name: "logs"}}

	target, found, err := resolveWindowTarget(windows, "herdr", "w2", "w2:t2")

	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "w2:t2", target)
}

func TestResolveWindowTargetRejectsDuplicateHerdrLabels(t *testing.T) {
	windows := []*model.TmuxWindow{
		{ID: "w2:t2", Name: "logs"},
		{ID: "w2:t3", Name: "logs"},
	}

	_, found, err := resolveWindowTarget(windows, "herdr", "w2", "logs")

	require.Error(t, err)
	assert.False(t, found)
	assert.ErrorContains(t, err, "w2:t2, w2:t3")
}

func TestResolveWindowTargetPreservesTmuxDuplicateSelection(t *testing.T) {
	windows := []*model.TmuxWindow{
		{ID: "1", Name: "logs"},
		{ID: "2", Name: "logs"},
	}

	target, found, err := resolveWindowTarget(windows, "tmux", "dev", "logs")

	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "dev:1", target)
}
func TestResolveWindowTargetPreservesZellijTargetSession(t *testing.T) {
	windows := []*model.TmuxWindow{{ID: "2", Name: "logs"}}

	target, found, err := resolveWindowTarget(windows, "zellij", "dev", "logs")

	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "dev:logs", target)
}

func TestResolveWindowTargetUsesHerdrIDForUniqueLabel(t *testing.T) {
	windows := []*model.TmuxWindow{{ID: "w2:t2", Name: "logs"}}

	target, found, err := resolveWindowTarget(windows, "herdr", "w2", "logs")

	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "w2:t2", target)
}

func TestCreateWindowFocusesNewHerdrTab(t *testing.T) {
	mux := tmux.NewMockTmux(t)
	mux.On("NewWindowInSession", "tmp", "/tmp", "w2", "").Return("w2:t2", nil).Once()
	mux.On("SelectWindow", "w2:t2").Return("w2:t2", nil).Once()

	require.NoError(t, createWindow(mux, "herdr", "tmp", "/tmp", "w2"))
}

func TestCreateWindowDoesNotReselectTmuxWindow(t *testing.T) {
	mux := tmux.NewMockTmux(t)
	mux.On("NewWindowInSession", "tmp", "/tmp", "dev", "").Return("", nil).Once()

	require.NoError(t, createWindow(mux, "tmux", "tmp", "/tmp", "dev"))
}

func TestCreateWindowReportsHerdrFocusFailure(t *testing.T) {
	mux := tmux.NewMockTmux(t)
	mux.On("NewWindowInSession", "tmp", "/tmp", "w2", "").Return("w2:t2", nil).Once()
	mux.On("SelectWindow", "w2:t2").Return("", errors.New("focus failed")).Once()

	err := createWindow(mux, "herdr", "tmp", "/tmp", "w2")
	require.Error(t, err)
	assert.ErrorContains(t, err, "failed to select window")
}
