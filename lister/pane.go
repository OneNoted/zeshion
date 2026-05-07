package lister

import (
	"fmt"
	"os"

	"github.com/OneNoted/zeshion/model"
)

func paneSource(muxName string) string {
	switch muxName {
	case "zellij", "herdr":
		return muxName + "-pane"
	default:
		return "tmux-pane"
	}
}

func paneKey(source, windowName, paneID string) string {
	return fmt.Sprintf("%s:%s/%s", source, windowName, paneID)
}

func tmuxPaneDisplayName(pane *model.TmuxPane) string {
	hostname, _ := os.Hostname()
	if pane.PaneTitle != "" && pane.PaneTitle != hostname {
		return pane.PaneTitle
	}
	return pane.PaneCommand
}

func listTmuxPanes(l *RealLister) (model.SeshSessions, error) {
	tmuxPanes, err := l.mux.ListTmuxPanes()
	if err != nil {
		return model.SeshSessions{}, fmt.Errorf("couldn't list tmux panes: %q", err)
	}
	source := paneSource(l.muxName)

	// Count raw names to detect duplicates needing .0, .1 suffixes
	type paneEntry struct {
		pane    *model.TmuxPane
		rawName string
	}
	entries := make([]paneEntry, len(tmuxPanes))
	nameCounts := make(map[string]int)
	for i, pane := range tmuxPanes {
		rawName := fmt.Sprintf("%s/%s", pane.WindowName, tmuxPaneDisplayName(pane))
		entries[i] = paneEntry{pane: pane, rawName: rawName}
		nameCounts[rawName]++
	}

	directory := make(map[string]model.SeshSession)
	orderedIndex := []string{}
	nameIndexes := make(map[string]int)

	for _, entry := range entries {
		name := entry.rawName
		if nameCounts[entry.rawName] > 1 {
			idx := nameIndexes[entry.rawName]
			name = fmt.Sprintf("%s.%d", entry.rawName, idx)
			nameIndexes[entry.rawName] = idx + 1
		}

		key := paneKey(source, entry.pane.WindowName, entry.pane.PaneID)
		orderedIndex = append(orderedIndex, key)
		directory[key] = model.SeshSession{
			ID:   entry.pane.PaneID,
			Src:  source,
			Name: name,
			Path: entry.pane.PanePath,
		}
	}

	return model.SeshSessions{
		Directory:    directory,
		OrderedIndex: orderedIndex,
	}, nil
}

func (l *RealLister) ListTmuxPanes() (model.SeshSessions, error) {
	return listTmuxPanes(l)
}
