package lister

import (
	"testing"

	"github.com/OneNoted/zeshion/model"
	"github.com/OneNoted/zeshion/tmux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListHerdrKeepsDuplicateLabelsSelectableByID(t *testing.T) {
	backend := tmux.NewMockTmux(t)
	backend.On("ListSessions").Return([]*model.TmuxSession{
		{ID: "w1", Name: "repo", Path: "/work/one"},
		{ID: "w2", Name: "repo", Path: "/work/two", Attached: 1},
	}, nil).Once()
	l := &RealLister{herdr: backend}

	sessions, err := listHerdr(l)
	require.NoError(t, err)
	assert.Equal(t, []string{"herdr:w1", "herdr:w2"}, sessions.OrderedIndex)
	assert.Equal(t, "w1", sessions.Directory["herdr:w1"].ID)
	assert.Equal(t, "w2", sessions.Directory["herdr:w2"].ID)
}

func TestFindHerdrSessionRequiresUniqueLabel(t *testing.T) {
	backend := tmux.NewMockTmux(t)
	backend.On("ListSessions").Return([]*model.TmuxSession{
		{ID: "w1", Name: "repo"},
		{ID: "w2", Name: "repo"},
	}, nil).Twice()
	l := &RealLister{herdr: backend}

	byID, found, err := l.FindHerdrSession("w2")
	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, "w2", byID.ID)

	_, found, err = l.FindHerdrSession("repo")
	assert.False(t, found)
	assert.ErrorContains(t, err, "ambiguous")
}
