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

// PubSub is like a post office that helps nodes send messages to each other
// It keeps track of all the mailboxes (subscribers) and makes sure messages get delivered
type PubSub struct {
	mu          sync.RWMutex    // This is like a lock to make sure only one person uses the mailroom at a time
	subscribers map[int]chan Message  // This maps each node's ID to their mailbox
	done        chan struct{}        // This is used to signal when the post office is closed
}

// NewPubSub creates a new post office (PubSub system)
func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[int]chan Message),
		done:        make(chan struct{}),
	}
}

// Subscribe is like signing up to receive mail
// It creates a new mailbox for a node with the given ID
func (ps *PubSub) Subscribe(id int) chan Message {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ch := make(chan Message, 10)
	ps.subscribers[id] = ch
	return ch
}

// Unsubscribe is like closing your mailbox
// It removes a node's mailbox when they don't want to receive messages anymore
func (ps *PubSub) Unsubscribe(id int) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if ch, exists := ps.subscribers[id]; exists {
		close(ch)
		delete(ps.subscribers, id)
	}
}

// Publish is like sending a letter to everyone
// When someone sends a message, it gets copied and delivered to all mailboxes
func (ps *PubSub) Publish(msg Message) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	for id, ch := range ps.subscribers {
		if id != msg.leaderID { // Don't send to self
			select {
			case ch <- msg:
			default:
				// If a mailbox is full, we skip it
			}
		}
	}
}

// Close is like shutting down the post office
// It closes all mailboxes and cleans up
func (ps *PubSub) Close() {
	close(ps.done)
	ps.mu.Lock()
	defer ps.mu.Unlock()
	for _, ch := range ps.subscribers {
		close(ch)
	}
	ps.subscribers = make(map[int]chan Message)
}

// State keeps track of whether we have chosen a leader
// It's like a bulletin board that everyone can check
type State struct {
	mu           sync.Mutex  // This lock makes sure only one person can update the bulletin board at a time
	leaderChosen bool       // Has a leader been chosen?
	leaderID     int        // Who is the current leader?
}

// The main function is where everything starts!
func main() {
	// Create our tools: post office, bulletin board, and random number generator
	const numGoroutines = 5
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	state := &State{}
	pubsub := NewPubSub()
	defer pubsub.Close()

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Launch goroutines
	for i := 1; i <= numGoroutines; i++ {
		go runNodeElection(i, state, pubsub, rng, &wg)
	}

	// Let it run for a while
	time.Sleep(10 * time.Second)
	pubsub.Close()
	wg.Wait()
}

// runNodeElection is like each person participating in the election
func runNodeElection(id int, state *State, pubsub *PubSub, rng *rand.Rand, wg *sync.WaitGroup) {
	defer wg.Done()

	// Get a mailbox to receive messages
	msgChan := pubsub.Subscribe(id)
	defer pubsub.Unsubscribe(id)

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
		pubsub.Publish(Message{leaderID: id, msgType: "leader"})
		runLeaderHeartbeat(id, pubsub)
		return
	}
	state.mu.Unlock()

	fmt.Printf("Goroutine %d starting as follower\n", id)
	runFollowerMessageHandling(id, msgChan, pubsub)
}

// runLeaderHeartbeat is like the leader saying "I'm still here!" every second
func runLeaderHeartbeat(id int, pubsub *PubSub) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pubsub.Publish(Message{leaderID: id, msgType: "heartbeat"})
			fmt.Printf("Leader %d sending heartbeat\n", id)
		case <-pubsub.done:
			return
		}
	}
}

// runFollowerMessageHandling is like followers listening for messages from the leader
func runFollowerMessageHandling(id int, msgChan <-chan Message, pubsub *PubSub) {
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return
			}
			handleFollowerMessage(id, msg)
		case <-pubsub.done:
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
