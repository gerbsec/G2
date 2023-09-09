package commands

import (
	"os"

	"github.com/gerbsec/G2/agent/models"
)

var _ models.AgentCommand = (*PrintWorkingDirectory)(nil)

type PrintWorkingDirectory struct{}

func (t *PrintWorkingDirectory) Name() string {
	return "pwd"
}

func (t *PrintWorkingDirectory) Execute(task *models.AgentTask) string {
	wd, _ := os.Getwd()
	return wd
}
