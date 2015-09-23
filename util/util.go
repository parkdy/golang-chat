package util

import (
	"os"

	"github.com/gorilla/websocket"
)

type UserConnection struct {
	UserName   string
	Connection *websocket.Conn
}

func CreateUserConnection(userName string, conn *websocket.Conn) UserConnection {
	return UserConnection{
		UserName:   userName,
		Connection: conn,
	}
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
