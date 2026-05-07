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
