package main

import (
	"fmt"
	"log"

	"github.com/rs/xid"

	"github.com/gorilla/websocket"
)

type wssClient struct {
	id   string
	hub  *hub
	conn *websocket.Conn
}

func (c *wssClient) getID() string {
	return c.id
}

func (c *wssClient) getDesc() string {
	return fmt.Sprintf(
		"wss %s -> %s",
		c.conn.LocalAddr().String(),
		c.conn.RemoteAddr().String())
}

func newWssClient(conn *websocket.Conn) client {
	// use xid to identify this client uniquely
	id := xid.New()
	return &wssClient{
		id:   id.String(),
		conn: conn,
	}
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
