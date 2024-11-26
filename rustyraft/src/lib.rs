// Required imports for the Raft implementation
use log::info;
use rand::Rng;
use std::sync::Arc;
use std::time::Duration;
use tokio::sync::broadcast;
use tokio::sync::Mutex;
use tokio::time::sleep;

// Server state enum representing possible states in Raft protocol
#[derive(Debug, Clone, PartialEq)]
pub enum ServerState {
    Follower,
    Leader,
}

// Main server structure representing a node in the Raft cluster
pub struct RaftServer {
    pub id: usize,
    pub port: u16,
    pub timeout_ms: u64,
    pub state: ServerState,
    pub last_heartbeat: u64,
}

// Shared state structure to coordinate leader election
pub struct SharedState {
    pub leader_elected: bool,
    pub leader_id: Option<usize>,
    pub heartbeat_tx: broadcast::Sender<usize>,
}

impl RaftServer {
    pub fn new(id: usize, port: u16, timeout_ms: u64) -> Self {
        RaftServer {
            id,
            port,
            timeout_ms,
            state: ServerState::Follower,
            last_heartbeat: 0,
        }
    }

    // Main server runtime function
    pub async fn run(&mut self, shared_state: Arc<Mutex<SharedState>>) {
        info!(
            "Server {} starting on port {} with timeout {}ms",
            self.id, self.port, self.timeout_ms
        );

        let state_clone = shared_state.clone();

        // Get heartbeat channel
        let mut heartbeat_rx = {
            let state = shared_state.lock().await;
            state.heartbeat_tx.subscribe()
        };

        loop {
            match self.state {
                ServerState::Follower => {
                    // Create RNG inside the loop to avoid Send trait issues
                    let timeout =
                        self.timeout_ms + rand::thread_rng().gen_range(0..self.timeout_ms);
                    let timeout_duration = Duration::from_millis(timeout);

                    tokio::select! {
                    //It's like having two parallel tasks:
                    // ⏰ Task 1: "Wake me up after timeout_duration"
                    // 📨 Task 2: "Tell me if you get a heartbeat message"
                    // Whichever happens first:
                    // If the timer wakes up first (no heartbeat received) → Try to become leader
                    // If a heartbeat arrives before timeout → Stay as follower
                        _ = sleep(timeout_duration) => { 
                            // First branch
                            // No heartbeat received, initiate election
                            // "Hey, the timeout happened! Let's try to become leader"
                            let mut state = state_clone.lock().await;
                            if !state.leader_elected {
                                state.leader_elected = true;
                                state.leader_id = Some(self.id);
                                self.state = ServerState::Leader;
                                info!("Server {} timed out and became Leader", self.id);
                            }
                        }
                        Ok(leader_id) = heartbeat_rx.recv() => { 
                            // Second branch
                            // "Hey, got a heartbeat! Stay as follower"
                            info!("Server {} received heartbeat from Leader {}", self.id, leader_id);
                            // Simply increment heartbeat counter
                            self.last_heartbeat += 1;
                        }
                    }
                }
                ServerState::Leader => {
                    // Send heartbeat to all followers
                    let state = state_clone.lock().await;
                    let _ = state.heartbeat_tx.send(self.id);
                    drop(state);
                    sleep(Duration::from_millis(self.timeout_ms / 2)).await; // Half of timeout for heartbeat interval
                }
            }
        }
    }
}

// Function to create and manage multiple Raft servers
pub async fn create_and_run_servers(num_servers: usize) -> Option<usize> {
    let (heartbeat_tx, _) = broadcast::channel(16);

    let shared_state = Arc::new(Mutex::new(SharedState {
        leader_elected: false,
        leader_id: None,
        heartbeat_tx,
    }));

    let handles: Vec<_> = (0..num_servers)
        .map(|id| {
            let mut server = RaftServer::new(
                id,
                8000 + id as u16,
                3000, // 3 seconds base timeout
            );
            let state_clone = shared_state.clone();
            
            // Understanding Ownership with 'move':
            // ---------------------------------
            // In Java, variable capture is implicit:
            //   executorService.submit(() -> {
            //       server.run(stateClone);  // Variables captured by reference
            //   });
            //
            // In Rust, we must be explicit with 'move':
            //   // ❌ Without move - Compiler Error:
            //   tokio::spawn(async {
            //       server.run(state_clone)  // Error: server might not live long enough
            //   });
            //
            //   // ✅ With move - Ownership transferred:
            //   tokio::spawn(async move {
            //       server.run(state_clone)  // Variables now owned by this task
            //   });
            //
            // What 'move' does (in Java terms):
            //   final Server serverCopy = server;          // Take ownership
            //   final State stateCloneCopy = stateClone;  // Take ownership
            //   executorService.submit(() -> {
            //       serverCopy.run(stateCloneCopy);       // Use owned copies
            //   });
            //   // Original variables can't be used anymore!
            
            tokio::spawn(async move { server.run(state_clone).await })
        })
        .collect();
    // Can't use server or state_clone here anymore!

    // Understanding the Server Handles:
    // ------------------------------
    // When we spawn async tasks with tokio::spawn, it returns a JoinHandle
    // (similar to Java's Future). This is a promise for a future result.
    //
    // Example in Rust:
    //   let handles: Vec<_> = (0..num_servers)
    //       .map(|id| {
    //           tokio::spawn(async move { ... })  // Returns JoinHandle
    //       })
    //       .collect();
    //
    // Java equivalent:
    //   List<Future<?>> futures = servers.stream()
    //       .map(server -> executorService.submit(() -> { ... }))
    //       .collect(Collectors.toList());
    //
    // Why We Need the Loop:
    // -------------------
    // Without the loop:
    //   pub async fn create_and_run_servers(num_servers: usize) -> Option<usize> {
    //       // ... spawn servers ...
    //       return final_state.leader_id;  // ❌ Returns immediately, servers die!
    //   }
    //
    // With the loop:
    //   pub async fn create_and_run_servers(num_servers: usize) -> Option<usize> {
    //       // ... spawn servers ...
    //       for handle in handles {
    //           let _ = handle.await;  // ✅ Keeps servers alive
    //       }
    //       return final_state.leader_id;
    //   }
    //
    // Benefits:
    // --------
    // 1. Keeps Program Alive: Prevents main thread from exiting
    // 2. Error Propagation: Allows handling of server failures
    // 3. Resource Management: Ensures proper cleanup
    
    // Wait for all servers to complete
    for handle in handles {
        let _ = handle.await;
    }

    // Return the leader ID
    let final_state = shared_state.lock().await;
    final_state.leader_id
}
