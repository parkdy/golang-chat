package main

import (
  "fmt"
  "net"
  "bufio"
  "os"
)

// Entry point
// Start TCP chat client
func main() {
  fmt.Println("Connecting to chat server")
  
  connection, err := net.Dial("tcp", "127.0.0.1:8080")
  if err != nil {
    fmt.Println("Error connecting to server:", err)
    return
  }

  go handleIncomingMessages(connection)

  handleUserInput(connection)
}

func handleIncomingMessages(connection net.Conn) {
  reader := bufio.NewReader(connection)
  for {
    line, err := reader.ReadString('\n')
    if err != nil {
      fmt.Println("Error reading line from server:", err)
      connection.Close()
      break
    }

    fmt.Printf("\r\033[K" + line + "Chat here: ")
  }
}

func handleUserInput(connection net.Conn) {
  writer := bufio.NewWriter(connection)
  reader := bufio.NewReader(os.Stdin)
  for {
    fmt.Printf("Chat here: ")

    line, err := reader.ReadString('\n')
    if err != nil {
      fmt.Println("Error getting user input:", err)
    }

    _, err = writer.WriteString(line)
    if err != nil {
      fmt.Println("Error writing line to server:", err)
    }
    writer.Flush()  
  }
}