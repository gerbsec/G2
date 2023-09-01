package models

type CommModule interface {
	Start() chan bool
	Stop()
	Init(metadata *AgentMetadata)
	RecvData() ([]*AgentTask, bool)
	SendData(result *AgentTaskResult)
	getOutbound() []*AgentTaskResult
}

type BaseCommModule struct {
	AgentMetadata *AgentMetadata
	Inbound       chan *AgentTask
	Outbound      chan *AgentTaskResult
}
