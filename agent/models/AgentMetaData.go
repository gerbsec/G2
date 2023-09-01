package models

type AgentMetadata struct {
	Id           string
	Hostname     string
	Username     string
	ProcessName  string
	ProcessId    int
	Integrity    string
	Architecture string
}
