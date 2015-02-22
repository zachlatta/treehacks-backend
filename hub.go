package main

import "fmt"

type hub struct {
	broadcast  chan []byte
	register   chan *conn
	unregister chan *conn
	conns      map[*conn]bool
}

var h = hub{
	broadcast:  make(chan []byte),
	register:   make(chan *conn),
	unregister: make(chan *conn),
	conns:      make(map[*conn]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			fmt.Println("register hub")
			h.conns[c] = true
			gm.register <- c
		case c := <-h.unregister:
			gm.unregister <- c
			if _, ok := h.conns[c]; ok {
				delete(h.conns, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for c := range h.conns {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.conns, c)
				}
			}
		}
	}
}
