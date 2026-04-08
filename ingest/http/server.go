package http

import (
	"fmt"
	"net/http"
	"os"
	"time"

	redis_lib "github.com/redis/go-redis/v9"
)

var IsRunning = false
var currentServer *http.Server

func SetupRoutes(rdb *redis_lib.Client) {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/GET_ALL_TREE_SENSOR", func(w http.ResponseWriter, r *http.Request) {
		getAllTreeSensorHandler(w, r, rdb)
	})
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
