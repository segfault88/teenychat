package main

import (
	"sync"
)

type client interface {
	joined(h *hub)
}

type hub struct {
	clients  []client
	wg       sync.WaitGroup
	shutdown chan bool
}

func (h *hub) start() {

}

func (h *hub) stop() {
	close(h.shutdown)
	h.wg.Wait()
}

func (h *hub) connect(c client) {
	h.clients = append(h.clients, c)
	h.wg.Add(1)
	c.joined(h)
}
