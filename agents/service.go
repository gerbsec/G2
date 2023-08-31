package agents

import (
	"sync"
)

type AgentService interface {
	AddAgent(agent *Agent)
	GetAgents() []*Agent
	GetAgent(id string) *Agent
	RemoveAgent(agent *Agent)
}

type Service struct {
	agents []*Agent
	mux    sync.Mutex
}

func NewService() AgentService {
	return &Service{
		agents: make([]*Agent, 0),
	}
}

func (s *Service) AddAgent(agent *Agent) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.agents = append(s.agents, agent)
}

func (s *Service) GetAgents() []*Agent {
	return s.agents
}

func (s *Service) GetAgent(id string) *Agent {
	for _, a := range s.agents {
		if a.Metadata.Id == id {
			return a
		}
	}
	return nil
}

func (s *Service) RemoveAgent(agent *Agent) {
	s.mux.Lock()
	defer s.mux.Unlock()
	for i, a := range s.agents {
		if a == agent {
			s.agents = append(s.agents[:i], s.agents[i+1:]...)
			return
		}
	}
}
