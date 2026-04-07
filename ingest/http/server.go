package http

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var IsRunning = false
var currentServer *http.Server

func SetupRoutes() {
	http.HandleFunc("/", homeHandler)
}

func StartServer() {
	currentServer = &http.Server{Addr: ":" + os.Getenv("HTTP_SERVER_PORT")}
	go func() {
		IsRunning = true
		fmt.Println("Server started on " + currentServer.Addr)
		currentServer.ListenAndServe()
	}()
	time.Sleep(100 * time.Millisecond) // Give it a moment to start
}

func StopServer() {
	if currentServer != nil {
		currentServer.Close()
		IsRunning = false
		fmt.Println("Server stopped")
	}
}
