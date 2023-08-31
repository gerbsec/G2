package agents

type AgentTask struct {
	Id        string
	Command   string
	Arguments []string
	File      []byte
}
