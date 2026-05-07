package lister

import (
	"log"
	"testing"

	"github.com/OneNoted/zeshion/home"
	"github.com/OneNoted/zeshion/model"
	"github.com/OneNoted/zeshion/tmux"
	"github.com/OneNoted/zeshion/tmuxinator"
	"github.com/OneNoted/zeshion/zoxide"
	"github.com/stretchr/testify/assert"
)

func TestListTmuxinatorConfigs(t *testing.T) {
	t.Run("should list tmuxinator configs", func(t *testing.T) {
		mockConfig := model.Config{}
		mockHome := new(home.MockHome)
		mockZoxide := new(zoxide.MockZoxide)
		mockTmux := new(tmux.MockTmux)
		mockTmuxinator := new(tmuxinator.MockTmuxinator)
		mockTmuxinator.On("List").Return([]*model.TmuxinatorConfig{
			{Name: "sesh"},
			{Name: "dotfiles"},
		}, nil)

		lister := NewLister(mockConfig, mockHome, mockTmux, mockZoxide, mockTmuxinator)

		realLister, ok := lister.(*RealLister)
		if !ok {
			log.Fatal("Cannot convert lister to *RealLister")
		}
		sessions, err := listTmuxinator(realLister)
		assert.Equal(t, "tmuxinator:sesh", sessions.OrderedIndex[0])
		assert.Equal(t, "tmuxinator:dotfiles", sessions.OrderedIndex[1])
		assert.Nil(t, err)
	})
}
