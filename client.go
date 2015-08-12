package main

import (
  "fmt"
  "net"
  "bufio"
  "os"
  "strings"

  "github.com/parkdy/golang-chat/shared"
)

// Entry point
// Start TCP chat client
func main() {
  host, port := shared.GetHostPort()

  fmt.Printf("Connecting to chat server on %s:%s\n", host, port)
  
  connection, err := net.Dial("tcp", host+":"+port)
  if err != nil {
    fmt.Println("Error connecting to server:", err)
    return
  }

  fmt.Println("Connected to chat server")

  clientInputReader := bufio.NewReader(os.Stdin)
  userName := getUserName(clientInputReader)

  userConnection := shared.CreateUserConnection(userName, connection)
  _, err = userConnection.Writer.WriteString(shared.GetHeaderFromUserName(userName))
  if err != nil {
    fmt.Println("Error writing header to server:", err)
  }
  userConnection.Writer.Flush()  

  go handleIncomingMessages(&userConnection)

  handleUserInput(&userConnection, clientInputReader)
}

func getUserName(clientInputReader *bufio.Reader) string {
  fmt.Printf("Enter user name: ")
  line, err := clientInputReader.ReadString('\n')
  if err != nil {
    fmt.Println("Error reading user name:", err)
    return "Default User"
  }

  return strings.TrimSpace(line)
}

func handleIncomingMessages(userConnection *shared.UserConnection) {
  for {
    line, err := userConnection.Reader.ReadString('\n')
    if err != nil {
      fmt.Println("Error reading line from server:", err)
      userConnection.Connection.Close()
      break
    }

    fmt.Printf("\r\033[K" + line + "Chat here: ")
  }
}

func handleUserInput(userConnection *shared.UserConnection, clientInputReader *bufio.Reader) {
  for {
    fmt.Printf("Chat here: ")

    line, err := clientInputReader.ReadString('\n')
    if err != nil {
      fmt.Println("Error getting user input:", err)
    }

    _, err = userConnection.Writer.WriteString(line)
    if err != nil {
      fmt.Println("Error writing line to server:", err)
    }
    userConnection.Writer.Flush()  
  }
}