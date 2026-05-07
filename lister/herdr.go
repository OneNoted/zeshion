package lister

import (
	"fmt"

	"github.com/OneNoted/zeshion/model"
)

func herdrKey(id string) string {
	return fmt.Sprintf("herdr:%s", id)
}

func listHerdr(l *RealLister) (model.SeshSessions, error) {
	if l.herdr == nil {
		return model.SeshSessions{Directory: model.SeshSessionMap{}}, nil
	}
	herdrSessions, err := l.herdr.ListSessions()
	if err != nil {
		return model.SeshSessions{}, fmt.Errorf("couldn't list Herdr workspaces: %q", err)
	}

	directory := make(model.SeshSessionMap, len(herdrSessions))
	orderedIndex := make([]string, 0, len(herdrSessions))
	for _, session := range herdrSessions {
		id := session.ID
		if id == "" {
			id = session.Name
		}
		key := herdrKey(id)
		orderedIndex = append(orderedIndex, key)
		directory[key] = model.SeshSession{
			ID:       id,
			Src:      "herdr",
			Name:     session.Name,
			Path:     session.Path,
			Attached: session.Attached,
			Windows:  session.Windows,
		}
	}
	return model.SeshSessions{Directory: directory, OrderedIndex: orderedIndex}, nil
}

func (l *RealLister) FindHerdrSession(name string) (model.SeshSession, bool, error) {
	sessions, err := listHerdr(l)
	if err != nil {
		return model.SeshSession{}, false, err
	}
	if session, exists := sessions.Directory[herdrKey(name)]; exists {
		return session, true, nil
	}

	var found model.SeshSession
	matches := 0
	for _, key := range sessions.OrderedIndex {
		session := sessions.Directory[key]
		if session.Name == name {
			found = session
			matches++
		}
	}
	if matches > 1 {
		return model.SeshSession{}, false, fmt.Errorf("Herdr workspace label %q is ambiguous; use a workspace ID", name)
	}
	return found, matches == 1, nil
}

func (l *RealLister) GetAttachedHerdrSession() (model.SeshSession, bool) {
	sessions, err := listHerdr(l)
	if err != nil {
		return model.SeshSession{}, false
	}
	for _, key := range sessions.OrderedIndex {
		session := sessions.Directory[key]
		if session.Attached != 0 {
			return session, true
		}
	}
	return model.SeshSession{}, false
}
