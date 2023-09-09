package commands

import (
	"fmt"
	"strings"

	"github.com/gerbsec/G2/agent/models"
	"github.com/shirou/gopsutil/process"
)

var _ models.AgentCommand = (*ListProcesses)(nil)

type ListProcesses struct{}

func (l *ListProcesses) Name() string {
	return "ps"
}

func (l *ListProcesses) Execute(task *models.AgentTask) string {
	processes, err := process.Processes()
	if err != nil {
		return fmt.Sprintf("Failed to get processes: %s", err)
	}

	var processList []string
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		pid := p.Pid // get the PID of the process

		processInfo := fmt.Sprintf("%d: %s", pid, name)
		processList = append(processList, processInfo)
	}

	return fmt.Sprintf("Processes:\n%s", strings.Join(processList, "\n"))
}
