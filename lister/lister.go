package lister

import (
	"github.com/OneNoted/zeshion/home"
	"github.com/OneNoted/zeshion/model"
	"github.com/OneNoted/zeshion/tmux"
	"github.com/OneNoted/zeshion/tmuxinator"
	"github.com/OneNoted/zeshion/zoxide"
)

type Lister interface {
	List(opts ListOptions) (model.SeshSessions, error)
	ListTmuxPanes() (model.SeshSessions, error)
	FindTmuxSession(name string) (model.SeshSession, bool)
	FindZellijSession(name string) (model.SeshSession, bool)
	FindHerdrSession(name string) (model.SeshSession, bool, error)
	GetAttachedTmuxSession() (model.SeshSession, bool)
	GetAttachedHerdrSession() (model.SeshSession, bool)
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
	herdr      tmux.Tmux
	mux        tmux.Tmux
	muxName    string
	zoxide     zoxide.Zoxide
	tmuxinator tmuxinator.Tmuxinator
}

func NewLister(config model.Config, home home.Home, tmux tmux.Tmux, zoxide zoxide.Zoxide, tmuxinator tmuxinator.Tmuxinator) Lister {
	return NewListerWithMux(config, home, tmux, tmux, nil, tmux, "tmux", zoxide, tmuxinator)
}

func NewListerWithMux(config model.Config, home home.Home, tmux tmux.Tmux, zellij tmux.Tmux, herdr tmux.Tmux, mux tmux.Tmux, muxName string, zoxide zoxide.Zoxide, tmuxinator tmuxinator.Tmuxinator) Lister {
	return &RealLister{
		config: config, home: home, tmux: tmux, zellij: zellij, herdr: herdr,
		mux: mux, muxName: muxName, zoxide: zoxide, tmuxinator: tmuxinator,
	}
}
