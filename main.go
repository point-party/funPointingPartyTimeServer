package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

var rooms = make(map[*Room]bool)

// HubHandler stores the hub
type HubHandler struct {
	hub *Hub
}

func (h *Hub) generateRoom(w http.ResponseWriter, r *http.Request) {
	room := createRoom(h)
	n := roomName{Name: room.name}
	go room.start()
	h.listRooms()
	res, err := json.Marshal(n)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (h *Hub) listRoomsAndClients(w http.ResponseWriter, r *http.Request) {
	h.listRooms()
	for _, v := range h.rooms {
		v.listClients()
	}
}

func main() {
	h := HubHandler{hub: newHub()}
	go h.hub.run()
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/wakeup", wakeUp)
	mux.HandleFunc("/generateRoom", h.hub.generateRoom)
	mux.HandleFunc("/listRoomsAndClients", h.hub.listRoomsAndClients)
	mux.HandleFunc("/joinRoom", func(w http.ResponseWriter, r *http.Request) {
		roomName := r.URL.Query().Get("room")
		playerName := r.URL.Query().Get("name")
		fmt.Println("roomName", roomName)
		fmt.Println("playerName", playerName)
		joinRoom(h.hub, roomName, playerName, w, r)
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

type roomName struct {
	Name string `json:"roomName"`
}
