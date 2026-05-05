package zeshcli

import (
	"testing"

	"github.com/notes/zesh/oswrap"
	"github.com/stretchr/testify/assert"
)

func TestSelectMuxNameAutoPrefersEnvironment(t *testing.T) {
	mockOs := new(oswrap.MockOs)
	mockOs.EXPECT().Getenv("ZELLIJ").Return("1")

	assert.Equal(t, "zellij", selectMuxName(mockOs, "auto", "tmux"))
}

func TestSelectMuxNameAutoFallsBackToConfigThenTmux(t *testing.T) {
	mockOs := new(oswrap.MockOs)
	mockOs.EXPECT().Getenv("ZELLIJ").Return("")
	mockOs.EXPECT().Getenv("ZELLIJ_SESSION_NAME").Return("")
	mockOs.EXPECT().Getenv("TMUX").Return("")

	assert.Equal(t, "zellij", selectMuxName(mockOs, "auto", "zellij"))

	mockOs = new(oswrap.MockOs)
	mockOs.EXPECT().Getenv("ZELLIJ").Return("")
	mockOs.EXPECT().Getenv("ZELLIJ_SESSION_NAME").Return("")
	mockOs.EXPECT().Getenv("TMUX").Return("")

	assert.Equal(t, "tmux", selectMuxName(mockOs, "auto", ""))
}

func TestSelectMuxNameOverrideWins(t *testing.T) {
	mockOs := new(oswrap.MockOs)

	assert.Equal(t, "zellij", selectMuxName(mockOs, "zellij", "tmux"))
	assert.Equal(t, "tmux", selectMuxName(mockOs, "tmux", "zellij"))
}
