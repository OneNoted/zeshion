package connector

import "github.com/OneNoted/zeshion/model"

func herdrStrategy(c *RealConnector, name string) (model.Connection, error) {
	if c.herdr == nil {
		return model.Connection{Found: false}, nil
	}
	session, exists, err := c.lister.FindHerdrSession(name)
	if err != nil {
		return model.Connection{}, err
	}
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

func connectToHerdr(c *RealConnector, connection model.Connection, opts model.ConnectOpts) (string, error) {
	return connectWith(c, c.herdr, connection, opts)
}
