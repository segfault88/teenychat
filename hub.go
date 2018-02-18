package main

import (
	"log"
	"sync"
)

type client interface {
	getID() string
	getDesc() string
	joined(h *hub)
	close()
	sendMessage(m interface{})
}

type hub struct {
	clients  map[string]client
	wg       sync.WaitGroup
	shutdown chan bool

	clientToTopics  map[string][]string
	topicsToClients map[string][]string
}

func (h *hub) start() {
	h.shutdown = make(chan bool)
}

func (h *hub) stop() {
	close(h.shutdown)
	for _, c := range h.clients {
		c.close()
	}
	h.clients = map[string]client{}
	h.wg.Wait()
}

func (h *hub) connect(c client) {
	h.clients[c.getID()] = c
	h.wg.Add(1)
	c.joined(h)

	log.Printf("Client %s %s joined\n", c.getID(), c.getDesc())
}
