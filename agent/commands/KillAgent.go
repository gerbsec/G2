package commands

import (
	"os"

	"github.com/gerbsec/G2/agent/models"
)

var _ models.AgentCommand = (*KillAgent)(nil)

type KillAgent struct{}

func (t *KillAgent) Name() string {
	return "KillAgent"
}

func (t *KillAgent) Execute(task *models.AgentTask) string {
	os.Exit(1)
	return ""
}
