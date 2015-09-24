package main

import (
	"errors"
	"fmt"
	"os"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/parkdy/golang-chat/util"
)

// Global
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var userConnections []util.UserConnection

var secretToken = []byte(os.Getenv("SECRET_TOKEN"))	
var sessionStore = sessions.NewCookieStore(secretToken)

// Entry point
// Start HTTP server
func main() {
	// Get host and port from command line arguments
	host, port := util.GetHostPort()

	fmt.Printf("Starting HTTP server on %s:%s\n", host, port)

	router := gin.Default() // Initialize router

	router.Static("/assets", "./assets") // Serve static assets
	router.LoadHTMLGlob("templates/*")   // Load HTML template directory

	// Handle requests
	router.GET("/", getRoot)
	router.GET("/login", getLogin)
	router.POST("/login", postLogin)
	router.POST("/logout", postLogout)
	router.POST("/messages", postMessage)
	router.GET("/ws", getWebSocket)

	// Start server
	router.Run(host + ":" + port)
}

func getCurrentUserName(writer http.ResponseWriter, request *http.Request) (string, error) {
	session, _ := sessionStore.Get(request, "session")
	userName := session.Values["username"]
	if userName != nil {
		return userName.(string), nil
	} else {
		return "", errors.New("Not logged in")
	}
}

func getRoot(context *gin.Context) {
	currentUserName, err := getCurrentUserName(context.Writer, context.Request)
	if err != nil {
		context.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	context.HTML(http.StatusOK, "index.tmpl", gin.H{
		"userName": currentUserName,
	})
}

func getLogin(context *gin.Context) {
	context.HTML(http.StatusOK, "login.tmpl", "")
}

func postLogin(context *gin.Context) {
	session, _ := sessionStore.Get(context.Request, "session")
	session.Values["username"] = context.PostForm("user[name]")
	session.Save(context.Request, context.Writer)
	context.Redirect(http.StatusMovedPermanently, "/")
}

func postLogout(context *gin.Context) {
	session, _ := sessionStore.Get(context.Request, "session")
	session.Values["username"] = nil
	session.Save(context.Request, context.Writer)
	context.Redirect(http.StatusMovedPermanently, "/login")
}

func postMessage(context *gin.Context) {
	currentUserName, err := getCurrentUserName(context.Writer, context.Request)
	if err != nil {
		context.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	message := context.PostForm("message")
	fullMessage := []byte(currentUserName + ": " + message)

	for _, userConnection := range userConnections {
		userConnection.Connection.WriteMessage(websocket.TextMessage, fullMessage)
	}

	context.JSON(http.StatusOK, gin.H{})
}

func getWebSocket(context *gin.Context) {
	handleWebSocket(context.Writer, context.Request)
}

func handleWebSocket(writer http.ResponseWriter, request *http.Request) {
	conn, err := wsUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("Error handling WebSocket:", err)
		return
	}

	userName, err := getCurrentUserName(writer, request)
	userConnection := util.CreateUserConnection(userName, conn)
	userConnections = append(userConnections, userConnection)
}
