package models

type AgentCommand interface {
	Name() string
	Execute(task *AgentTask) string
}
