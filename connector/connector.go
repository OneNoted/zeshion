package connector

import (
	"github.com/notes/zeshion/dir"
	"github.com/notes/zeshion/home"
	"github.com/notes/zeshion/lister"
	"github.com/notes/zeshion/model"
	"github.com/notes/zeshion/namer"
	"github.com/notes/zeshion/startup"
	"github.com/notes/zeshion/tmux"
	"github.com/notes/zeshion/tmuxinator"
	"github.com/notes/zeshion/zoxide"
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
