package lister

import (
	"github.com/notes/zesh/home"
	"github.com/notes/zesh/model"
	"github.com/notes/zesh/tmux"
	"github.com/notes/zesh/tmuxinator"
	"github.com/notes/zesh/zoxide"
)

type Lister interface {
	List(opts ListOptions) (model.SeshSessions, error)
	ListTmuxPanes() (model.SeshSessions, error)
	FindTmuxSession(name string) (model.SeshSession, bool)
	FindZellijSession(name string) (model.SeshSession, bool)
	GetAttachedTmuxSession() (model.SeshSession, bool)
	GetLastTmuxSession() (model.SeshSession, bool)
	GetLastZellijSession() (model.SeshSession, bool)
	FindConfigSession(name string) (model.SeshSession, bool)
	FindConfigWildcard(path string) (model.WildcardConfig, bool)
	FindZoxideSession(name string) (model.SeshSession, bool)
	FindTmuxinatorConfig(name string) (model.SeshSession, bool)
}

type RealLister struct {
	config     model.Config
	home       home.Home
	tmux       tmux.Tmux
	zellij     tmux.Tmux
	mux        tmux.Tmux
	zoxide     zoxide.Zoxide
	tmuxinator tmuxinator.Tmuxinator
}

func NewLister(config model.Config, home home.Home, tmux tmux.Tmux, zellij tmux.Tmux, mux tmux.Tmux, zoxide zoxide.Zoxide, tmuxinator tmuxinator.Tmuxinator) Lister {
	return &RealLister{config, home, tmux, zellij, mux, zoxide, tmuxinator}
}
