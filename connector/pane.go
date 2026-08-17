package connector

import (
	"fmt"
	"strings"

	"github.com/OneNoted/zeshion/model"
)

func isMuxPaneFormat(name string) bool {
	return strings.Contains(name, "/") && !strings.HasPrefix(name, "/")
}

func muxPaneStrategy(c *RealConnector, name string) (model.Connection, error) {
	if !c.mux.IsAttached() || !isMuxPaneFormat(name) {
		return model.Connection{Found: false}, nil
	}

	sessions, err := c.lister.ListTmuxPanes()
	if err != nil {
		return model.Connection{Found: false}, nil
	}

	for _, key := range sessions.OrderedIndex {
		session := sessions.Directory[key]
		if session.Name == name {
			return model.Connection{
				Found:       true,
				Session:     session,
				New:         false,
				AddToZoxide: false,
			}, nil
		}
	}

	return model.Connection{Found: false}, nil
}

func connectToMuxPane(c *RealConnector, connection model.Connection, _ model.ConnectOpts) (string, error) {
	panes, err := c.mux.ListTmuxPanes()
	if err != nil {
		return "", err
	}
	if connection.Session.ID != "" {
		for _, pane := range panes {
			if pane.PaneID == connection.Session.ID {
				return c.mux.SelectPane(pane.WindowIndex, pane.PaneIndex)
			}
		}
		return "", fmt.Errorf("pane not found: %s", connection.Session.ID)
	}

	sessions, err := c.lister.ListTmuxPanes()
	if err != nil {
		return "", err
	}

	var targetPaneKey string
	for _, key := range sessions.OrderedIndex {
		if sessions.Directory[key].Name == connection.Session.Name {
			targetPaneKey = key
			break
		}
	}

	for _, pane := range panes {
		paneKey := fmt.Sprintf("%s:%s/%s", connection.Session.Src, pane.WindowName, pane.PaneID)
		if paneKey == targetPaneKey {
			return c.mux.SelectPane(pane.WindowIndex, pane.PaneIndex)
		}
	}

	return "", fmt.Errorf("pane not found: %s", connection.Session.Name)
}
