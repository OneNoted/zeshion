package connector

import (
	"github.com/OneNoted/zeshion/model"
	"github.com/OneNoted/zeshion/tmux"
)

func tmuxStrategy(c *RealConnector, name string) (model.Connection, error) {
	session, exists := c.lister.FindTmuxSession(name)
	if !exists {
		return model.Connection{Found: false}, nil
	}
	return model.Connection{
		Found:       true,
		Session:     session,
		New:         false,
		AddToZoxide: true,
	}, nil
}

func connectToTmux(c *RealConnector, connection model.Connection, opts model.ConnectOpts) (string, error) {
	return connectWith(c, c.tmux, connection, opts)
}

func connectToMux(c *RealConnector, connection model.Connection, opts model.ConnectOpts) (string, error) {
	return connectWith(c, c.mux, connection, opts)
}

func connectWith(c *RealConnector, target tmux.Tmux, connection model.Connection, opts model.ConnectOpts) (string, error) {
	targetSession := connection.Session.Name
	if connection.New {
		// Resolve the startup command before creation so the first pane receives it atomically.
		var rawCmd string
		if opts.Command != "" {
			rawCmd = opts.Command
		} else {
			resolved, err := c.startup.ResolveCommand(connection.Session)
			if err != nil {
				return "", err
			}
			rawCmd = resolved
		}
		shellCmd := rawCmd
		if target != c.herdr {
			shellCmd = c.startup.WrapForShell(rawCmd)
		}
		createdID, err := target.NewSession(connection.Session.Name, connection.Session.Path, shellCmd)
		if err != nil {
			return "", err
		}
		if createdID != "" {
			targetSession = createdID
		}
		if opts.Command == "" {
			startupSession := connection.Session
			startupSession.ID = createdID
			if _, err := c.startup.Exec(startupSession); err != nil {
				return "", err
			}
		}
	} else if connection.Session.ID != "" {
		targetSession = connection.Session.ID
	}
	return target.SwitchOrAttach(targetSession, opts)
}
