# golang-chat
A simple chat server for learning Go. It uses the Gin HTTP framework and Gorilla toolkit for sessions and WebSockets. Currently there is only one room per server, and users are uniquely identified by their name.

## Usage
1. Build the server binary with `go build`.
2. Change the secret token in the `config/secret_token` file.
3. Launch the server by running `./golang-chat`.
4. Visit `http://localhost:8080` on separate browser sessions to see the chat room in action.
