# rustyraft & goraft
Raft implementation in Rust and Golang as a learning project

## 👨‍💻 Author
Vineed Kaladharan

## 🚀 Implementation Status

### Phase 1: Simple Multi-Node Implementation 🔄
> **Status:** `🟢 In Progress`

🔄 Basic state machine `🟢 In Progress`
  - Implementation of core state machine that processes and applies commands
  - Handles read/write operations and maintains consistent state
  - Ensures atomic execution of commands

⏳ Term management
  - Implementation of logical clock using term numbers
  - Handling term updates and comparisons
  - Term-based leader election protocol

⏳ Log storage
  - Persistent storage of log entries
  - Log entry format with term numbers and commands
  - Log consistency checking mechanisms

⏳ State transitions
  - Implementation of state persistence
  - Safe state transfer between nodes
  - Crash recovery handling

⏳ Role transitions (Follower → Candidate → Leader)
  - Election timeout management
  - Vote request and response handling
  - Leader heartbeat implementation
  - Role-specific behavior implementation

### Phase 2: Multi-Node Communication 🌐
> **Status:** `⚪ Planned`

- [ ] Node communication setup
- [ ] Leader election implementation
- [ ] Basic log replication
- [ ] Simple membership management
- [ ] Message passing infrastructure

### Phase 3: Full Cluster Features ⚡
> **Status:** `⚪ Planned`

- [ ] Complete log replication
- [ ] Snapshot handling
- [ ] Membership changes
- [ ] Recovery mechanisms
- [ ] Fault tolerance

## 📌 Tasks & Notes

### Tasks
- Created basic func to accept number of servers to run in random ports as async `🟢 Done`
- Added timeout on each servers at random milliseconds. When any one of the servers' gets timed out, it becomes 'Leader'. Rest of the nodes becomes 'Follower' using global concurrent variable `🟢 Done`
- Added func to generate random election timeout between 150-300ms `🟢 Done`
- Added func to generate random heartbeat timeout between 100-200ms `🟢 Done`

## 🛠️ Build and Run

### Go Implementation

```bash
# Build the Go implementation
cd raft-go
go build
 ./raft-go

# Or run directly without building
# Run with default 3 servers
go run main.go

# Run with custom number of servers (e.g., 5)
go run main.go -n 5

Example output:
```bash
2024/11/26 14:17:17 Creating Raft cluster with 5 servers...
2024/11/26 14:17:17 Server 1 started on port 52000 as Follower
2024/11/26 14:17:17 Server 4 started on port 52003 as Follower
2024/11/26 14:17:17 Started Raft server 1 on port 52000
2024/11/26 14:17:17 Server 3 started on port 52002 as Follower
2024/11/26 14:17:17 Server 5 started on port 52004 as Follower
2024/11/26 14:17:17 Started Raft server 5 on port 52004
2024/11/26 14:17:17 Server 2 started on port 52001 as Follower
2024/11/26 14:17:17 Started Raft server 4 on port 52003
2024/11/26 14:17:17 Started Raft server 3 on port 52002
2024/11/26 14:17:17 Started Raft server 2 on port 52001
2024/11/26 14:17:17 Server 5: Yay, I am first! Becoming Leader
2024/11/26 14:17:17 Server 3: Another leader was elected, becoming Follower
2024/11/26 14:17:18 Server 2: Another leader was elected, becoming Follower
2024/11/26 14:17:18 Server 4: Another leader was elected, becoming Follower
2024/11/26 14:17:18 Server 1: Another leader was elected, becoming Follower
2024/11/26 14:17:32 Sending: are you alive? id: 5
2024/11/26 14:17:32 Received: are you alive? id: 5
2024/11/26 14:17:32 Sending: yes i am alive, id: 1
2024/11/26 14:17:32 Received: yes i am alive from server at port 52000
2024/11/26 14:17:32 Received: are you alive? id: 5
2024/11/26 14:17:32 Sending: yes i am alive, id: 2
2024/11/26 14:17:32 Received: yes i am alive from server at port 52001
2024/11/26 14:17:32 Received: are you alive? id: 5
2024/11/26 14:17:32 Sending: yes i am alive, id: 3
2024/11/26 14:17:32 Received: yes i am alive from server at port 52002
2024/11/26 14:17:32 Received: are you alive? id: 5
2024/11/26 14:17:32 Sending: yes i am alive, id: 4
2024/11/26 14:17:32 Received: yes i am alive from server at port 52003
2024/11/26 14:17:33 Server 2: Another leader was elected, becoming Follower
2024/11/26 14:17:33 Server 3: Another leader was elected, becoming Follower
2024/11/26 14:17:33 Server 4: Another leader was elected, becoming Follower
2024/11/26 14:17:33 Server 1: Another leader was elected, becoming Follower
2024/11/26 14:17:40 Shutting down Raft cluster...
2024/11/26 14:17:40 rpc.Serve: accept:accept tcp 127.0.0.1:52000: use of closed network connection
2024/11/26 14:17:40 rpc.Serve: accept:accept tcp 127.0.0.1:52001: use of closed network connection
2024/11/26 14:17:40 rpc.Serve: accept:accept tcp 127.0.0.1:52002: use of closed network connection
2024/11/26 14:17:40 rpc.Serve: accept:accept tcp 127.0.0.1:52003: use of closed network connection
2024/11/26 14:17:40 rpc.Serve: accept:accept tcp 127.0.0.1:52004: use of closed network connection
2024/11/26 14:17:41 Raft cluster shutdown complete.

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests verbosely
go test -v ./...
```

### Rust Implementation

```bash
# Build and run the Rust implementation
cd rustyraft
cargo build
cargo run


