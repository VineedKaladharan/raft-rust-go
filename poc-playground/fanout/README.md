# Leader Election System with Fan-out Pattern

This system implements a distributed leader election mechanism using Go's concurrency features and the Fan-out pattern for message distribution.

## Architecture

```
Architecture Diagram:

+----------------+     +-----------------+     +----------------+
|     State      |     |     Fanout     |     |    Message    |
|----------------|     |-----------------|     |----------------|
| mu: Mutex      |     | mu: RWMutex    |     | leaderID: int |
| leaderChosen   |<--->| subscribers    |<--->| msgType: str  |
| leaderID       |     | input channel  |     +----------------+
+----------------+     | done channel   |
                      +-----------------+
                            ^  ^
                            |  |
                    +-------+  +-------+
                    |                  |
              +-----------+      +-----------+
              |  Leader   |      | Follower  |
              +-----------+      +-----------+

Flow:
[1] Node Election Process
    main() → runNodeElection() → [Becomes Leader or Follower]
                ↓
[2] Leader Path:        [3] Follower Path:
    ↓                       ↓
    runLeaderHeartbeat     runFollowerMessageHandling
    ↓                       ↓
    Publish Messages       handleFollowerMessage
```

## Component Relationships

### 1. Core Components

- **Message**: Data structure for communication
  - Contains leaderID and message type
- **State**: Shared state for leader election
  - Tracks current leader and election status
- **Fanout**: Message distribution system
  - Manages subscriber channels and message broadcasting
- **Subscriber**: Node representation
  - Represents each participant in the election

### 2. Communication Flow
```
Leader → Fanout.Publish() → [input channel] → run() → [subscriber channels] → Followers
```

### 3. Key Interactions

#### a. Node Election
```
runNodeElection()
├── Checks State.leaderChosen
├── If false: Becomes Leader
└── If true: Becomes Follower
```

#### b. Leader Operations
```
runLeaderHeartbeat()
└── Periodically sends heartbeat messages via Fanout
```

#### c. Follower Operations
```
runFollowerMessageHandling()
└── Listens for messages from leader
    └── handleFollowerMessage()
```

### 4. Concurrency Management

- `State.mu`: Protects leader election state
- `Fanout.mu`: Protects subscriber management
- Channels: Handle async communication

### 5. Message Types

- `"leader"`: Leader announcement
- `"heartbeat"`: Leader keepalive

## System Features

The system is designed to be:
- **Thread-safe**: Using mutexes and channels
- **Scalable**: Using non-blocking message delivery
- **Fault-tolerant**: Using heartbeat mechanism
- **Decentralized**: Any node can become leader

## Implementation Details

1. **Node Election Process**:
   - Nodes start simultaneously and compete to become leader
   - First node to acquire the state lock becomes leader

2. **Leader Responsibilities**:
   - Sends regular heartbeats
   - Maintains leadership status

3. **Follower Responsibilities**:
   - Listen for leader messages
   - Acknowledge leader heartbeats
   - Ready to participate in new election if leader fails

4. **Fan-out Pattern Benefits**:
   - Efficient message distribution
   - Non-blocking message delivery
   - Better handling of slow subscribers
   - Centralized message distribution through single goroutine

## Usage

To run the system:

```bash
go run fanout.go
```

This will start a simulation with multiple nodes competing for leadership, followed by the leader sending heartbeats and followers acknowledging them.
