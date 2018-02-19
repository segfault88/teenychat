package main

import (
	"fmt"
	"log"

	"github.com/rs/xid"

	"github.com/gorilla/websocket"
)

type wsClient struct {
	id   string
	hub  *hub
	conn *websocket.Conn
}

func (c *wsClient) getID() string {
	return c.id
}

func (c *wsClient) getDesc() string {
	return fmt.Sprintf(
		"wss %s -> %s",
		c.conn.LocalAddr().String(),
		c.conn.RemoteAddr().String())
}

func newWsClient(conn *websocket.Conn) client {
	// use xid to identify this client uniquely
	id := xid.New()
	return &wsClient{
		id:   id.String(),
		conn: conn,
	}
}

func (c *wsClient) joined(h *hub) {
	c.hub = h
	go c.handleInput()

	m := map[string]interface{}{"hello": "world!"}
	c.sendMessage(m)
}

func (c *wsClient) close() {
	c.conn.Close()
	c.conn = nil
	c.hub = nil
}

func (c *wsClient) handleInput() {
	defer c.hub.wg.Done()
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage failed", err)
			return
		}

		log.Printf("Message type: %d message: %s", messageType, message)
		c.hub.send(messageType, string(message))
	}
}

func (c *wsClient) sendMessage(m interface{}) error {
	return c.conn.WriteJSON(m)
}
