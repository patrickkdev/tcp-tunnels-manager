# TCP Tunnels Manager

This project implements a TCP tunnels manager in Go. It allows for the creation, management, and monitoring of TCP tunnels, facilitating communication between different network endpoints.

## Why MySQL?

The system uses MySQL because a client required tunnel configurations and logs to be stored in a database.

In many cases a project like this could be configured through a CLI or simple configuration files. However, the client preferred database-backed configuration so the tunnels could integrate with their existing systems and operational workflows.

Because of that requirement, MySQL serves as the control and persistence layer for configuration and logs.

## Why `socat`?

Actual TCP forwarding is handled by `socat`.

Rather than reimplementing a TCP proxy in Go, the project delegates the networking layer to `socat`, a mature and widely used Unix utility designed for socket and TCP bridging.

The Go application focuses on orchestration:

- starting and stopping tunnel processes
- supervising process lifecycle
- restarting tunnels if they exit
- collecting and storing logs
- loading configuration from the database

This keeps the implementation simple while relying on a battle-tested networking tool.

## Features

- **Tunnel Management** — Create, update, and remove tunnel configurations.
- **Process Supervision** — Automatically restarts tunnels if they exit unexpectedly.
- **Database Integration** — Stores configuration and logs in MySQL.
- **Tunnel Logging** — Captures process output and records events.
- **Clean Architecture** — Clear separation between domain, application logic, and infrastructure.

## Project Structure

The project follows a clean architecture style layout:

- `cmd/app` — Application entrypoint
- `configs` — Configuration and initialization
- `internal/application` — Core business logic and use cases
- `internal/domain` — Core entities and interfaces
- `internal/infrastructure` — Implementations for external systems (database, TCP tunnels)

## Requirements

- Go
- MySQL
- socat

### Installing socat

Debian / Ubuntu:

```bash
apt install socat
```

Arch Linux:

```bash
pacman -S socat
```

## Setup

### 1. Clone the repository

```bash
git clone https://github.com/patrickkdev/tcp-tunnels-manager.git
cd tcp-tunnels
```

### 2. Configure environment variables

Create a `.env` file in the project root with the required configuration values.

Refer to `configs/init.go` for the expected environment variables.

### 3. Run the application

```bash
go run cmd/app/main.go
```

## Database Schema

The `schema.sql` file defines the tables used by the system:

- `tcp_tunnels` — Stores the configuration for each tunnel
- `tcp_tunnel_logs` — Stores tunnel activity and process logs

## Contributing

Contributions are welcome. Feel free to open issues or submit pull requests.
