# ChatGO

ChatGO is a real-time chat application built with Go, featuring secure message encryption, room-based communication, and a modern command-line interface.

## Features

- 🔐 Secure Authentication

  - User registration and login
  - Password encryption
  - JWT-based session management

- 💬 Real-time Messaging

  - WebSocket-based communication
  - Support for multiple chat rooms
  - Message history with customizable limits
  - Message encryption for enhanced security

- 🏠 Room Management

  - Create and join chat rooms
  - Default room support
  - Room member management
  - Room-specific message history

- 🎨 Rich CLI Interface
  - Color-coded messages
  - Command-based interaction
  - Real-time message updates
  - User-friendly formatting

## Architecture

The project follows a clean architecture pattern with separate client and server components:

### Server

- Built with Go and Gin framework
- PostgreSQL database for data persistence
- Redis for caching and real-time features
- WebSocket support for real-time communication
- JWT-based authentication
- Message encryption

### Client

- Command-line interface
- WebSocket client for real-time messaging
- Color-coded output for better readability
- Command-based interaction system

## Project Structure

```
ChatGO/
├── server/
│   ├── cmd/            # Server entry point
│   ├── internal/       # Internal packages
│   │   ├── db/        # Database operations
│   │   ├── models/    # Data models
│   │   ├── services/  # Business logic
│   │   └── transport/ # HTTP/WebSocket handlers
│   ├── router/        # Route definitions
│   ├── build/         # Build configurations
│   └── docs/          # Documentation
└── client/
    ├── main.go        # Client entry point
    ├── main_test.go   # Client tests
    └── color/         # Color formatting utilities
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14
- Redis 6
- Make (for build automation)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/ChatGO.git
cd ChatGO
```

2. Set up the database:

```bash
cd server
make db-setup
```

3. Build the server:

```bash
make build
```

4. Build the client:

```bash
cd ../client
go build
```

## Running the Application

1. Start the server:

```bash
cd server
make run
```

2. Run the client:

```bash
cd ../client
./client -username your_username -password your_password
```

## Client Commands

- `/help` - Display available commands
- `/history [limit]` - View chat history (default: 10 messages)
- `/room [room_id]` - Switch to a different room
- `/create [room_name]` - Create a new room
- `/exit` - Exit the chat

## Testing

Run the test suite:

```bash
# Server tests
cd server
make test

# Client tests
cd ../client
go test -v
```

## Security Features

- End-to-end message encryption
- Secure password storage
- JWT-based authentication
- WebSocket connection validation
- Room access control

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [GORM](https://gorm.io/)
- [JWT-Go](https://github.com/golang-jwt/jwt)
