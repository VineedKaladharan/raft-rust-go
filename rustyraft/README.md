# Rusty Raft Server

A Rust implementation of multiple Raft servers with random ports and timeouts.

## Features
- Spawn multiple Raft servers with random port assignments
- Each server has a random timeout between 1-5 seconds
- Asynchronous server execution using Tokio
- Logging of server events

## Prerequisites
- Rust and Cargo installed
- Required dependencies (automatically installed via Cargo):
  - tokio
  - rand
  - log
  - env_logger
  - futures

## Running the Program

To run the program with logging enabled, use one of the following commands:

### Windows PowerShell
```powershell
$env:RUST_LOG='info'; cargo run -- <number_of_servers>
```

### Windows Command Prompt (cmd)
```cmd
set RUST_LOG=info && cargo run -- <number_of_servers>
```

### Linux/MacOS
```bash
RUST_LOG=info cargo run -- <number_of_servers>
```

Replace `<number_of_servers>` with the desired number of servers to run.

Example to run 5 servers:
```powershell
$env:RUST_LOG='info'; cargo run -- 5
```

## To run the tests:
```bash
cargo test
```

## Output
The program will show:
- Server startup messages with randomly assigned ports
- Timeout events for each server
- Program completion message

## Example Output
```
Server 0 starting on port 4532
Server 1 starting on port 6789
Server 2 starting on port 3456
Server 0 timed out after 2345ms
Server 2 timed out after 3456ms
Server 1 timed out after 4567ms
All servers completed
```
