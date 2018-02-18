package main

import (
	"sync"
)

type client interface {
	joined(h *hub)
	close()
	sendMessage(m interface{})
}

type hub struct {
	clients  []client
	wg       sync.WaitGroup
	shutdown chan bool
}

func (h *hub) start() {
	h.shutdown = make(chan bool)
}

func (h *hub) stop() {
	close(h.shutdown)
	for _, c := range h.clients {
		c.close()
	}
	h.wg.Wait()
}

func (h *hub) connect(c client) {
	h.clients = append(h.clients, c)
	h.wg.Add(1)
	c.joined(h)
}
