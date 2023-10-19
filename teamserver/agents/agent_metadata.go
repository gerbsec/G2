package agents

type AgentMetadata struct {
	Id           string `json:"id"`
	Hostname     string `json:"hostname"`
	Username     string `json:"username"`
	Ip           string `json:"ip"`
	ProcessName  string `json:"processName"`
	ProcessId    int    `json:"processId"`
	Architecture string `json:"architecture"`
}
