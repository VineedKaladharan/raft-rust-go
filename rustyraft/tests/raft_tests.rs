use rustyraft::{RaftServer, ServerState, SharedState};
use std::sync::Arc;
use tokio::sync::Mutex;
use std::time::Duration;

#[tokio::test]
async fn test_server_initialization() {
    // Test server creation with specific values
    let server = RaftServer::new(1, 3000, 1000);
    assert_eq!(server.id, 1);
    assert_eq!(server.port, 3000);
    assert_eq!(server.timeout_ms, 1000);
    assert!(matches!(server.state, ServerState::Follower));
}

#[tokio::test]
async fn test_leader_election() {
    // Create two servers with different timeouts
    let mut server1 = RaftServer::new(1, 3000, 100);  // Fast timeout
    let mut server2 = RaftServer::new(2, 3001, 200);  // Slow timeout

    // Create shared state
    let shared_state = Arc::new(Mutex::new(SharedState {
        leader_elected: false,
        leader_id: None,
    }));

    // Run both servers concurrently
    let state_clone1 = Arc::clone(&shared_state);
    let state_clone2 = Arc::clone(&shared_state);

    let (_, _) = tokio::join!(
        server1.run(state_clone1),
        server2.run(state_clone2)
    );

    // Verify that server1 became leader (due to shorter timeout)
    assert!(matches!(server1.state, ServerState::Leader));
    assert!(matches!(server2.state, ServerState::Follower));
}

#[tokio::test]
async fn test_multiple_servers() {
    // Test with 5 servers
    let result = rustyraft::create_and_run_servers(5).await;
    
    // Verify that a leader was elected
    assert!(result.is_some());
    assert!(result.unwrap() < 5); // Leader ID should be valid
}

#[tokio::test]
async fn test_server_state_transitions() {
    let mut server = RaftServer::new(1, 3000, 100);
    
    // Server should start as follower
    assert!(matches!(server.state, ServerState::Follower));
    
    // Create shared state where no leader is elected
    let shared_state = Arc::new(Mutex::new(SharedState {
        leader_elected: false,
        leader_id: None,
    }));
    
    // Run the server
    server.run(Arc::clone(&shared_state)).await;
    
    // Server should have become leader
    assert!(matches!(server.state, ServerState::Leader));
}

#[tokio::test]
async fn test_concurrent_timeout_handling() {
    let shared_state = Arc::new(Mutex::new(SharedState {
        leader_elected: false,
        leader_id: None,
    }));
    
    // Create and run servers concurrently
    let handles: Vec<_> = (0..3)
        .map(|id| {
            let state_clone = Arc::clone(&shared_state);
            tokio::spawn(async move {
                let mut server = RaftServer::new(id, 3000 + id as u16, 100);
                server.run(state_clone).await;
            })
        })
        .collect();
    
    // Wait for all servers to complete
    for handle in handles {
        let _ = handle.await;
    }
    
    // Verify that exactly one leader was elected
    let final_state = shared_state.lock().await;
    assert!(final_state.leader_elected);
    assert!(final_state.leader_id.is_some());
}
