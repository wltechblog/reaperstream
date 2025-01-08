package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//go:embed index.html
var indexhtml []byte
var upgrader = websocket.Upgrader{} // use default options

func httpserver() {
	log.Printf("Starting http server. Open http://%s in your browser!", c.listen)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(indexhtml)
	})
	// add websocket handler
	http.HandleFunc("/data/", http.FileServer(http.Dir(".")).ServeHTTP)
	http.HandleFunc("/ws", wsHandler)

	http.ListenAndServe(c.listen, nil)

}

// websocket handeler that receives json from our car channel and sends it to the client
func wsHandler(w http.ResponseWriter, r *http.Request) {
	if c.debug {
		log.Println("new connection on websocket")
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// Register subscription with cars
	rec := Subscribe()
	for {
		c := <-rec
		err = ws.WriteJSON(c)
		if err != nil { //dead client
			Unsubscribe(rec)
			return
		}
	}

}
