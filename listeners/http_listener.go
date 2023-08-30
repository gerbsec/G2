package http_listener

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type HttpListener struct {
	Name     string `json:"name"`
	BindPort string `json:"bind_port"`
	server   *http.Server
	stopChan chan bool
	wg       *sync.WaitGroup
}

var listenersMap = make(map[string]*HttpListener)

func (s *HttpListener) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from listener %s!\n", s.Name)
}

func (s *HttpListener) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handler)

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
	if _, exists := listenersMap[name]; exists {
		return fmt.Errorf("a listener with the name %s already exists", name)
	}

	l := &HttpListener{
		Name:     name,
		BindPort: bindPort,
		stopChan: make(chan bool),
		wg:       &sync.WaitGroup{},
	}
	l.wg.Add(1)
	listenersMap[name] = l
	l.Start()
	return nil
}

func GetListenerInfoByName(name string) (string, error) {
	if listener, exists := listenersMap[name]; exists {
		data, err := json.Marshal(listener)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", fmt.Errorf("No listener found with name: %s", name)
}

func GetAllListenersInfo() string {
	var infos []*HttpListener
	for _, l := range listenersMap {
		infos = append(infos, l)
	}

	data, _ := json.Marshal(infos)
	return string(data)
}

func StopListenerByName(name string) {
	if listener, exists := listenersMap[name]; exists {
		listener.Stop()
		listener.wg.Wait()
		delete(listenersMap, name)
	}
}
