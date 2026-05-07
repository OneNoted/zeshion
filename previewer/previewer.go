package previewer

import (
	"github.com/notes/zeshion/dir"
	"github.com/notes/zeshion/home"
	"github.com/notes/zeshion/icon"
	"github.com/notes/zeshion/lister"
	"github.com/notes/zeshion/ls"
	"github.com/notes/zeshion/model"
	"github.com/notes/zeshion/shell"
	"github.com/notes/zeshion/tmux"
)

type Previewer interface {
	// Previews a session or directory
	Preview(name string) (string, error)
}

type RealPreviewer struct {
	icon       icon.Icon
	strategies []PreviewStrategy
}

func NewPreviewer(
	lister lister.Lister,
	tmux tmux.Tmux,
	icon icon.Icon,
	dir dir.Dir,
	home home.Home,
	ls ls.Ls,
	config model.Config,
	shell shell.Shell,
) Previewer {
	strategies := []PreviewStrategy{
		NewTmuxStrategy(lister, tmux),
		NewConfigStrategy(lister, shell),
		NewDefaultConfigStrategy(lister, config, ls),
		NewDirectoryStrategy(home, dir, ls),
	}

	return &RealPreviewer{
		icon:       icon,
		strategies: strategies,
	}
}

func (p *RealPreviewer) Preview(name string) (string, error) {
	trimmedName := p.icon.RemoveIcon(name)

	for _, strategy := range p.strategies {
		output, err := strategy.Execute(trimmedName)
		if err != nil {
			return "", err
		}
		if output != "" {
			return output, nil
		}
	}
	return "", nil
}
