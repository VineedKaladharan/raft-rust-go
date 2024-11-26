use std::env;
use log::info;
use rustyraft::create_and_run_servers;

#[tokio::main]
async fn main() {
    // Initialize the logger for output
    env_logger::init();

    // Parse command line arguments
    let args: Vec<String> = env::args().collect();
    if args.len() != 2 {
        eprintln!("Usage: {} <number_of_servers>", args[0]);
        std::process::exit(1);
    }

    // Parse the number of servers from command line
    let num_servers = match args[1].parse::<usize>() {
        Ok(n) => n,
        Err(_) => {
            eprintln!("Error: Please provide a valid number");
            std::process::exit(1);
        }
    };

    info!("Starting Raft cluster with {} servers...", num_servers);
    
    // Run the servers and keep the program alive
    if let Some(leader_id) = create_and_run_servers(num_servers).await {
        info!("Initial leader elected: Server {}", leader_id);
    } else {
        eprintln!("Failed to elect a leader");
    }
}
