// Package server implements the core Raft consensus protocol
package server

import (
	"fmt"       // Package for formatted I/O
	"log"       // Package for logging
	"math/rand" // Package for random number generation
	"net"       // Package for network I/O
	"net/rpc"   // Package for RPC implementation
	"sync"      // Package for synchronization primitives
	"time"      // Package for time-related functions
)

const (
	minTimeout = 150 // Minimum election timeout in milliseconds
	maxTimeout = 300 // Maximum election timeout in milliseconds
)

// ServerState represents the current state of a Raft server
type ServerState int

const (
	Follower  ServerState = iota // Server is in follower state
	Candidate                    // Server is in candidate state
	Leader                       // Server is in leader state
)

var (
	// leaderElected is a global channel to ensure only one leader is elected
	leaderElected = make(chan struct{}, 1)
	// globalMutex protects state changes across all servers
	globalMutex = sync.Mutex{}
)

// String returns the string representation of ServerState
func (s ServerState) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	default:
		return "Unknown"
	}
}

// RaftServer represents a single node in the Raft cluster
type RaftServer struct {
	ID            int           // Unique identifier for the server
	Port          int           // Port number the server listens on
	state         ServerState   // Current state of the server
	currentTerm   int          // Latest term server has seen
	votedFor      int          // CandidateId that received vote in current term
	electionTimer *time.Timer  // Timer for election timeout
	listener      net.Listener // Network listener for RPC
	server        *rpc.Server  // RPC server instance
	mu            sync.Mutex   // Protects server state
	stopChan      chan struct{} // Channel for coordinating shutdown
}

// RaftRPC holds RPC methods for the Raft protocol
type RaftRPC struct {
	server *RaftServer // Reference to parent server
}

// RequestVoteArgs contains arguments for the RequestVote RPC
type RequestVoteArgs struct {
	Term         int // Candidate's term
	CandidateID  int // Candidate requesting vote
	LastLogIndex int // Index of candidate's last log entry
	LastLogTerm  int // Term of candidate's last log entry
}

// RequestVoteReply contains the response for the RequestVote RPC
type RequestVoteReply struct {
	Term        int  // Current term, for candidate to update itself
	VoteGranted bool // True means candidate received vote
}

// RequestVote handles vote requests from candidates
func (r *RaftRPC) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	r.server.mu.Lock()
	defer r.server.mu.Unlock()
	// TODO: Implement vote request logic
	return nil
}

// getRandomTimeout returns a random duration between minTimeout and maxTimeout
func getRandomTimeout() time.Duration {
	return time.Duration(rand.Intn(maxTimeout-minTimeout)+minTimeout) * time.Millisecond
}

// becomeLeader attempts to transition the server to leader state
func (s *RaftServer) becomeLeader() {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	select {
	case leaderElected <- struct{}{}: // Try to signal leader election
		s.state = Leader
		log.Printf("Server %d: Yay, I am first! Becoming Leader", s.ID)
		if s.electionTimer != nil {
			s.electionTimer.Stop()
		}
	default: // Channel is full, someone else is already leader
		s.state = Follower
		log.Printf("Server %d: Another leader was elected, becoming Follower", s.ID)
	}
}

// resetElectionTimer resets the election timer with a random timeout
func (s *RaftServer) resetElectionTimer() {
	if s.electionTimer != nil {
		s.electionTimer.Stop()
	}
	timeout := getRandomTimeout()
	s.electionTimer = time.AfterFunc(timeout, func() {
		select {
		case <-s.stopChan: // Check if server is stopping
			return
		default:
			s.becomeLeader()
		}
	})
}

// NewRaftServer creates a new Raft server instance
func NewRaftServer(id int, port int) (*RaftServer, error) {
	rand.Seed(time.Now().UnixNano() + int64(id)) // Ensure different random sequences
	
	server := &RaftServer{
		ID:       id,
		Port:     port,
		state:    Follower,
		stopChan: make(chan struct{}),
	}
	
	// Create and register RPC server
	rpcServer := rpc.NewServer()
	rpcRaft := &RaftRPC{server: server}
	err := rpcServer.RegisterName("Raft", rpcRaft)
	if err != nil {
		return nil, fmt.Errorf("failed to register RPC server: %v", err)
	}
	
	server.server = rpcServer
	return server, nil
}

// Start starts the RPC server and election timer
func (s *RaftServer) Start() error {
	addr := fmt.Sprintf(":%d", s.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start listener on port %d: %v", s.Port, err)
	}
	
	s.listener = listener
	log.Printf("Server %d started on port %d as %s", s.ID, s.Port, s.state)
	
	go s.server.Accept(listener)
	
	// Start election timer
	s.resetElectionTimer()
	
	return nil
}

// Stop stops the RPC server and election timer
func (s *RaftServer) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	close(s.stopChan)
	
	if s.electionTimer != nil {
		s.electionTimer.Stop()
	}
	
	if s.listener != nil {
		s.listener.Close()
	}
	
	// Clear leader election channel if this server was the leader
	if s.state == Leader {
		select {
		case <-leaderElected: // Clear the channel
		default:
		}
	}
}
