package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type wssClient struct {
	hub  *hub
	conn *websocket.Conn
}

func (c *wssClient) joined(h *hub) {
	c.hub = h
	go c.handleInput()
}

func (c *wssClient) handleInput() {
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage failed", err)
			return
		}

		log.Printf("Message type: %d message: %s", messageType, message)
		response := append([]byte("pong: "), message...)
		if err = c.conn.WriteMessage(1, response); err != nil {
			log.Println("WriteMessage 'pong' failed", err)
			return
		}
	}
}
