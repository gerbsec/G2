package commands

import (
	"os"

	"github.com/gerbsec/D2/agent/models"
)

var _ models.AgentCommand = (*ChangeDirectory)(nil)

type ChangeDirectory struct{}

func (t *ChangeDirectory) Name() string {
	return "cd"
}

func (t *ChangeDirectory) Execute(task *models.AgentTask) string {
	var path string
	if task.Arguments == nil || len(task.Arguments) == 0 {
		path, _ = os.UserHomeDir()
	} else {
		path = task.Arguments[0]
	}
	os.Chdir(path)
	wd, _ := os.Getwd()
	return wd
}
