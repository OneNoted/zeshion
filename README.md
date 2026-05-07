# zeshion

`zeshion` is a session manager for tmux, Zellij, and Herdr. It keeps the `sesh` command shape where practical while mapping sessions and windows onto each multiplexer's native model.

The project is derived from the MIT-licensed [`joshmedeski/sesh`](https://github.com/joshmedeski/sesh) codebase.

## Status

`zeshion` is usable locally, but first public release packaging is still being prepared. The current install path is Go-based.

## Install

Install from the module path:

```sh
go install github.com/OneNoted/zeshion@latest
```

Build from a checkout:

```sh
go build -o zeshion .
```

Build into `GOPATH/bin` with version metadata:

```sh
just build dev
```

If you want existing `sesh` muscle memory:

```sh
alias sesh=zeshion
```

## Requirements

- Go 1.25 or newer to build from source.
- `tmux` for tmux session management.
- `zellij` for Zellij session management.
- `herdr` for Herdr workspace management.
- `zoxide` for zoxide-backed directory sessions.
- A picker-compatible terminal for `zeshion picker`.

## Commands

Most commands follow the `sesh` shape:

```sh
zeshion list
zeshion list --tmux
zeshion list --zellij
zeshion list --herdr
zeshion list --config --zoxide
zeshion connect my-session
zeshion connect --root "$PWD"
zeshion picker
zeshion preview my-session
zeshion last
zeshion window
zeshion window ~/projects/my-app
```

Aliases are preserved where practical:

| Command | Aliases |
| --- | --- |
| `list` | `l` |
| `connect` | `cn` |
| `clone` | `cl` |
| `picker` | `pick`, `pk` |
| `preview` | `p` |
| `root` | `r` |
| `last` | `L` |
| `window` | `w` |

## Multiplexer Selection

`zeshion` defaults to `--multiplexer auto`.

Auto mode chooses:

1. Herdr when running inside a Herdr pane.
2. Zellij when running inside Zellij.
3. tmux when running inside tmux.
4. The `multiplexer` value from config.
5. tmux as the compatibility fallback.

Override the multiplexer per command:

```sh
zeshion --multiplexer tmux list
zeshion --multiplexer zellij list
zeshion --multiplexer herdr list
zeshion --multiplexer zellij connect my-session
```

## Multiplexer Behavior

| Action | tmux | Zellij | Herdr |
| --- | --- | --- | --- |
| zeshion session | Session | Session | Workspace |
| zeshion window | Window | Tab | Tab |
| List live sessions | `tmux list-sessions` | `zellij list-sessions` | `herdr api snapshot` |
| Connect outside a session | Attach | Attach in the current terminal | Focus the workspace, then attach the Herdr UI |
| Connect inside a session | Switch client | Switch session | Focus workspace |
| Create session | New session | Attach/create flow | Create workspace and run startup command in its root pane |
| tmuxinator/tmuxp | Supported | Not used | Not used |

Herdr workspace and tab IDs are retained internally. The picker therefore selects the exact workspace even when multiple workspaces share a label. A direct `zeshion connect <name>` requires the label or derived directory name to be unique; use the workspace ID from `zeshion list --herdr --json` otherwise. `zeshion window` likewise accepts a tab ID from its JSON output and reports duplicate tab labels instead of selecting one arbitrarily.

`zeshion last` is unavailable for Herdr because Herdr 0.8 does not expose previous-workspace history. Herdr panes can be listed, but Herdr 0.8 does not expose direct pane-ID focus through its CLI, so selecting a pane is not supported.

## Picker

Use the interactive picker:

```sh
zeshion picker
```

Filter the picker input sources:

```sh
zeshion picker --tmux
zeshion picker --zellij
zeshion picker --herdr
zeshion picker --config --zoxide
zeshion picker --hide-duplicates
```

Every picker row includes its source, such as `[tmux]`, `[zellij]`, or `[herdr]`.

Customize picker text:

```sh
zeshion picker --prompt "session> " --placeholder "Find session"
```

## Configuration

`zeshion` looks for config files in this order:

1. The path passed with `--config`.
2. `~/.config/zeshion/zeshion.toml`.
3. Legacy `~/.config/zesh/zesh.toml`.
4. Legacy `~/.config/sesh/sesh.toml`.
5. Built-in defaults when no config exists.

Minimal config:

```toml
multiplexer = "auto"
tmux_command = "tmux"
zellij_command = "zellij"
herdr_command = "herdr"
# Optional named Herdr server session used when outside Herdr.
herdr_session = ""

[default_session]
startup_command = "nvim"
preview_command = "ls -la {}"

[[session]]
name = "dotfiles"
path = "~/dotfiles"
startup_command = "nvim"
```

Use the JSON schema from `zeshion.schema.json`.

## Migration

From `sesh`:

- Keep your existing `~/.config/sesh/sesh.toml`; `zeshion` will read it if no `zeshion` or `zesh` config exists.
- Add `multiplexer = "auto"`, `"tmux"`, `"zellij"`, or `"herdr"` when you want the choice to be explicit.
- tmux-specific features such as tmuxinator stay on the tmux path.
- Set `zellij_command` or `herdr_command` if you use a wrapper or custom binary name.

From pre-release `zesh`:

- Rename `~/.config/zesh/zesh.toml` to `~/.config/zeshion/zeshion.toml` when convenient.
- The old `zesh` config path remains a fallback for local migration.
- The binary is now `zeshion`.

## Troubleshooting

**Selecting a Zellij session from a normal terminal should open the Zellij UI.**

If it only sends actions to an existing Zellij pane, make sure you are running an up-to-date `zeshion` binary and not an old `zesh` binary. Check:

```sh
command -v zeshion
zeshion --version
```

**`--multiplexer auto` chooses the wrong multiplexer.**

Use an explicit override while debugging:

```sh
zeshion --multiplexer tmux list
zeshion --multiplexer zellij list
zeshion --multiplexer herdr list
```

Then set `multiplexer = "tmux"`, `"zellij"`, or `"herdr"` in config if you do not want auto-selection.

**No config is found.**

Run with an explicit config path:

```sh
zeshion --config ~/.config/zeshion/zeshion.toml list
```

## Development

```sh
just mock
go test ./...
go test -race ./...
go vet ./...
go build -o zeshion .
```

Generate the man page:

```sh
just man
```

Manual release checks live in [`docs/release-smoke-test.md`](docs/release-smoke-test.md).

Version-control work in this repo uses Jujutsu (`jj`).

## Attribution

`zeshion` is based on [`joshmedeski/sesh`](https://github.com/joshmedeski/sesh), which is distributed under the MIT License. The original copyright notice is preserved in `LICENSE`; the upstream derivation is also recorded in `NOTICE`.