## Example Output:
```bash
[2024-11-26T10:24:12Z INFO  rustyraft] Starting Raft cluster with 5 servers...
[2024-11-26T10:24:12Z INFO  rustyraft] Server 2 starting on port 8002 with timeout 3000ms    
[2024-11-26T10:24:12Z INFO  rustyraft] Server 3 starting on port 8003 with timeout 3000ms    
[2024-11-26T10:24:12Z INFO  rustyraft] Server 0 starting on port 8000 with timeout 3000ms    
[2024-11-26T10:24:12Z INFO  rustyraft] Server 1 starting on port 8001 with timeout 3000ms    
[2024-11-26T10:24:12Z INFO  rustyraft] Server 4 starting on port 8004 with timeout 3000ms    
[2024-11-26T10:24:16Z INFO  rustyraft] Server 1 timed out and became Leader
[2024-11-26T10:24:16Z INFO  rustyraft] Server 4 received heartbeat from Leader 1
[2024-11-26T10:24:16Z INFO  rustyraft] Server 0 received heartbeat from Leader 1
[2024-11-26T10:24:16Z INFO  rustyraft] Server 2 received heartbeat from Leader 1
[2024-11-26T10:24:16Z INFO  rustyraft] Server 3 received heartbeat from Leader 1
[2024-11-26T10:24:17Z INFO  rustyraft] Server 3 received heartbeat from Leader 1
[2024-11-26T10:24:17Z INFO  rustyraft] Server 2 received heartbeat from Leader 1
[2024-11-26T10:24:17Z INFO  rustyraft] Server 0 received heartbeat from Leader 1
[2024-11-26T10:24:17Z INFO  rustyraft] Server 0 received heartbeat from Leader 1
[2024-11-26T10:24:19Z INFO  rustyraft] Server 3 received heartbeat from Leader 1
[2024-11-26T10:24:19Z INFO  rustyraft] Server 0 received heartbeat from Leader 1
[2024-11-26T10:24:19Z INFO  rustyraft] Server 2 received heartbeat from Leader 1
[2024-11-26T10:24:19Z INFO  rustyraft] Server 4 received heartbeat from Leader 1
[2024-11-26T10:24:20Z INFO  rustyraft] Server 4 received heartbeat from Leader 1
[2024-11-26T10:24:20Z INFO  rustyraft] Server 2 received heartbeat from Leader 1
[2024-11-26T10:24:20Z INFO  rustyraft] Server 3 received heartbeat from Leader 1
[2024-11-26T10:24:20Z INFO  rustyraft] Server 0 received heartbeat from Leader 1
[2024-11-26T10:24:22Z INFO  rustyraft] Server 3 received heartbeat from Leader 1
[2024-11-26T10:24:22Z INFO  rustyraft] Server 2 received heartbeat from Leader 1
[2024-11-26T10:24:22Z INFO  rustyraft] Server 0 received heartbeat from Leader 1
[2024-11-26T10:24:22Z INFO  rustyraft] Server 4 received heartbeat from Leader 1
[2024-11-26T10:24:23Z INFO  rustyraft] Server 3 received heartbeat from Leader 1
[2024-11-26T10:24:23Z INFO  rustyraft] Server 2 received heartbeat from Leader 1
[2024-11-26T10:24:23Z INFO  rustyraft] Server 0 received heartbeat from Leader 1
[2024-11-26T10:24:23Z INFO  rustyraft] Server 4 received heartbeat from Leader 1
[2024-11-26T10:24:25Z INFO  rustyraft] Server 2 received heartbeat from Leader 1
[2024-11-26T10:24:25Z INFO  rustyraft] Server 3 received heartbeat from Leader 1
[2024-11-26T10:24:25Z INFO  rustyraft] Server 0 received heartbeat from Leader 1
[2024-11-26T10:24:25Z INFO  rustyraft] Server 4 received heartbeat from Leader 1
[2024-11-26T10:24:26Z INFO  rustyraft] Server 3 received heartbeat from Leader 1
[2024-11-26T10:24:26Z INFO  rustyraft] Server 2 received heartbeat from Leader 1
[2024-11-26T10:24:26Z INFO  rustyraft] Server 0 received heartbeat from Leader 1
[2024-11-26T10:24:26Z INFO  rustyraft] Server 4 received heartbeat from Leader 1
```

## 📚 References

- [etcd-io/raft](https://github.com/etcd-io/raft)
- [Raft Paper](https://raft.github.io/raft.pdf)
- [Raft Consensus Simulator](https://observablehq.com/@stwind/raft-consensus-simulator)
- [Raft Visualization](https://observablehq.com/d/c695a57a32e50025)
