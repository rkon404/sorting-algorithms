package main

import (
	websocket "backend/pkg"
	"log"
	"net/http"
)

func main() {
	wsHandler := websocket.NewHandler()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(wsHandler, w, r)
	})

	port := "8080"

	log.Printf("Starting server on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
