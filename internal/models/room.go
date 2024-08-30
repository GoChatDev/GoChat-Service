package models

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	users map[*user]bool
	join  chan *user
	leave chan *user

	forward chan []byte
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *user),
		leave:   make(chan *user),
		users:   make(map[*user]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case user := <-r.join:
			r.users[user] = true
		case user := <-r.leave:
			delete(r.users, user)
			close(user.recieve)
		case msg := <-r.forward:
			for user := range r.users {
				user.recieve <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatalf("ServeHTTP: %v", err)
		return
	}

	user := &user{
		socket:  socket,
		recieve: make(chan []byte, messageBufferSize),
		room:    r,
	}

	r.join <- user
	defer func() { r.leave <- user }()
	go user.write()
	user.read()
}
