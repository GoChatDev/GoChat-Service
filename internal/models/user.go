package models

import (
	"github.com/gorilla/websocket"
)

type user struct {
	socket *websocket.Conn

	recieve chan []byte

	room *room
}

func (u *user) read() {
	defer u.socket.Close()
	for {
		_, msg, err := u.socket.ReadMessage()
		if err != nil {
			return
		}
		u.room.forward <- msg
	}
}

func (u *user) write() {
	defer u.socket.Close()
	for msg := range u.recieve {
		err := u.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
