// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"net"
// 	"os"
// 	"sync"
// 	"time"
// )

// const (
// 	defaultPort    = "8989"
// 	maxConnections = 10
// 	usageMessage   = "[USAGE]: ./TCPChat $port"
// )

// type Client struct {
// 	conn     net.Conn
// 	name     string
// 	messages chan string // Channel for sending messages to this client
// }

// var (
// 	clients     = make(map[net.Conn]*Client)
// 	history     []string
// 	clientMutex sync.Mutex
// )

// func main() {
// 	// Set the port
// 	port := defaultPort
// 	if len(os.Args) == 2 {
// 		port = os.Args[1]
// 	} else if len(os.Args) > 2 {
// 		fmt.Println(usageMessage)
// 		return
// 	}

// 	// Start listening for incoming connections
// 	listener, err := net.Listen("tcp", ":"+port)
// 	if err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}
// 	defer listener.Close()

// 	fmt.Printf("Listening on port :%s\n", port)

// 	// Handle incoming connections
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Printf("Failed to accept connection: %v", err)
// 			continue
// 		}

// 		clientMutex.Lock()
// 		if len(clients) >= maxConnections {
// 			clientMutex.Unlock()
// 			conn.Write([]byte("Max connections reached, try again later.\n"))
// 			conn.Close()
// 			continue
// 		}
// 		clientMutex.Unlock()

// 		go handleConnection(conn)
// 	}
// }

// // Handle an individual client connection
// func handleConnection(conn net.Conn) {
// 	defer conn.Close()

// 	// Prompt for a client name
// 	conn.Write([]byte("[ENTER YOUR NAME]: "))
// 	scanner := bufio.NewScanner(conn)
// 	if !scanner.Scan() {
// 		return
// 	}

// 	name := scanner.Text()
// 	if name == "" {
// 		conn.Write([]byte("Invalid name, connection closing...\n"))
// 		return
// 	}

// 	client := &Client{
// 		conn:     conn,
// 		name:     name,
// 		messages: make(chan string),
// 	}

// 	clientMutex.Lock()
// 	clients[conn] = client
// 	clientMutex.Unlock()

// 	// Send chat history to the new client
// 	sendHistory(client)

// 	// Announce the new client to others
// 	broadcastMessage(fmt.Sprintf("%s has joined our chat...\n", client.name), client)

// 	// Handle incoming messages in a separate goroutine
// 	go clientWriter(client)

// 	// Read incoming messages from the client
// 	for scanner.Scan() {
// 		message := scanner.Text()
// 		if message == "" {
// 			continue
// 		}
// 		formattedMessage := formatMessage(client.name, message)
// 		addMessageToHistory(formattedMessage)
// 		broadcastMessage(formattedMessage, nil)
// 	}

// 	clientMutex.Lock()
// 	delete(clients, conn)
// 	clientMutex.Unlock()

// 	broadcastMessage(fmt.Sprintf("%s has left our chat...\n", client.name), nil)
// }

// // Send the chat history to a new client
// func sendHistory(client *Client) {
// 	clientMutex.Lock()
// 	defer clientMutex.Unlock()
// 	for _, msg := range history {
// 		client.conn.Write([]byte(msg))
// 	}
// }

// // Broadcast a message to all clients, except optionally one to exclude
// func broadcastMessage(message string, excludeClient *Client) {
// 	clientMutex.Lock()
// 	defer clientMutex.Unlock()

// 	for _, client := range clients {
// 		if client != excludeClient {
// 			client.messages <- message
// 		}
// 	}
// }

// // Client writer goroutine to send messages to the client's connection
// func clientWriter(client *Client) {
// 	for msg := range client.messages {
// 		_, err := client.conn.Write([]byte(msg))
// 		if err != nil {
// 			log.Printf("Error writing to client %s: %v", client.name, err)
// 			return
// 		}
// 	}
// }

// // Add a message to the chat history
// func addMessageToHistory(message string) {
// 	clientMutex.Lock()
// 	defer clientMutex.Unlock()
// 	history = append(history, message)
// }

// // Format a message with a timestamp and username
// func formatMessage(name, message string) string {
// 	timestamp := time.Now().Format("2006-01-02 15:04:05")
// 	return fmt.Sprintf("[%s][%s]: %s\n", timestamp, name, message)
// }
//--------------------------------------------------------------------------------
// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"net"
// 	"os"
// 	"strings"
// 	"sync"
// 	"time"
// )

// const (
// 	defaultPort     = "8989"
// 	maxConnections  = 10
// 	usageMessage    = "[USAGE]: ./TCPChat $port"
// 	logFileName     = "server.log"
// )

// type Client struct {
// 	conn     net.Conn
// 	name     string
// 	messages chan string // Channel for sending messages to this client
// }

