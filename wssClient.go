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

	m := map[string]interface{}{"hello": "world!"}
	c.sendMessage(m)
}

func (c *wssClient) close() {
	c.conn.Close()
	c.conn = nil
	c.hub = nil
}

func (c *wssClient) handleInput() {
	defer c.hub.wg.Done()
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

func (c *wssClient) sendMessage(m interface{}) {
	c.conn.WriteJSON(m)
}
