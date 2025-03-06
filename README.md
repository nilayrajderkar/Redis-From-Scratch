# Redis Server Implementation in Go

## Overview
This project is a Redis server implementation written in Go, aiming to replicate core Redis functionality. It implements the RESP (Redis Serialization Protocol) for client-server communication and provides a command-line interface for interacting with the server.

## Current Features

### RESP Protocol Implementation
- Full RESP (Redis Serialization Protocol) serialization and deserialization
- Supports all RESP data types:
  - Simple Strings (prefixed with "+")
  - Errors (prefixed with "-")
  - Integers (prefixed with ":")
  - Bulk Strings (prefixed with "$")
  - Arrays (prefixed with "*")

### Supported Commands
Currently implemented commands:
- `PING` - Returns PONG
- `ECHO <message>` - Returns the message

### Client Interface
- Interactive command-line interface
- Support for quoted arguments
- Real-time command processing
- Error handling and display

## Planned Features
- Key-Value Store Operations:
  - SET
  - GET
  - DEL
  - EXISTS
- Data Structure Commands:
  - Lists (LPUSH, RPUSH, LPOP, RPOP)
  - Sets (SADD, SREM, SMEMBERS)
  - Hashes (HSET, HGET, HDEL)
- Persistence
- TTL Support
- Pub/Sub System

## Project Structure 
```
├── client/         # Client implementation
├── resp/           # RESP serializer and deserializer
├── server/         # Server implementation
├── main.go         # Entry point for the server
```

## Requirements
- Go 1.22.1 or higher

## Installation
1. Clone the repository:

```bash
git clone https://github.com/nilayrajderkar/redis-implementation.git
cd redis-implementation
```

2. Build the project:
```bash
go build
```

## Running the Server
To start the server:
```bash
go run main.go
```

The server will start on localhost:6379 (default Redis port).

## Running Tests
To run all tests:
```bash
go test ./...
```

To run tests with coverage:
```bash
go test ./... -cover
```

## Development

### Running Lint Checks
The project uses golangci-lint for code quality. To run lint checks:
```bash
golangci-lint run
```

### Debugging
A VSCode launch configuration is provided for debugging. You can start the debugger using the "Debug Redis Client" configuration in VSCode.

## Protocol Specification

### RESP Protocol Implementation
The RESP protocol is implemented in the `resp` package with the following components:

1. Serializer (`resp/resp_serializer.go`):
   - Converts Go data types to RESP format
   - Handles strings, integers, arrays, and errors

2. Deserializer (`resp/resp_deserializer.go`):
   - Parses RESP format into Go data types
   - Handles all RESP data types with error checking

## Contributing
1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments
- Redis Protocol Specification: [Redis Protocol](https://redis.io/topics/protocol)
- Go Programming Language: [golang.org](https://golang.org)

