# Release Smoke Test

Run this checklist before tagging a public `zeshion` release.

## Build Checks

```sh
go test ./...
go test -race ./...
go vet ./...
go build -o /tmp/zeshion .
/tmp/zeshion --version
/tmp/zeshion --help
```

## Generated Shell Output

```sh
/tmp/zeshion completion zsh >/tmp/zeshion.zsh
/tmp/zeshion completion bash >/tmp/zeshion.bash
/tmp/zeshion completion fish >/tmp/zeshion.fish
/tmp/zeshion man >/tmp/zeshion.1
```

Confirm each generated file names `zeshion`, not an old pre-release name.

## Config Lookup

Check explicit config:

```sh
/tmp/zeshion --config ~/.config/zeshion/zeshion.toml list
```

Check fallback order manually when those files exist:

1. `~/.config/zeshion/zeshion.toml`
2. `~/.config/zesh/zesh.toml`
3. `~/.config/sesh/sesh.toml`

## tmux Smoke

From outside tmux:

```sh
/tmp/zeshion --multiplexer tmux list --tmux
/tmp/zeshion --multiplexer tmux connect <tmux-session>
```

From inside tmux:

```sh
/tmp/zeshion --multiplexer tmux picker --tmux
/tmp/zeshion --multiplexer tmux last
/tmp/zeshion --multiplexer tmux window
```

Expected behavior:

- Outside tmux, selecting an existing tmux session attaches to it.
- Inside tmux, selecting an existing tmux session switches the client.
- `window` targets tmux windows.

## Zellij Smoke

From outside Zellij in a normal terminal:

```sh
/tmp/zeshion --multiplexer zellij list --zellij
/tmp/zeshion --multiplexer zellij picker --zellij
```

Select an existing Zellij session in the picker. It should open the Zellij UI in the current terminal.

From inside Zellij:

```sh
/tmp/zeshion --multiplexer zellij picker --zellij
/tmp/zeshion --multiplexer zellij last
/tmp/zeshion --multiplexer zellij window
```

Expected behavior:

- Outside Zellij, selecting an existing Zellij session attaches to it in the current terminal.
- Inside Zellij, selecting an existing Zellij session switches the active session.
- `window` targets Zellij tabs.

## Picker Sources

```sh
/tmp/zeshion picker --tmux
/tmp/zeshion picker --zellij
/tmp/zeshion picker --config
/tmp/zeshion picker --zoxide
/tmp/zeshion picker --config --zoxide --hide-duplicates
```

Confirm source filtering works and picker selection connects to the expected session or directory.

## Migration Checks

Confirm the binary is named `zeshion`:

```sh
command -v zeshion
zeshion --version
```

Confirm old local names do not appear in generated public artifacts except documented migration paths:

```sh
rg -n "zesh\\b|Zesh\\b" README.md docs zeshion.schema.json .github justfile
```
