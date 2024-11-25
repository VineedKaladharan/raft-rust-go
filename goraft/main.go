// Package main is the entry point for the Raft consensus implementation
package main

import (
	"flag"      // Package for command-line flag parsing
	"fmt"       // Package for formatted I/O
	"goraft/server" // Local package containing Raft server implementation
	"log"       // Package for logging
	"os"        // Package for OS functionality
	"os/signal" // Package for handling OS signals
	"syscall"   // Package for system call primitives
	"time"      // Package for time-related functions
)

const (
	// basePort is the starting port number for Raft servers
	basePort = 52000
)

// createRaftCluster creates and initializes a cluster of Raft servers
// numServers: number of servers to create in the cluster
// returns: slice of created servers and any error encountered
func createRaftCluster(numServers int) ([]*server.RaftServer, error) {
	// Initialize slice to hold server instances
	servers := make([]*server.RaftServer, numServers)
	
	// Create and start each server in the cluster
	for i := 0; i < numServers; i++ {
		// Calculate port for this server
		port := basePort + i
		// Create new server instance
		srv, err := server.NewRaftServer(i+1, port)
		if err != nil {
			return nil, fmt.Errorf("failed to create server %d: %v", i+1, err)
		}
		
		// Start the server
		if err := srv.Start(); err != nil {
			return nil, fmt.Errorf("failed to start server %d: %v", i+1, err)
		}
		
		// Store server in slice
		servers[i] = srv
		log.Printf("Started Raft server %d on port %d", i+1, port)
	}
	
	return servers, nil
}

// main is the entry point of the program
func main() {
	// Define and parse command-line flags
	numServers := flag.Int("n", 3, "Number of Raft servers to create")
	flag.Parse()

	// Validate input
	if *numServers < 1 {
		log.Fatal("Number of servers must be at least 1")
	}

	// Create the Raft cluster
	log.Printf("Creating Raft cluster with %d servers...", *numServers)
	servers, err := createRaftCluster(*numServers)
	if err != nil {
		log.Fatalf("Failed to create Raft cluster: %v", err)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down Raft cluster...")

	// Stop all servers in the cluster
	for _, srv := range servers {
		srv.Stop()
	}
	
	// Wait for servers to shutdown gracefully
	time.Sleep(time.Second)
	log.Println("Raft cluster shutdown complete")
}
