package main

import (
	"fmt"
	ingest_http "ias/automation/ingest/http"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	ingest_http.SetupRoutes()
	for {
		fmt.Print("\nCommand (start/stop/restart/quit): ")
		var cmd string
		fmt.Scanln(&cmd)

		switch cmd {
		case "start":
			if !ingest_http.IsRunning {
				ingest_http.StartServer()
			} else {
				fmt.Println("Already running")
			}
		case "stop":
			if ingest_http.IsRunning {
				ingest_http.StopServer()
			} else {
				fmt.Println("Not running")
			}
		case "restart":
			ingest_http.StopServer()
			time.Sleep(1 * time.Second)
			ingest_http.StartServer()
		case "quit":
			ingest_http.StopServer()
			return
		}
	}
}
