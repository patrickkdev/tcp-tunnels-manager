# TCP Tunnels Manager

This project implements a TCP tunnels manager in Go. It allows for the creation, management, and monitoring of TCP tunnels, facilitating communication between different network endpoints.

## Why MySQL?

This project uses a MySQL database because it was a requirement from a client who requested that tunnel configurations and logs be persisted in a database.

The original intention for a project like this would typically be a simpler CLI-driven configuration or file-based approach. However, the client preferred storing configuration and activity logs in database tables for integration with their existing systems and workflows.

As a result, MySQL is used as the control and persistence layer for tunnel configurations and logs.

## Features

- **Tunnel Management:** Create, update, and delete TCP tunnel configurations.
- **Database Integration:** Stores tunnel configurations and logs in a MySQL database.
- **Tunnel Logging:** Records connection events and data transfer for each tunnel.
- **Client-Side Tunneling:** Provides functionality for establishing and managing TCP tunnel clients.

## Project Structure

The project follows a clean architecture pattern, separating concerns into distinct layers:

- `cmd/app`: Contains the main application entry point.
- `configs`: Handles application configuration and initialization.
- `internal/application`: Implements the core business logic and use cases (e.g., `tunnels_manager.go`).
- `internal/domain`: Defines the core entities and interfaces (e.g., `tunnel_log.go`, `tunnel_row.go`).
- `internal/infrastructure`: Provides implementations for external concerns like database access (`db`) and TCP tunneling (`tcptunnels`).

## Getting Started

### Prerequisites

- Go
- MySQL Database

### Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/patrickkdev/tcp-tunnels.git
   cd tcp-tunnels
   ```

2. **Set up environment variables:**

   Create a `.env` file in the project root with your database connection string and other necessary configurations. Refer to `configs/init.go` for required variables.

3. **Run the application:**
   ```bash
   go run cmd/app/main.go
   ```

## Database Schema

The `schema.sql` file defines the database tables for `tunnel_rows` and `tunnel_logs`:

- `tunnel_rows`: Stores the configuration for each TCP tunnel.
- `tunnel_logs`: Records events and statistics related to tunnel activity.

## Contributing

Contributions are welcome. Feel free to submit pull requests or open issues.