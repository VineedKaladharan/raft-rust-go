package main

import (
	"testing"
	"time"
)

func TestCreateRaftCluster(t *testing.T) {
	tests := []struct {
		name       string
		numServers int
		wantErr    bool
	}{
		{
			name:       "Create single server",
			numServers: 1,
			wantErr:    false,
		},
		{
			name:       "Create multiple servers",
			numServers: 3,
			wantErr:    false,
		},
		{
			name:       "Create zero servers",
			numServers: 0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			servers, err := createRaftCluster(tt.numServers)
			if (err != nil) != tt.wantErr {
				t.Errorf("createRaftCluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(servers) != tt.numServers {
				t.Errorf("createRaftCluster() created %v servers, want %v", len(servers), tt.numServers)
			}

			// Verify each server is running on the correct port
			for i, server := range servers {
				expectedPort := basePort + i
				if server.Port != expectedPort {
					t.Errorf("Server %d running on port %d, want %d", i+1, server.Port, expectedPort)
				}
			}

			// Clean up
			for _, server := range servers {
				server.Stop()
			}
			// Give time for servers to shut down
			time.Sleep(100 * time.Millisecond)
		})
	}
}

func TestClusterShutdown(t *testing.T) {
	// Create a small cluster
	numServers := 2
	servers, err := createRaftCluster(numServers)
	if err != nil {
		t.Fatalf("Failed to create cluster: %v", err)
	}

	// Stop all servers
	for _, srv := range servers {
		srv.Stop()
	}

	// Verify all servers are stopped by checking ports
	time.Sleep(100 * time.Millisecond) // Give servers time to shut down
	for i := 0; i < numServers; i++ {
		port := basePort + i
		if isPortOpen(port) {
			t.Errorf("Port %d still open after shutdown", port)
		}
	}
}

// Helper function to check if a port is open
func isPortOpen(port int) bool {
	servers, err := createRaftCluster(1)
	if err != nil {
		return false
	}
	servers[0].Stop()
	return err == nil
}
