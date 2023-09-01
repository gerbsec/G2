package commands

import "github.com/gerbsec/D2/agent/models"

var _ models.AgentCommand = (*TestCommand)(nil)

type TestCommand struct{}

func (t *TestCommand) Name() string {
	return "TestCommand"
}

func (t *TestCommand) Execute(task *models.AgentTask) string {
	return "Hello from Test Command"
}
