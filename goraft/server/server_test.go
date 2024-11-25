// Package server_test contains tests for the Raft server implementation
package server

import (
	"fmt"
	"net/rpc"
	"sync"
	"testing"
	"time"
)

// TestServerState_String tests the string representation of server states
func TestServerState_String(t *testing.T) {
	// Test cases for each server state
	tests := []struct {
		state ServerState
		want  string
	}{
		{Follower, "Follower"},
		{Candidate, "Candidate"},
		{Leader, "Leader"},
		{ServerState(99), "Unknown"}, // Test unknown state
	}

	// Run each test case
	for _, tt := range tests {
		got := tt.state.String()
		if got != tt.want {
			t.Errorf("ServerState.String() = %v, want %v", got, tt.want)
		}
	}
}

// TestNewRaftServer tests the creation of a new Raft server
func TestNewRaftServer(t *testing.T) {
	// Test cases for server creation
	tests := []struct {
		name    string
		id      int
		port    int
		wantErr bool
	}{
		{
			name:    "Valid server creation",
			id:      1,
			port:    51000,
			wantErr: false,
		},
		{
			name:    "Valid server with different port",
			id:      2,
			port:    51001,
			wantErr: false,
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create server
			server, err := NewRaftServer(tt.id, tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRaftServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if server != nil {
				if server.ID != tt.id {
					t.Errorf("NewRaftServer() server.ID = %v, want %v", server.ID, tt.id)
				}
				if server.Port != tt.port {
					t.Errorf("NewRaftServer() server.Port = %v, want %v", server.Port, tt.port)
				}
				if server.state != Follower {
					t.Errorf("NewRaftServer() server.state = %v, want Follower", server.state)
				}
			}
		})
	}
}

// TestLeaderElection tests the leader election process
func TestLeaderElection(t *testing.T) {
	// Clear the leader election channel before test
	select {
	case <-leaderElected:
	default:
	}

	// Create multiple servers for testing
	numServers := 3
	servers := make([]*RaftServer, numServers)
	ports := []int{51100, 51101, 51102}
	var wg sync.WaitGroup

	// Initialize and start each server
	for i := 0; i < numServers; i++ {
		server, err := NewRaftServer(i+1, ports[i])
		if err != nil {
			t.Fatalf("Failed to create server %d: %v", i+1, err)
		}
		servers[i] = server
		wg.Add(1)
		go func(s *RaftServer) {
			defer wg.Done()
			err := s.Start()
			if err != nil {
				t.Errorf("Failed to start server: %v", err)
			}
		}(server)
	}

	// Wait for leader election to complete
	time.Sleep(500 * time.Millisecond)

	// Verify only one leader is elected
	leaderCount := 0
	var leaderId int
	for _, s := range servers {
		if s.state == Leader {
			leaderCount++
			leaderId = s.ID
		}
	}

	// Check leader count
	if leaderCount != 1 {
		t.Errorf("Expected exactly one leader, got %d", leaderCount)
	}

	// Verify others are followers
	for _, s := range servers {
		if s.ID != leaderId && s.state != Follower {
			t.Errorf("Server %d should be Follower, got %s", s.ID, s.state)
		}
	}

	// Clean up servers
	for _, s := range servers {
		s.Stop()
	}
	wg.Wait()
}

// TestRaftServerStartStop tests server startup and shutdown
func TestRaftServerStartStop(t *testing.T) {
	// Test cases for server start/stop
	tests := []struct {
		name    string
		id      int
		port    int
		wantErr bool
	}{
		{
			name:    "Valid server creation",
			id:      1,
			port:    51000,
			wantErr: false,
		},
		{
			name:    "Valid server with different port",
			id:      2,
			port:    51001,
			wantErr: false,
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create server
			server, err := NewRaftServer(tt.id, tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRaftServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if server != nil {
				if server.ID != tt.id {
					t.Errorf("NewRaftServer() server.ID = %v, want %v", server.ID, tt.id)
				}
				if server.Port != tt.port {
					t.Errorf("NewRaftServer() server.Port = %v, want %v", server.Port, tt.port)
				}
			}
		})
	}
}

// TestRaftServer_StartStop tests server startup and shutdown
func TestRaftServer_StartStop(t *testing.T) {
	// Create server
	server, err := NewRaftServer(1, 51200)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test Start
	err = server.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Test if server is actually listening
	time.Sleep(100 * time.Millisecond)
	client, err := rpc.Dial("tcp", fmt.Sprintf(":%d", 51200))
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	client.Close()

	// Test Stop
	server.Stop()

	// Verify server is stopped
	time.Sleep(100 * time.Millisecond)
	_, err = rpc.Dial("tcp", fmt.Sprintf(":%d", 51200))
	if err == nil {
		t.Error("Server still accepting connections after Stop()")
	}
}

// TestGetRandomTimeout tests the random timeout generation
func TestGetRandomTimeout(t *testing.T) {
	// Test multiple random timeouts
	for i := 0; i < 100; i++ {
		timeout := getRandomTimeout()
		if timeout < time.Duration(minTimeout)*time.Millisecond || timeout > time.Duration(maxTimeout)*time.Millisecond {
			t.Errorf("getRandomTimeout() = %v, want between %v and %v", timeout, minTimeout*time.Millisecond, maxTimeout*time.Millisecond)
		}
	}
}

// TestBecomeLeader tests the leader transition process
func TestBecomeLeader(t *testing.T) {
	// Clear the leader election channel before test
	select {
	case <-leaderElected:
	default:
	}

	// Create test servers
	s1, err := NewRaftServer(1, 51300)
	if err != nil {
		t.Fatalf("Failed to create server 1: %v", err)
	}

	s2, err := NewRaftServer(2, 51301)
	if err != nil {
		t.Fatalf("Failed to create server 2: %v", err)
	}

	// Test first server becoming leader
	s1.becomeLeader()
	if s1.state != Leader {
		t.Errorf("Server 1 state = %v, want Leader", s1.state)
	}

	// Test second server becoming follower
	s2.becomeLeader()
	if s2.state != Follower {
		t.Errorf("Server 2 state = %v, want Follower", s2.state)
	}

	// Clean up
	s1.Stop()
	s2.Stop()
}

// TestRequestVote tests the request vote process
func TestRequestVote(t *testing.T) {
	// Create test server
	testPort := 51003
	server, err := NewRaftServer(1, testPort)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer server.Stop()

	// Start server
	err = server.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Create RPC client
	rpc := &RaftRPC{server: server}

	// Create request vote arguments
	args := &RequestVoteArgs{
		Term:         1,
		CandidateID:  2,
		LastLogIndex: 0,
		LastLogTerm:  0,
	}

	// Create reply
	reply := &RequestVoteReply{}

	// Test request vote
	err = rpc.RequestVote(args, reply)
	if err != nil {
		t.Errorf("RequestVote() error = %v", err)
	}
}
