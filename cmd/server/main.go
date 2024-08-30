package main

import (
	"log"

	"github.com/GoChatDev/GoChat-Service/internal/server"
)

func main() {
	// Initialize and start the WebSocket server
	s := server.NewServer()
	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
