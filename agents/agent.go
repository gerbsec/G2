package agents

import (
	"sync"
	"time"
)

type Agent struct {
	Metadata     *AgentMetadata
	LastSeen     time.Time
	pendingTasks []*AgentTask
	taskResults  []*AgentTaskResult
	mutex        sync.Mutex
}

func NewAgent(metadata *AgentMetadata) *Agent {
	return &Agent{
		Metadata:     metadata,
		pendingTasks: make([]*AgentTask, 0),
		taskResults:  make([]*AgentTaskResult, 0),
	}
}

func (a *Agent) CheckIn() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.LastSeen = time.Now().UTC()
}

func (a *Agent) QueueTask(task *AgentTask) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.pendingTasks = append(a.pendingTasks, task)
}

func (a *Agent) GetPendingTasks() []*AgentTask {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	tasks := a.pendingTasks
	a.pendingTasks = []*AgentTask{} // reset the pending tasks after fetching
	return tasks
}

func (a *Agent) GetTaskResult(taskId string) *AgentTaskResult {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for _, result := range a.taskResults {
		if result.Id == taskId {
			return result
		}
	}
	return nil
}

func (a *Agent) GetTaskResults() []*AgentTaskResult {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.taskResults
}
