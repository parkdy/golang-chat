# golang-chat
A simple TCP chat server and client for learning Go. Currently there is only one room per server, and users are uniquely identified by their name.


Launch the server by running `./server`. You can pass in the desired ip and port using `./server 127.0.0.1 8080`. You can recompile using `go build server.go`.

Launch a CLI chat client by running `./client`. You can specify the server to connect to using `./client 10.9.8.7 6543`. Recompile by running `go build client.go`