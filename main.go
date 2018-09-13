package main

import (
	"fmt"
	"funPointingPartyTime/socketroom"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	h := socketroom.NewHub()
	go h.Run()
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/wakeup", wakeUp)
	mux.HandleFunc("/generateRoom", h.GenerateRoom)
	mux.HandleFunc("/listRoomsAndClients", h.ListRoomsAndClients)
	mux.HandleFunc("/joinRoom", func(w http.ResponseWriter, r *http.Request) {
		roomName := r.URL.Query().Get("room")
		playerName := r.URL.Query().Get("name")
		fmt.Println("roomName", roomName)
		fmt.Println("playerName", playerName)
		socketroom.JoinRoom(h, roomName, playerName, w, r)
		fmt.Println("joined room")
	})
	handler := cors.AllowAll().Handler(mux)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func wakeUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("API is up and running"))
}
