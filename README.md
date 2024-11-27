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

#### Architecture and Flow

##### 1. Cluster Creation Flow
```
main() → createRaftCluster()
├── For each server (concurrent):
│   ├── NewRaftServer(id, port)
│   │   ├── Initialize RaftServer struct
│   │   └── Set initial state as Follower
│   └── server.Start()
│       ├── Start RPC server
│       └── resetElectionTimer()
```

##### 2. Election Process Flow
```
Election Timer Expires
├── State changes to Candidate
├── Increment currentTerm
├── Vote for self
├── Send RequestVote RPCs to all peers
└── If majority votes received:
    └── becomeLeader()
        ├── Change state to Leader
        ├── Cancel election timer
        └── startHeartbeat()
```

##### 3. Heartbeat Process Flow
```
Leader's startHeartbeat()
├── Create heartbeat timer (interval: 15ms)
└── On timer tick:
    └── sendHeartbeats()
        ├── For each peer:
        │   ├── Prepare HeartbeatArgs
        │   └── Send Heartbeat RPC
        └── Handle responses:
            ├── Update term if needed
            └── Step down if higher term seen
```

#### Key Components

##### Server States
- **Follower**: Initial state of all servers
- **Candidate**: State during election process
- **Leader**: State after winning an election

##### Timeouts
- **Election Timeout**: Random duration between 150-300ms
- **Heartbeat Interval**: Fixed at 15ms

##### Key Mechanisms
1. **Leader Election**
   - Triggered by election timeout
   - Requires majority votes to become leader
   - Uses term numbers to maintain consistency

2. **Heartbeat System**
   - Regular heartbeats from leader to maintain authority
   - Prevents unnecessary elections
   - Resets follower election timers

3. **Term Management**
   - Monotonically increasing term numbers
   - Used to detect stale leaders
   - Helps maintain cluster consistency

#### Implementation Details

The implementation uses Go's standard libraries:
- `net/rpc` for RPC communication
- `sync` for mutex and synchronization
- `time` for timer management

##### Key Structures
- `RaftServer`: Main server structure
- `RaftRPC`: RPC interface implementation
- Various RPC argument and reply structures for voting and heartbeat mechanisms

#### Running the Go Implementation

To start a Raft cluster:
```bash
cd goraft
go run main.go -n <number_of_servers>
```
Default number of servers is 3 if not specified.

### Go Implementation Example Output
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
```

### Rust Implementation

#### Architecture and Flow

##### 1. Node Creation and Initialization
```
RawNode::new(config, storage, logger)
├── Create Raft instance
│   ├── Initialize ProgressTracker
│   ├── Setup RaftLog with storage
│   └── Set initial state as Follower
├── Load initial state
│   ├── Load HardState if exists
│   └── Apply committed entries
└── Start election timer
```

##### 2. Election Process
```
Election Timer Expires
├── Increment term
├── Change state to Candidate
├── Vote for self
├── Send RequestVote RPCs
│   ├── Include last log index/term
│   └── Wait for responses
└── If majority votes received:
    ├── Change state to Leader
    ├── Initialize volatile leader state
    └── Start sending heartbeats
```

##### 3. Heartbeat Mechanism
```
Leader's tick_heartbeat()
├── Increment heartbeat counter
├── If heartbeat_elapsed >= timeout:
│   ├── Reset heartbeat timer
│   ├── Send MsgBeat to self
│   └── Broadcast heartbeats to followers
└── Handle responses:
    ├── Update progress
    ├── Check term
    └── Step down if higher term seen
```

#### Key Components

##### Core Structures
- `RawNode`: Main interface for external interaction
- `Raft`: Core consensus logic implementation
- `Storage`: Interface for persistent state
- `ProgressTracker`: Tracks peer progress and voting

##### Timeouts and Intervals
- **Heartbeat Tick**: Every 3 ticks
- **Election Tick**: 10 ticks
- **Base Tick**: 100ms (configurable)

##### State Management
1. **Volatile State**
   - Current term
   - Vote history
   - Log entries

2. **Persistent State**
   - Hard state (term, vote, commit)
   - Log entries
   - Snapshot metadata

#### Running the Rust Implementation

To start a Raft cluster:
```bash
cd rustyraft
cargo run
```

The implementation supports:
- Custom storage backends
- Configurable timeouts
- Snapshot and log compaction
- Leader election and transfer
- Dynamic membership changes

### Rust Implementation Example Output
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
