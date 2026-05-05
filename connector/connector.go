package connector

import (
	"github.com/notes/zesh/dir"
	"github.com/notes/zesh/home"
	"github.com/notes/zesh/lister"
	"github.com/notes/zesh/model"
	"github.com/notes/zesh/namer"
	"github.com/notes/zesh/startup"
	"github.com/notes/zesh/tmux"
	"github.com/notes/zesh/tmuxinator"
	"github.com/notes/zesh/zoxide"
)

type Connector interface {
	Connect(name string, opts model.ConnectOpts) (string, error)
}

type RealConnector struct {
	config     model.Config
	dir        dir.Dir
	home       home.Home
	lister     lister.Lister
	namer      namer.Namer
	startup    startup.Startup
	tmux       tmux.Tmux
	zellij     tmux.Tmux
	mux        tmux.Tmux
	muxName    string
	zoxide     zoxide.Zoxide
	tmuxinator tmuxinator.Tmuxinator
}

func NewConnector(
	config model.Config,
	dir dir.Dir,
	home home.Home,
	lister lister.Lister,
	namer namer.Namer,
	startup startup.Startup,
	tmux tmux.Tmux,
	zellij tmux.Tmux,
	mux tmux.Tmux,
	muxName string,
	zoxide zoxide.Zoxide,
	tmuxinator tmuxinator.Tmuxinator,
) Connector {
	return &RealConnector{
		config,
		dir,
		home,
		lister,
		namer,
		startup,
		tmux,
		zellij,
		mux,
		muxName,
		zoxide,
		tmuxinator,
	}
}
