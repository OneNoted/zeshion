package previewer

import (
	"github.com/OneNoted/zeshion/lister"
	"github.com/OneNoted/zeshion/tmux"
)

type HerdrPreviewStrategy struct {
	lister lister.Lister
	herdr  tmux.Tmux
}

func NewHerdrStrategy(lister lister.Lister, herdr tmux.Tmux) *HerdrPreviewStrategy {
	return &HerdrPreviewStrategy{lister: lister, herdr: herdr}
}

func (s *HerdrPreviewStrategy) Execute(name string) (string, error) {
	session, exists, err := s.lister.FindHerdrSession(name)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", nil
	}
	return s.herdr.CapturePane(session.ID)
}
