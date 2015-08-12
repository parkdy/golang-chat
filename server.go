package main

import (
  "fmt"
  "net"

  "github.com/parkdy/golang-chat/shared"
)

// Entry point
// Start TCP chat server
func main() {
  host, port := shared.GetHostPort()

  fmt.Printf("Starting chat server on %s:%s\n", host, port)

  listener, err := net.Listen("tcp", host+":"+port)
  if err != nil {
    fmt.Println("Error starting server:", err)
  }

  fmt.Println("Started chat server")

  // Maintain list of active user connections
  userConnections := make([]shared.UserConnection, 0)

  for {
    connection, err := listener.Accept()
    if err != nil {
      fmt.Println("Error accepting connection:", err)
    }

    go handleConnection(connection, &userConnections)
  }
}

func handleConnection(connection net.Conn, userConnections *[]shared.UserConnection) {
  fmt.Println("Handling connection from:", connection.RemoteAddr())
  
  // Get user name and create user connection
  userConnection := shared.CreateUserConnection("placeholder", connection)
  header, err := userConnection.Reader.ReadString('\n')
  if err != nil {
    fmt.Println("Error reading header from connection:", err)
  }
  userName := shared.GetUserNameFromHeader(header)
  userConnection.UserName = userName

  // Disconnect if user already exists, otherwise add the new user connection to our list
  for _, existingUserConnection := range *userConnections {
    if userName == existingUserConnection.UserName {
      errorMessage := fmt.Sprintf("Error: the user '%s' already exists\n", userName)
      userConnection.Writer.WriteString(errorMessage)
      userConnection.Connection.Close()
      return
    }
  }
  *userConnections = append(*userConnections, userConnection)
  fmt.Printf("The user '%s' has connected to the server\n", userName)

  handleIncomingMessages(&userConnection, userConnections)
}

func handleIncomingMessages(userConnection *shared.UserConnection, userConnections *[]shared.UserConnection) {
  for {
    line, err := userConnection.Reader.ReadString('\n')
    if err != nil {
      fmt.Println("Error reading line from connection:", err)
      userConnection.Connection.Close()
      break
    }

    pushMessageToClients(line, userConnection.UserName, userConnections)
  }
}

func pushMessageToClients(message string, senderName string, userConnections *[]shared.UserConnection) {
  for _, userConnection := range *userConnections {
    fullMessage := fmt.Sprintf("%s: %s", senderName, message)
    _, err := userConnection.Writer.WriteString(fullMessage)
    if err != nil {
      fmt.Println("Error writing to client:", err)
    }
    userConnection.Writer.Flush()
  }
}