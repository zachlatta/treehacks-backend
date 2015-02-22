package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Server started...")

	go h.run()
	go gm.run()
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(":8889", nil))
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &conn{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writePump()
	c.readPump()
}
