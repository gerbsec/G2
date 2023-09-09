package commands

import (
	"os/exec"
	"runtime"

	"github.com/gerbsec/G2/agent/models"
)

var _ models.AgentCommand = (*ListDirectory)(nil)

type ListDirectory struct{}

func (t *ListDirectory) Name() string {
	return "ls"
}

func (t *ListDirectory) Execute(task *models.AgentTask) string {
	var cmd []string
	if runtime.GOOS == "windows" {
		cmd = []string{"powershell.exe", "/c"}
	} else {
		cmd = []string{"sh", "-c"}
	}
	command := exec.Command(cmd[0], cmd[1], "ls")
	output, err := command.Output()
	if err != nil {
		return "Error getting contents"
	}
	return string(output)
}
