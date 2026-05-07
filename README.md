# zeshion

`zeshion` is a Go CLI for managing terminal multiplexer sessions across tmux and Zellij. It is based on the MIT-licensed `joshmedeski/sesh` codebase and keeps the `sesh` command shape wherever practical while adding first-class Zellij support.

## Status

This repository currently builds a `zeshion` binary only. Users who want command-name muscle memory can add their own shell alias, for example:

```sh
alias sesh=zeshion
```

## Compatibility Goals

- Keep `sesh`-style subcommands and aliases: `list/l`, `connect/cn`, `clone/cl`, `picker/pick/pk`, `preview/p`, `root`, `last/L`, and `window/w`.
- Keep existing session config concepts: configured sessions, wildcard sessions, startup commands, preview commands, windows, zoxide, tmuxinator, caching, icons, and JSON output.
- Support `~/.config/zeshion/zeshion.toml` first, then legacy `~/.config/zesh/zesh.toml`, then `~/.config/sesh/sesh.toml`.
- Keep tmux-specific integrations such as `tmuxinator` and `tmuxp` on the tmux path. Use native Zellij sessions, tabs, panes, and layouts for Zellij.

## Multiplexer Selection

`zeshion` defaults to `--multiplexer auto`.

Auto mode chooses:

1. Zellij when running inside Zellij.
2. tmux when running inside tmux.
3. The `multiplexer` value from config.
4. tmux as the compatibility fallback.

Override it per command:

```sh
zeshion --multiplexer tmux list
zeshion --multiplexer zellij connect my-session
```

## Common Commands

```sh
zeshion list
zeshion list --tmux
zeshion list --zellij
zeshion list --config --zoxide
zeshion connect my-session
zeshion connect --root "$PWD"
zeshion picker
zeshion window
zeshion window ~/projects/my-app
zeshion preview my-session
```

`zeshion window` maps to tmux windows when tmux is selected and Zellij tabs when Zellij is selected.

## Configuration

Create `~/.config/zeshion/zeshion.toml`:

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

The schema lives at `zeshion.schema.json`.

## Development

```sh
just mock
go test ./...
go build -o zeshion .
```

Version-control work in this repo uses Jujutsu (`jj`).
