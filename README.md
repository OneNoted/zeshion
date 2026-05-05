# zesh

`zesh` is a Go CLI for managing terminal multiplexer sessions across tmux and Zellij. It is based on the MIT-licensed `joshmedeski/sesh` codebase and keeps the `sesh` command shape wherever practical while adding first-class Zellij support.

## Status

This repository currently builds a `zesh` binary only. Users who want command-name muscle memory can add their own shell alias, for example:

```sh
alias sesh=zesh
```

## Compatibility Goals

- Keep `sesh`-style subcommands and aliases: `list/l`, `connect/cn`, `clone/cl`, `picker/pick/pk`, `preview/p`, `root`, `last/L`, and `window/w`.
- Keep existing session config concepts: configured sessions, wildcard sessions, startup commands, preview commands, windows, zoxide, tmuxinator, caching, icons, and JSON output.
- Support `~/.config/zesh/zesh.toml` first and fall back to `~/.config/sesh/sesh.toml` when no zesh config exists.
- Keep tmux-specific integrations such as `tmuxinator` and `tmuxp` on the tmux path. Use native Zellij sessions, tabs, panes, and layouts for Zellij.

## Multiplexer Selection

`zesh` defaults to `--multiplexer auto`.

Auto mode chooses:

1. Zellij when running inside Zellij.
2. tmux when running inside tmux.
3. The `multiplexer` value from config.
4. tmux as the compatibility fallback.

Override it per command:

```sh
zesh --multiplexer tmux list
zesh --multiplexer zellij connect my-session
```

## Common Commands

```sh
zesh list
zesh list --tmux
zesh list --zellij
zesh list --config --zoxide
zesh connect my-session
zesh connect --root "$PWD"
zesh picker
zesh window
zesh window ~/projects/my-app
zesh preview my-session
```

`zesh window` maps to tmux windows when tmux is selected and Zellij tabs when Zellij is selected.

## Configuration

Create `~/.config/zesh/zesh.toml`:

```toml
multiplexer = "auto"
tmux_command = "tmux"
zellij_command = "zellij"

[default_session]
startup_command = "nvim"
preview_command = "ls -la {}"

[[session]]
name = "dotfiles"
path = "~/dotfiles"
startup_command = "nvim"
```

The schema lives at `zesh.schema.json`.

## Development

```sh
just mock
go test ./...
go build -o zesh .
```

Version-control work in this repo uses Jujutsu (`jj`).
