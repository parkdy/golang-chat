package shared

import (
  "fmt"
  "net"
  "bufio"
  "os"
  "strings"
)

type UserConnection struct {
  UserName string
  Connection net.Conn
  Reader bufio.Reader
  Writer bufio.Writer
}

func CreateUserConnection(userName string, connection net.Conn) UserConnection {
  return UserConnection {
    UserName: userName,
    Connection: connection,
    Reader: *bufio.NewReader(connection),
    Writer: *bufio.NewWriter(connection),
  }
}

func GetHeaderFromUserName(userName string) string {
  return fmt.Sprintf("UserName: %s\n", userName)
}

func GetUserNameFromHeader(header string) string {
  userName := strings.SplitAfterN(header, ":", 2)[1]
  return strings.TrimSpace(userName)
}

func GetHostPort() (string, string) {
  host, port := "127.0.0.1", "8080"
  
  if len(os.Args) > 2 {
    port = os.Args[2]
  }
  if len(os.Args) > 1 {
    host = os.Args[1]
  }

  return host, port
}