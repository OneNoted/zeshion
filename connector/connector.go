package connector

import (
	"github.com/OneNoted/zeshion/dir"
	"github.com/OneNoted/zeshion/home"
	"github.com/OneNoted/zeshion/lister"
	"github.com/OneNoted/zeshion/model"
	"github.com/OneNoted/zeshion/namer"
	"github.com/OneNoted/zeshion/startup"
	"github.com/OneNoted/zeshion/tmux"
	"github.com/OneNoted/zeshion/tmuxinator"
	"github.com/OneNoted/zeshion/zoxide"
)

type Connector interface {
	Connect(name string, opts model.ConnectOpts) (string, error)
	ConnectSession(session model.SeshSession, opts model.ConnectOpts) (string, error)
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
	herdr      tmux.Tmux
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
	herdr tmux.Tmux,
	mux tmux.Tmux,
	muxName string,
	zoxide zoxide.Zoxide,
	tmuxinator tmuxinator.Tmuxinator,
) Connector {
	return &RealConnector{
		config: config, dir: dir, home: home, lister: lister, namer: namer, startup: startup,
		tmux: tmux, zellij: zellij, herdr: herdr, mux: mux, muxName: muxName,
		zoxide: zoxide, tmuxinator: tmuxinator,
	}
}
