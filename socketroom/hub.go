// Package socketroom contains the implementation structs and functions to create a multiple rooms for users.
package socketroom

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Hub maintains the set of active rooms.
type Hub struct {
	// Registered rooms.
	rooms map[string]*Room

	// Register requests from the rooms.
	register chan *Room

	// Unregister requests from rooms.
	unregister chan *Room
}

//NewHub returns a new Hub.
func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Room),
		unregister: make(chan *Room),
		rooms:      make(map[string]*Room),
	}
}

// Run initiates the select that listens to the channels.
func (h *Hub) Run() {
	for {
		select {
		case room := <-h.register:
			h.rooms[room.name] = room
		case room := <-h.unregister:
			if _, ok := h.rooms[room.name]; ok {
				delete(h.rooms, room.name)
			}
		}
	}
}

// ListRooms is a helper function to print what rooms are currently registered with the hub.
func (h *Hub) ListRooms() {
	for k := range h.rooms {
		fmt.Println("Rooms", k)
	}
}

// GenerateRoom creates a new room and returns the room name to the client.
func (h *Hub) GenerateRoom(w http.ResponseWriter, r *http.Request) {
	room := CreateRoom(h)
	n := RoomName{Name: room.name}
	go room.start()
	h.ListRooms()
	res, err := json.Marshal(n)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

//ListRoomsAndClients is a helper endpoint that lists all rooms registerd with the hub and all the clients registered in the rooms.
func (h *Hub) ListRoomsAndClients(w http.ResponseWriter, r *http.Request) {
	h.ListRooms()
	for _, v := range h.rooms {
		v.listClients()
	}
}
