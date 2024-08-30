package server

import (
	"log"
	"net/http"

	"github.com/GoChatDev/GoChat-Service/internal/auth"
	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader  websocket.Upgrader
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins for simplicity
		},
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/ws", s.handleConnections)
	go s.handleMessages()

	log.Println("Server started on :8080")
	return http.ListenAndServe(":8080", nil)
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Authenticate the user
	userID, err := auth.AuthenticateUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	s.clients[ws] = true
	log.Printf("User %s connected", userID)

	// Handle incoming messages
	for {
		var msg []byte
		if err := ws.ReadMessage(&msg); err != nil {
			log.Printf("Error reading message: %v", err)
			delete(s.clients, ws)
			ws.Close()
			break
		}
		s.broadcast <- msg
	}
}

func (s *Server) handleMessages() {
	for {
		msg := <-s.broadcast
		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				client.Close()
				delete(s.clients, client)
			}
		}
	}
}
