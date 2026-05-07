# zeshion

`zeshion` is a session manager for tmux and Zellij. It keeps the `sesh` command shape where practical, while making the selected multiplexer explicit and supporting both tmux sessions and Zellij sessions from one CLI.

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
- `zoxide` for zoxide-backed directory sessions.
- A picker-compatible terminal for `zeshion picker`.

## Commands

Most commands follow the `sesh` shape:

```sh
zeshion list
zeshion list --tmux
zeshion list --zellij
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

1. Zellij when running inside Zellij.
2. tmux when running inside tmux.
3. The `multiplexer` value from config.
4. tmux as the compatibility fallback.

Override the multiplexer per command:

```sh
zeshion --multiplexer tmux list
zeshion --multiplexer zellij list
zeshion --multiplexer zellij connect my-session
```

## tmux And Zellij Behavior

| Action | tmux | Zellij |
| --- | --- | --- |
| List live sessions | `tmux list-sessions` | `zellij list-sessions` |
| Connect outside a session | Attach to the tmux session | Attach to the Zellij session in the current terminal |
| Connect inside a session | Switch tmux client | Switch Zellij session |
| Create session | New tmux session | Zellij attach/create flow |
| Window command | tmux windows | Zellij tabs |
| tmuxinator/tmuxp | tmux-only path | Not used for native Zellij sessions |

For Zellij, the important rule is simple: outside Zellij, selecting a Zellij session opens that session in the current terminal; inside Zellij, selecting a Zellij session switches the active Zellij session.

## Picker

Use the interactive picker:

```sh
zeshion picker
```

Filter the picker input sources:

```sh
zeshion picker --tmux
zeshion picker --zellij
zeshion picker --config --zoxide
zeshion picker --hide-duplicates
```

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
- Add `multiplexer = "auto"` or `multiplexer = "tmux"` if you want the choice to be explicit.
- tmux-specific features such as tmuxinator stay on the tmux path.
- Add `zellij_command = "zellij"` if you use a wrapper or custom Zellij binary name.

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
```

Then set `multiplexer = "tmux"` or `multiplexer = "zellij"` in config if you do not want auto-selection.

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
