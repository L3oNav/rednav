# Rednav Server

Rednav is a Redis-like in-memory database server written in Go. It offers high-performance data storage and retrieval and can operate in both standalone and replication modes for high availability and data redundancy. Below are the instructions and details for getting started with Rednav Server.

## Features

- In-memory data storage with optional expiration.
- Basic commands like `SET`, `GET`, `ECHO`, `PING`, and `INFO`.
- Replication support with a master-replica configuration.
- Concurrent client connections handling.
- Custom RESP (REdis Serialization Protocol) parsing.
- Unit and integration testing examples.
- Simple TCP server implementation.

## Getting Started

### Prerequisites

Before running Rednav, ensure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.21 or higher)

### Installation

Clone the Rednav repository to your local machine:

```bash
git clone https://github.com/your_github_username/rednav.git
cd rednav
```

### Running The Server

Run the server in the main directory:

```bash
go run ./main.go
```

By default, it will start listening on `localhost` at port `3312`.

### Configuration Options

To configure the Rednav server, you can use command-line flags:

- To specify a host and port for the server to listen on, use:

```bash
go run ./main.go --host your_host --port your_port
```

- To run Rednav as a replica of another instance, use:

```bash
go run ./main.go --replica_of "master_host master_port"
```

### Testing

To run the existing tests, use the following command from the project root:

```bash
go test ./...
```

### Shutting Down

Rednav can be shut down gracefully by sending an interrupt signal (Ctrl+C) in the terminal.

## Architecture

Rednav's architecture consists of several primary components:

- The `main.go` file initializes the server configuration and manages the startup process.
- The `server` package contains logic for accepting and processing client connections.
- The `commands` package contains implementatons for RESP commands.
- The `app` package encapsulates the in-memory data storage and application logic.
- The `utils` package provides utility functions like RESP parsing and logging.

## Contributing

Contributions are welcome! To contribute, please follow these steps:

1. Fork the repository on GitHub.
2. Create a new feature branch from the master branch.
3. Commit your changes with meaningful commit messages.
4. Push your changes to your fork and submit a pull request.

Before submitting a pull request, ensure you've added tests for your changes and that all existing tests pass.

## License

This project is open-sourced under the MIT License. See the [LICENSE](https://platform.openai.com/playground/LICENSE) file for more information.

## Acknowledgments

- Special thanks to the Redis team for inspiration.
- Contributions from the open-source community.

For more information or queries, please open an issue on the GitHub repository.
