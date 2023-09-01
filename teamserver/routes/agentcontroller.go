package routes

import (
	"net/http"

	"github.com/gerbsec/D2/teamserver/agents"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var agentService = agents.NewService()

func SetupAgentRoutes(r *gin.Engine) {
	r.GET("/Agents", getAgents)
	r.GET("/Agents/:agentId", getAgent)
	r.GET("/Agents/:agentId/tasks", getTaskResults)
	r.GET("/Agents/:agentId/tasks/:taskId", getTaskResult)
	r.POST("/Agents/:agentId", taskAgent)
}

func getAgents(c *gin.Context) {
	allAgents := agentService.GetAgents()
	c.JSON(http.StatusOK, allAgents)
}

func getAgent(c *gin.Context) {
	agentId := c.Param("agentId")
	agent := agentService.GetAgent(agentId)
	if agent == nil {
		c.String(http.StatusNotFound, "Agent not found")
		return
	}
	c.JSON(http.StatusOK, agent)
}

func getTaskResults(c *gin.Context) {
	agentId := c.Param("agentId")
	agent := agentService.GetAgent(agentId)
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
	agent := agentService.GetAgent(agentId)
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

	agent := agentService.GetAgent(agentId)
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
