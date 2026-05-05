package tmux

func (t *RealTmux) SelectWindow(targetWindow string) (string, error) {
	return t.shell.Cmd(t.command(), "select-window", "-t", targetWindow)
}
