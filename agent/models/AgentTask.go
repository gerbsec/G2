package models

type AgentTask struct {
	Id        string   `json:"id"`
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
	File      []byte   `json:"file"`
}
