package herdr

type apiResponse[T any] struct {
	Result T `json:"result"`
}

type snapshotResult struct {
	Snapshot snapshot `json:"snapshot"`
}

type snapshot struct {
	FocusedWorkspaceID string      `json:"focused_workspace_id"`
	FocusedTabID       string      `json:"focused_tab_id"`
	FocusedPaneID      string      `json:"focused_pane_id"`
	Workspaces         []workspace `json:"workspaces"`
	Tabs               []tab       `json:"tabs"`
	Panes              []pane      `json:"panes"`
}

type workspace struct {
	ID          string    `json:"workspace_id"`
	ActiveTabID string    `json:"active_tab_id"`
	Label       string    `json:"label"`
	Number      int       `json:"number"`
	PaneCount   int       `json:"pane_count"`
	TabCount    int       `json:"tab_count"`
	Focused     bool      `json:"focused"`
	Worktree    *worktree `json:"worktree"`
}

type worktree struct {
	CheckoutPath string `json:"checkout_path"`
}

type tab struct {
	ID          string `json:"tab_id"`
	WorkspaceID string `json:"workspace_id"`
	Label       string `json:"label"`
	Number      int    `json:"number"`
	PaneCount   int    `json:"pane_count"`
	Focused     bool   `json:"focused"`
}

type pane struct {
	ID            string `json:"pane_id"`
	WorkspaceID   string `json:"workspace_id"`
	TabID         string `json:"tab_id"`
	Cwd           string `json:"cwd"`
	ForegroundCwd string `json:"foreground_cwd"`
	Title         string `json:"terminal_title_stripped"`
	Agent         string `json:"agent"`
	Focused       bool   `json:"focused"`
}

type createWorkspaceResult struct {
	Workspace workspace `json:"workspace"`
	Tab       tab       `json:"tab"`
	RootPane  pane      `json:"root_pane"`
}

type createTabResult struct {
	Tab      tab  `json:"tab"`
	RootPane pane `json:"root_pane"`
}
