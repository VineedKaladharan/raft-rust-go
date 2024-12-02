// This program simulates a basic leader election process, like choosing a team captain!
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Message is like a letter that nodes send to each other
// It contains who the leader is and what kind of message it is
type Message struct {
	leaderID int
	msgType  string
}

// Subscriber is like a mailbox for each node
// Each node has an ID and a channel (msgs) to receive messages
type Subscriber struct {
	id   int
	msgs chan Message
}

// Fanout manages message distribution to multiple subscribers
type Fanout struct {
	mu          sync.RWMutex
	subscribers map[int]chan Message
	input       chan Message
	done        chan struct{}
}

// NewFanout creates a new Fanout system
func NewFanout() *Fanout {
	f := &Fanout{
		subscribers: make(map[int]chan Message),
		input:       make(chan Message),
		done:        make(chan struct{}),
	}
	go f.run()
	return f
}

// run handles the message distribution
func (f *Fanout) run() {
	for {
		select {
		case msg := <-f.input:
			f.mu.RLock()
			for _, ch := range f.subscribers {
				select {
				case ch <- msg:
				default:
					// Skip if channel is full
				}
			}
			f.mu.RUnlock()
		case <-f.done:
			return
		}
	}
}

// Subscribe adds a new subscriber
func (f *Fanout) Subscribe(id int) chan Message {
	f.mu.Lock()
	defer f.mu.Unlock()
	ch := make(chan Message, 10)
	f.subscribers[id] = ch
	return ch
}

// Unsubscribe removes a subscriber
func (f *Fanout) Unsubscribe(id int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if ch, exists := f.subscribers[id]; exists {
		close(ch)
		delete(f.subscribers, id)
	}
}

// Publish sends a message to all subscribers
func (f *Fanout) Publish(msg Message) {
	select {
	case f.input <- msg:
	case <-f.done:
	}
}

// Close shuts down the Fanout system
func (f *Fanout) Close() {
	close(f.done)
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, ch := range f.subscribers {
		close(ch)
	}
	f.subscribers = make(map[int]chan Message)
}

// State keeps track of whether we have chosen a leader
// It's like a bulletin board that everyone can check
type State struct {
	mu           sync.Mutex // This lock makes sure only one person can update the bulletin board at a time
	leaderChosen bool       // Has a leader been chosen?
	leaderID     int        // Who is the current leader?
}

// The main function is where everything starts!
func main() {
	// Create our tools: post office, bulletin board, and random number generator
	const numGoroutines = 5
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	state := &State{}
	fanout := NewFanout()
	defer fanout.Close()

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch goroutines
	for i := 1; i <= numGoroutines; i++ {
		go runNodeElection(i, state, fanout, rng, &wg)
	}

	// Let it run for a while
	time.Sleep(10 * time.Second)
	fanout.Close()
	wg.Wait()
}

// runNodeElection is like each person participating in the election
func runNodeElection(id int, state *State, fanout *Fanout, rng *rand.Rand, wg *sync.WaitGroup) {
	defer wg.Done()

	// Get a mailbox to receive messages
	msgChan := fanout.Subscribe(id)
	defer fanout.Unsubscribe(id)

	// Each node waits a random time before trying to become leader
	timeout := time.Duration(rng.Intn(1000)+500) * time.Millisecond
	time.Sleep(timeout)

	state.mu.Lock()
	if !state.leaderChosen {
		// If no leader yet, this node becomes leader
		state.leaderChosen = true
		state.leaderID = id
		state.mu.Unlock()

		fmt.Printf("Goroutine %d became the leader!\n", id)
		fanout.Publish(Message{leaderID: id, msgType: "leader"})
		runLeaderHeartbeat(id, fanout)
		return
	}
	state.mu.Unlock()

	fmt.Printf("Goroutine %d starting as follower\n", id)
	runFollowerMessageHandling(id, msgChan, fanout)
}

// runLeaderHeartbeat is like the leader saying "I'm still here!" every second
func runLeaderHeartbeat(id int, fanout *Fanout) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fanout.Publish(Message{leaderID: id, msgType: "heartbeat"})
			fmt.Printf("Leader %d sending heartbeat\n", id)
		case <-fanout.done:
			return
		}
	}
}

// runFollowerMessageHandling is like followers listening for messages from the leader
func runFollowerMessageHandling(id int, msgChan <-chan Message, fanout *Fanout) {
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return
			}
			handleFollowerMessage(id, msg)
		case <-fanout.done:
			return
		}
	}
}

// handleFollowerMessage is like a follower processing a message they received
func handleFollowerMessage(id int, msg Message) {
	switch msg.msgType {
	case "leader":
		fmt.Printf("Follower %d acknowledging leader %d\n", id, msg.leaderID)
	case "heartbeat":
		fmt.Printf("Follower %d received heartbeat from leader %d\n", id, msg.leaderID)
	}
}
