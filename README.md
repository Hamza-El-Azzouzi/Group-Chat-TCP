# Group Chat TCP Application

## Overview

This project is a Go-based implementation of a NetCat-like group chat application using TCP socket programming. It provides a robust client-server architecture that allows multiple clients to connect, communicate, and share messages in real-time.

## Features

- 🌐 TCP connection supporting multiple clients (1-to-many relationship)
- 👥 Client name authentication
- 📝 Real-time message broadcasting
- 🕒 Timestamped messages
- 📜 Chat history preservation
- 🔒 Connection limit (maximum 10 connections)
- 🖥️ Linux terminal-style welcome screen

## Prerequisites

- Go (Golang) installed
- Basic understanding of networking concepts

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/group-chat-tcp.git
cd group-chat-tcp
```

2. Build the application:
```bash
go build
```

## Usage

### Server Mode

Run the server with a default or custom port:

```bash
# Default port (8989)
go run . 

# Custom port
go run . 2525
```

### Client Mode

Connect to the server using netcat:

```bash
nc localhost 8989
```

### Example Interaction

```console
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]: YourName
[2024-06-06 15:30:45][YourName]:Hello, everyone!
```

## Key Specifications

- Supports TCP connections
- Utilizes Go-routines for concurrent handling
- Implements channels or mutexes for synchronization
- Robust error handling
- Follows Go best practices

## Learning Objectives

This project helps developers understand:
- Network socket programming
- Go concurrency patterns
- Goroutines and channels
- Mutex synchronization
- Error handling in distributed systems

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Contact

Hamza-El_Azzouzi - elazzouzi.hamza20@gmail.com
