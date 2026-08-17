package connector

import (
	"fmt"

	"github.com/OneNoted/zeshion/model"
)

// TODO: send to logging (local txt file?)
func (c *RealConnector) Connect(name string, opts model.ConnectOpts) (string, error) {
	// TODO: make it configurable to change the order of connection establishments?
	// ["tmux", "config", "dir", "zoxide"]
	// TODO: make it configurable to disable certain strategies (including flags for optimized fzf commands)
	// sesh connect --config (sesh list --config | fzf)
	strategies := []func(*RealConnector, string) (model.Connection, error){muxPaneStrategy}
	switch c.muxName {
	case "herdr":
		strategies = append(strategies, herdrStrategy, tmuxStrategy, zellijStrategy)
	case "zellij":
		strategies = append(strategies, zellijStrategy, tmuxStrategy, herdrStrategy)
	default:
		strategies = append(strategies, tmuxStrategy, zellijStrategy, herdrStrategy)
	}
	strategies = append(
		strategies,
		tmuxinatorStrategy,
		configStrategy,
		configWildcardStrategy,
		dirStrategy,
		zoxideStrategy,
	)

	for _, strategy := range strategies {
		if connection, err := strategy(c, name); err != nil {
			return "", fmt.Errorf("failed to establish connection: %w", err)
		} else if connection.Found {
			return c.connectSession(connection, opts)
		}
	}

	return "", fmt.Errorf("no connection found for '%s'", name)
}

func shouldAddToZoxide(source string) bool {
	switch source {
	case "tmuxinator", "tmux-pane", "zellij-pane", "herdr-pane":
		return false
	default:
		return true
	}
}

func (c *RealConnector) revalidateLiveSession(session model.SeshSession) (model.SeshSession, error) {
	switch session.Src {
	case "tmux":
		live, exists := c.lister.FindTmuxSession(session.Name)
		if !exists {
			return model.SeshSession{}, fmt.Errorf("tmux session %q no longer exists", session.Name)
		}
		return live, nil
	case "zellij":
		live, exists := c.lister.FindZellijSession(session.Name)
		if !exists {
			return model.SeshSession{}, fmt.Errorf("Zellij session %q no longer exists", session.Name)
		}
		return live, nil
	case "herdr":
		if session.ID != "" {
			live, exists, err := c.lister.FindHerdrSession(session.ID)
			if err != nil {
				return model.SeshSession{}, err
			}
			if exists && live.Name == session.Name {
				return live, nil
			}
		}
		live, exists, err := c.lister.FindHerdrSession(session.Name)
		if err != nil {
			return model.SeshSession{}, err
		}
		if !exists {
			return model.SeshSession{}, fmt.Errorf("Herdr workspace %q no longer exists", session.Name)
		}
		return live, nil
	default:
		return session, nil
	}
}

func (c *RealConnector) ConnectSession(session model.SeshSession, opts model.ConnectOpts) (string, error) {
	switch session.Src {
	case "config", "config_wildcard", "dir", "zoxide", "tmuxinator":
		// These sources describe how to create a session. Preserve normal lookup
		// semantics so an existing live session wins before creating another one.
		return c.Connect(session.Name, opts)
	}
	var err error
	session, err = c.revalidateLiveSession(session)
	if err != nil {
		return "", err
	}

	connection := model.Connection{
		Found:       true,
		Session:     session,
		New:         session.Src == "config" || session.Src == "config_wildcard" || session.Src == "dir" || session.Src == "zoxide",
		AddToZoxide: shouldAddToZoxide(session.Src),
	}
	return c.connectSession(connection, opts)
}

func (c *RealConnector) connectSession(connection model.Connection, opts model.ConnectOpts) (string, error) {
	if connection.AddToZoxide && connection.Session.Path != "" {
		c.zoxide.Add(connection.Session.Path)
	}
	connectStrategy := map[string]func(c *RealConnector, connection model.Connection, opts model.ConnectOpts) (string, error){
		"tmux-pane":       connectToMuxPane,
		"zellij-pane":     connectToMuxPane,
		"herdr-pane":      connectToMuxPane,
		"tmux":            connectToTmux,
		"zellij":          connectToZellij,
		"herdr":           connectToHerdr,
		"tmuxinator":      connectToTmuxinator,
		"config":          connectToMux,
		"config_wildcard": connectToMux,
		"dir":             connectToMux,
		"zoxide":          connectToMux,
	}
	connect, exists := connectStrategy[connection.Session.Src]
	if !exists {
		return "", fmt.Errorf("unsupported session source %q", connection.Session.Src)
	}
	return connect(c, connection, opts)
}
