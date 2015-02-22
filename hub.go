package main

type hub struct {
	incoming   chan *connMsg
	broadcast  chan interface{}
	register   chan *conn
	unregister chan *conn
	conns      map[*conn]bool
}

var h = hub{
	incoming:   make(chan *connMsg),
	broadcast:  make(chan interface{}),
	register:   make(chan *conn),
	unregister: make(chan *conn),
	conns:      make(map[*conn]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.conns[c] = true
			gm.register <- c
		case c := <-h.unregister:
			gm.unregister <- c
			if _, ok := h.conns[c]; ok {
				delete(h.conns, c)
				close(c.send)
			}
		case m := <-h.incoming:
			gm.incoming <- m
		case i := <-h.broadcast:
			for c := range h.conns {
				select {
				case c.send <- i:
				default:
					close(c.send)
					delete(h.conns, c)
				}
			}
		}
	}
}
