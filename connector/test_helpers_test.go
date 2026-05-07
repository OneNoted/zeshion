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

func testConnector(mockDir dir.Dir, mockHome home.Home, mockLister lister.Lister, mockNamer namer.Namer, mockStartup startup.Startup, mockTmux tmux.Tmux, mockZoxide zoxide.Zoxide, mockTmuxinator tmuxinator.Tmuxinator) *RealConnector {
	return &RealConnector{
		config:     model.Config{},
		dir:        mockDir,
		home:       mockHome,
		lister:     mockLister,
		namer:      mockNamer,
		startup:    mockStartup,
		tmux:       mockTmux,
		zellij:     mockTmux,
		mux:        mockTmux,
		muxName:    "tmux",
		zoxide:     mockZoxide,
		tmuxinator: mockTmuxinator,
	}
}