// var (
// 	clients     = make(map[net.Conn]*Client)
// 	history     []string
// 	clientMutex sync.Mutex
// 	logger      *log.Logger
// 	logFile     *os.File
// )

// func main() {
// 	// Open log file
// 	var err error
// 	logFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// 	if err != nil {
// 		fmt.Printf("Failed to open log file: %v\n", err)
// 		return
// 	}
// 	defer logFile.Close()
// 	logger = log.New(logFile, "", log.LstdFlags)

// 	// Set the port
// 	port := defaultPort
// 	if len(os.Args) == 2 {
// 		port = os.Args[1]
// 	} else if len(os.Args) > 2 {
// 		fmt.Println(usageMessage)
// 		return
// 	}

// 	// Start listening for incoming connections
// 	listener, err := net.Listen("tcp", ":"+port)
// 	if err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}
// 	defer listener.Close()

// 	fmt.Printf("Listening on port :%s\n", port)
// 	logger.Printf("Server started on port %s", port)

// 	// Handle incoming connections
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			logger.Printf("Failed to accept connection: %v", err)
// 			continue
// 		}

// 		clientMutex.Lock()
// 		if len(clients) >= maxConnections {
// 			clientMutex.Unlock()
// 			conn.Write([]byte("Max connections reached, try again later.\n"))
// 			conn.Close()
// 			continue
// 		}
// 		clientMutex.Unlock()

// 		go handleConnection(conn)
// 	}
// }

// // Handle an individual client connection
// func handleConnection(conn net.Conn) {
// 	defer conn.Close()

// 	// Prompt for a client name
// 	conn.Write([]byte("[ENTER YOUR NAME]: "))
// 	scanner := bufio.NewScanner(conn)
// 	if !scanner.Scan() {
// 		return
// 	}

// 	name := scanner.Text()
// 	if name == "" {
// 		conn.Write([]byte("Invalid name, connection closing...\n"))
// 		return
// 	}

// 	client := &Client{
// 		conn:     conn,
// 		name:     name,
// 		messages: make(chan string),
// 	}

// 	clientMutex.Lock()
// 	clients[conn] = client
// 	clientMutex.Unlock()

// 	// Send chat history to the new client
// 	sendHistory(client)

// 	// Announce the new client to others
// 	broadcastMessage(fmt.Sprintf("%s has joined our chat...\n", client.name), client)
// 	logger.Printf("Client %s connected", client.name)

// 	// Handle incoming messages in a separate goroutine
// 	go clientWriter(client)

// 	// Read incoming messages from the client
// 	for scanner.Scan() {
// 		message := scanner.Text()
// 		if message == "" {
// 			continue
// 		}

// 		if strings.HasPrefix(message, "/name ") {
// 			newName := strings.TrimSpace(strings.TrimPrefix(message, "/name "))
// 			if newName != "" {
// 				changeClientName(client, newName)
// 			}
// 		} else {
// 			formattedMessage := formatMessage(client.name, message)
// 			addMessageToHistory(formattedMessage)
// 			broadcastMessage(formattedMessage, nil)
// 			logger.Printf("Message from %s: %s", client.name, message)
// 		}
// 	}

// 	clientMutex.Lock()
// 	delete(clients, conn)
// 	clientMutex.Unlock()

// 	broadcastMessage(fmt.Sprintf("%s has left our chat...\n", client.name), nil)
// 	logger.Printf("Client %s disconnected", client.name)
// }

// // Change a client's name and inform others
// func changeClientName(client *Client, newName string) {
// 	clientMutex.Lock()
// 	oldName := client.name
// 	client.name = newName
// 	clientMutex.Unlock()

// 	broadcastMessage(fmt.Sprintf("%s is now known as %s\n", oldName, newName), nil)
// 	logger.Printf("Client %s changed name to %s", oldName, newName)
// }

// // Send the chat history to a new client
// func sendHistory(client *Client) {
// 	clientMutex.Lock()
// 	defer clientMutex.Unlock()
// 	for _, msg := range history {
// 		client.conn.Write([]byte(msg))
// 	}
// }

// // Broadcast a message to all clients, except optionally one to exclude
// func broadcastMessage(message string, excludeClient *Client) {
// 	clientMutex.Lock()
// 	defer clientMutex.Unlock()

// 	for _, client := range clients {
// 		if client != excludeClient {
// 			client.messages <- message
// 		}
// 	}
// }

// // Client writer goroutine to send messages to the client's connection
// func clientWriter(client *Client) {
// 	for msg := range client.messages {
// 		_, err := client.conn.Write([]byte(msg))
// 		if err != nil {
// 			log.Printf("Error writing to client %s: %v", client.name, err)
// 			return
// 		}
// 	}
// }

