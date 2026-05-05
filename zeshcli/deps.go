package zeshcli

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/notes/zesh/cache"
	"github.com/notes/zesh/cloner"
	"github.com/notes/zesh/configurator"
	"github.com/notes/zesh/connector"
	"github.com/notes/zesh/dir"
	"github.com/notes/zesh/execwrap"
	"github.com/notes/zesh/git"
	"github.com/notes/zesh/home"
	"github.com/notes/zesh/icon"
	"github.com/notes/zesh/json"
	"github.com/notes/zesh/lister"
	"github.com/notes/zesh/ls"
	"github.com/notes/zesh/model"
	"github.com/notes/zesh/namer"
	"github.com/notes/zesh/oswrap"
	"github.com/notes/zesh/pathwrap"
	"github.com/notes/zesh/picker"
	"github.com/notes/zesh/previewer"
	"github.com/notes/zesh/replacer"
	"github.com/notes/zesh/runtimewrap"
	"github.com/notes/zesh/shell"
	"github.com/notes/zesh/startup"
	"github.com/notes/zesh/tmux"
	"github.com/notes/zesh/tmuxinator"
	"github.com/notes/zesh/zellij"
	"github.com/notes/zesh/zoxide"
)

// BaseDeps holds config-free dependencies that can be constructed eagerly.
type BaseDeps struct {
	Exec       execwrap.Exec
	Os         oswrap.Os
	Path       pathwrap.Path
	Runtime    runtimewrap.Runtime
	Home       home.Home
	Shell      shell.Shell
	Json       json.Json
	Replacer   replacer.Replacer
	Git        git.Git
	Dir        dir.Dir
	Zoxide     zoxide.Zoxide
	Tmuxinator tmuxinator.Tmuxinator
}

// Deps holds all dependencies including config-dependent ones.
type Deps struct {
	BaseDeps
	Config        model.Config
	Tmux          tmux.Tmux
	Mux           tmux.Tmux
	MuxName       string
	Lister        lister.Lister
	Picker        picker.Picker
	CachingLister *lister.CachingLister
	Startup       startup.Startup
	Namer         namer.Namer
	Connector     connector.Connector
	Icon          icon.Icon
	Previewer     previewer.Previewer
	Cloner        cloner.Cloner
}

// NewBaseDeps constructs all config-free dependencies.
func NewBaseDeps() *BaseDeps {
	exec := execwrap.NewExec()
	os := oswrap.NewOs()
	path := pathwrap.NewPath()
	runtime := runtimewrap.NewRunTime()

	h := home.NewHome(os)
	sh := shell.NewShell(exec, h)
	j := json.NewJson()
	r := replacer.NewReplacer()

	g := git.NewGit(sh)
	d := dir.NewDir(os, g, path)
	z := zoxide.NewZoxide(sh)
	ti := tmuxinator.NewTmuxinator(sh)

	return &BaseDeps{
		Exec:       exec,
		Os:         os,
		Path:       path,
		Runtime:    runtime,
		Home:       h,
		Shell:      sh,
		Json:       j,
		Replacer:   r,
		Git:        g,
		Dir:        d,
		Zoxide:     z,
		Tmuxinator: ti,
	}
}

// BuildAll loads config and constructs all config-dependent dependencies.
func (b *BaseDeps) BuildAll(configPath string, muxOverride string) (*Deps, error) {
	config, err := configurator.NewConfiguratorWithPath(b.Os, b.Path, b.Runtime, configPath).GetConfig()
	if err != nil {
		return nil, err
	}

	slog.Debug("deps: BuildAll", "config", config)

	t := tmux.NewTmux(b.Os, b.Shell, config.TmuxCommand)
	zj := zellij.NewZellij(b.Os, b.Shell, config.ZellijCommand)
	muxName := selectMuxName(b.Os, muxOverride, config.Multiplexer)
	var selectedMux tmux.Tmux = t
	if muxName == "zellij" {
		selectedMux = zj
	}

	l := ls.NewLs(config, b.Shell)
	li := lister.NewListerWithMux(config, b.Home, t, zj, selectedMux, b.Zoxide, b.Tmuxinator)

	var usedLister lister.Lister = li
	var cachedLi *lister.CachingLister
	if config.Cache {
		fc := cache.NewFileCache()
		cachedLi = lister.NewCachingLister(li, fc)
		usedLister = cachedLi
	}

	s := startup.NewStartup(b.Os, config, usedLister, selectedMux, b.Home, b.Replacer)
	n := namer.NewNamer(b.Path, b.Git, b.Home, config)
	c := connector.NewConnector(config, b.Dir, b.Home, usedLister, n, s, t, selectedMux, muxName, b.Zoxide, b.Tmuxinator)
	ic := icon.NewIcon(config)
	p := previewer.NewPreviewer(usedLister, selectedMux, ic, b.Dir, b.Home, l, config, b.Shell)
	cl := cloner.NewCloner(c, b.Git)
	pk := picker.NewPicker(config)

	return &Deps{
		BaseDeps:      *b,
		Config:        config,
		Tmux:          t,
		Mux:           selectedMux,
		MuxName:       muxName,
		Lister:        usedLister,
		Picker:        pk,
		CachingLister: cachedLi,
		Startup:       s,
		Namer:         n,
		Connector:     c,
		Icon:          ic,
		Previewer:     p,
		Cloner:        cl,
	}, nil
}

// buildDeps reads the --config flag from cobra and builds all dependencies.
func buildDeps(cmd *cobra.Command, base *BaseDeps) (*Deps, error) {
	configPath, _ := cmd.Root().PersistentFlags().GetString("config")
	muxOverride, _ := cmd.Root().PersistentFlags().GetString("multiplexer")
	deps, err := base.BuildAll(configPath, muxOverride)
	if err != nil {
		var human *configurator.ConfigError
		if errors.As(err, &human) {
			fmt.Printf("Couldn't parse config, err: %v\n details:\n %s\n", err.Error(), human.Human())
		}
		slog.Error("buildDeps", "error", err)
		return nil, err
	}
	return deps, nil
}

func selectMuxName(os oswrap.Os, override string, configured string) string {
	choice := override
	if choice == "" || choice == "auto" {
		if os.Getenv("ZELLIJ") != "" || os.Getenv("ZELLIJ_SESSION_NAME") != "" {
			return "zellij"
		}
		if os.Getenv("TMUX") != "" {
			return "tmux"
		}
		choice = configured
	}
	if choice == "" || choice == "auto" {
		return "tmux"
	}
	if choice == "zellij" {
		return "zellij"
	}
	return "tmux"
}
