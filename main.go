package main

import (
	"fmt"
	"sync"
)

// Agent is a simple pub/sub agent
type Agent struct {
	mu    sync.Mutex
	subs  map[string][]chan string
	quit  chan struct{}
	closed bool
}

// NewAgent creates a new Agent
func NewAgent() *Agent {
	return &Agent{
		subs: make(map[string][]chan string),
		quit: make(chan struct{}),
	}
}

// Publish publishes a message to a topic
func (b *Agent) Publish(topic string, msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	for _, ch := range b.subs[topic] {
		ch <- msg
	}
}

// Subscribe subscribes to a topic
func (b *Agent) Subscribe(topic string) <-chan string {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	ch := make(chan string)
	b.subs[topic] = append(b.subs[topic], ch)
	return ch
}

// Close closes the agent
func (b *Agent) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	b.closed = true
	close(b.quit)

	for _, ch := range b.subs {
		for _, sub := range ch {
			close(sub)
		}
	}
}

func main() {
	// Create a new agent
	agent := NewAgent()

	// Subscribe to a topic
	sub := agent.Subscribe("foo")

	// Publish a message to the topic
	go agent.Publish("foo", "hello world")

	// Print the message
	fmt.Println(<-sub)

	// Close the agent
	agent.Close()
}
