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

func testConnector(mockDir dir.Dir, mockHome home.Home, mockLister lister.Lister, mockNamer namer.Namer, mockStartup startup.Startup, mockTmux tmux.Tmux, mockZoxide zoxide.Zoxide, mockTmuxinator tmuxinator.Tmuxinator) *RealConnector {
	return &RealConnector{
		config:     model.Config{},
		dir:        mockDir,
		home:       mockHome,
		lister:     mockLister,
		namer:      mockNamer,
		startup:    mockStartup,
		tmux:       mockTmux,
		mux:        mockTmux,
		muxName:    "tmux",
		zoxide:     mockZoxide,
		tmuxinator: mockTmuxinator,
	}
}
