package main

import (
  "fmt"
  "net"
  "bufio"
)

// Entry point
// Start TCP chat server
func main() {
  fmt.Println("Starting chat server")
  
  listener, err := net.Listen("tcp", "127.0.0.1:8080")
  if err != nil {
    fmt.Println("Error starting server:", err)
  }

  // Maintain list of active connections
  activeConnections := make([]net.Conn, 0)

  for {
    connection, err := listener.Accept()
    if err != nil {
      fmt.Println("Error accepting connection:", err)
    }

    go handleConnection(connection, &activeConnections)
  }
}

func handleConnection(connection net.Conn, activeConnections *[]net.Conn) {
  fmt.Println("Handling connection")
  
  *activeConnections = append(*activeConnections, connection)

  handleIncomingMessages(connection, activeConnections)
}

func handleIncomingMessages(connection net.Conn, activeConnections *[]net.Conn) {
  reader := bufio.NewReader(connection)
  for {
    line, err := reader.ReadString('\n')
    if err != nil {
      fmt.Println("Error reading line from connection:", err)
      connection.Close()
      break
    }

    pushMessageToClients(string(line), activeConnections)
  }
}

func pushMessageToClients(message string, activeConnections *[]net.Conn) {
  for _, connection := range *activeConnections {
    writer := bufio.NewWriter(connection)
    _, err := writer.WriteString(message)
    if err != nil {
      fmt.Println("Error writing to client:", err)
    }
    writer.Flush()
  }
}