package previewer

import (
	"testing"

	"github.com/OneNoted/zeshion/lister"
	"github.com/OneNoted/zeshion/model"
	"github.com/OneNoted/zeshion/tmux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHerdrPreviewUsesWorkspaceID(t *testing.T) {
	mockLister := lister.NewMockLister(t)
	backend := tmux.NewMockTmux(t)
	mockLister.On("FindHerdrSession", "repo").Return(model.SeshSession{ID: "w2", Src: "herdr", Name: "repo"}, true, nil).Once()
	backend.On("CapturePane", "w2").Return("screen contents", nil).Once()

	output, err := NewHerdrStrategy(mockLister, backend).Execute("repo")
	require.NoError(t, err)
	assert.Equal(t, "screen contents", output)
}
