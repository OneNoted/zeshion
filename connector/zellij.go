package connector

import "github.com/notes/zesh/model"

func zellijStrategy(c *RealConnector, name string) (model.Connection, error) {
	session, exists := c.lister.FindZellijSession(name)
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

func connectToZellij(c *RealConnector, connection model.Connection, opts model.ConnectOpts) (string, error) {
	return connectWith(c, c.zellij, connection, opts)
}
