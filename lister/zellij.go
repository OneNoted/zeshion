package lister

import (
	"fmt"

	"github.com/notes/zeshion/model"
)

func zellijKey(name string) string {
	return fmt.Sprintf("zellij:%s", name)
}

func listZellij(l *RealLister) (model.SeshSessions, error) {
	zellijSessions, err := l.zellij.ListSessions()
	if err != nil {
		return model.SeshSessions{}, fmt.Errorf("couldn't list zellij sessions: %q", err)
	}

	directory := make(map[string]model.SeshSession)
	orderedIndex := []string{}

	for _, session := range zellijSessions {
		key := zellijKey(session.Name)
		orderedIndex = append(orderedIndex, key)
		directory[key] = model.SeshSession{
			Src:      "zellij",
			Name:     session.Name,
			Path:     session.Path,
			Attached: session.Attached,
			Windows:  session.Windows,
		}
	}

	return model.SeshSessions{
		Directory:    directory,
		OrderedIndex: orderedIndex,
	}, nil
}

func (l *RealLister) FindZellijSession(name string) (model.SeshSession, bool) {
	sessions, err := listZellij(l)
	if err != nil {
		return model.SeshSession{}, false
	}
	key := zellijKey(name)
	if session, exists := sessions.Directory[key]; exists {
		return session, exists
	}
	return model.SeshSession{}, false
}

func (l *RealLister) GetLastZellijSession() (model.SeshSession, bool) {
	sessions, err := listZellij(l)
	if err != nil {
		return model.SeshSession{}, false
	}
	if len(sessions.OrderedIndex) < 2 {
		return model.SeshSession{}, false
	}
	return sessions.Directory[sessions.OrderedIndex[1]], true
}