// // Add a message to the chat history
// func addMessageToHistory(message string) {
// 	clientMutex.Lock()
// 	defer clientMutex.Unlock()
// 	history = append(history, message)
// }

// // Format a message with a timestamp and username
// func formatMessage(name, message string) string {
// 	timestamp := time.Now().Format("2006-01-02 15:04:05")
// 	return fmt.Sprintf("[%s][%s]: %s\n", timestamp, name, message)
// }
//--------------------------------------------------------------------------------
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	defaultPort    = "8989"
	maxConnections = 10
	usageMessage   = "[USAGE]: ./TCPChat $port [-v]"
	logFileName    = "server.log" // File for logging
)

type Client struct {
	conn     net.Conn
	name     string
	messages chan string
}

var (
	clients     = make(map[net.Conn]*Client)
	history     []string
	clientMutex sync.Mutex
	logger      *log.Logger
	verbose     bool
)

func main() {
	// Check for verbose flag
	if len(os.Args) > 2 && os.Args[2] == "-v" {
		verbose = true
	}

	// Set up logging
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer logFile.Close()

	// Initialize logger
	logger = log.New(logFile, "", log.LstdFlags)
	if verbose {
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(logFile)
	}

	// Set the port
	port := defaultPort
	if len(os.Args) >= 2 {
		port = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Println(usageMessage)
		return
	}

	// Start TCP server
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	logger.Printf("Listening on port :%s\n", port)

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Printf("Failed to accept connection: %v", err)
			continue
		}

		clientMutex.Lock()
		if len(clients) >= maxConnections {
			clientMutex.Unlock()
			conn.Write([]byte("Max connections reached, try again later.\n"))
			conn.Close()
			continue
		}
		clientMutex.Unlock()

		go handleConnection(conn)
	}
}

// Handle an individual client connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Prompt for a client name
	_, err := conn.Write([]byte("[ENTER YOUR NAME]: "))
	if err != nil {
		logger.Printf("Error writing name prompt to client %s: %v", conn.RemoteAddr().String(), err)
		return
	}

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		logActivity(fmt.Sprintf("Error reading name from client %s: %v", conn.RemoteAddr().String(), scanner.Err()))
		return
	}

	name := scanner.Text()
	if name == "" {
		conn.Write([]byte("Invalid name, connection closing...\n"))
		logActivity(fmt.Sprintf("Invalid name received from client %s", conn.RemoteAddr().String()))
		return
	}

	client := &Client{
		conn:     conn,
		name:     name,
		messages: make(chan string),
	}

	clientMutex.Lock()
	clients[conn] = client
	clientMutex.Unlock()

	// Send chat history to the new client
	sendHistory(client)

	// Announce the new client to others
	broadcastMessage(fmt.Sprintf("%s has joined our chat...\n", client.name), client)
	logActivity(fmt.Sprintf("%s joined the chat", client.name))

	// Handle incoming messages in a separate goroutine
	go clientWriter(client)

	// Read incoming messages from the client
	for scanner.Scan() {
		message := scanner.Text()
		if strings.HasPrefix(message, "/name ") {
			// Handle name change
			newName := strings.TrimSpace(strings.TrimPrefix(message, "/name "))
			if newName != "" {
				oldName := client.name
				client.name = newName
				broadcastMessage(fmt.Sprintf("%s changed their name to %s\n", oldName, newName), nil)
				logActivity(fmt.Sprintf("%s changed their name to %s", oldName, newName))
			}
		} else if message != "" {
			formattedMessage := formatMessage(client.name, message)
			addMessageToHistory(formattedMessage)
			broadcastMessage(formattedMessage, nil)
		}
	}

	clientMutex.Lock()
	delete(clients, conn)
	clientMutex.Unlock()

	broadcastMessage(fmt.Sprintf("%s has left our chat...\n", client.name), nil)
	logActivity(fmt.Sprintf("%s left the chat", client.name))
}

// Send the chat history to a new client
func sendHistory(client *Client) {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	for _, msg := range history {
		client.conn.Write([]byte(msg))
	}
}

// Broadcast a message to all clients, except optionally one to exclude
func broadcastMessage(message string, excludeClient *Client) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	for _, client := range clients {
		if client != excludeClient {
			client.messages <- message
		}
	}
}

// Client writer goroutine to send messages to the client's connection
func clientWriter(client *Client) {
	for msg := range client.messages {
		_, err := client.conn.Write([]byte(msg))
		if err != nil {
			logger.Printf("Error writing to client %s: %v", client.name, err)
			return
		}
	}
}

// Add a message to the chat history
func addMessageToHistory(message string) {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	history = append(history, message)
}

// Format a message with a timestamp and username
func formatMessage(name, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]: %s\n", timestamp, name, message)
}

// Log client activities and messages
func logActivity(activity string) {
	logger.Println(activity)
}
