// Package socketroom contains the implementation structs and functions to create a multiple rooms for users.
package socketroom

import (
	"fmt"
	"net/http"
)

// Hub maintains the set of active rooms.
type Hub struct {
	// Registered rooms.
	Rooms map[string]*Room

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
		Rooms:      make(map[string]*Room),
	}
}

// Run initiates the select that listens to the channels.
func (h *Hub) Run() {
	for {
		select {
		case room := <-h.register:
			h.Rooms[room.Name] = room
		case room := <-h.unregister:
			if _, ok := h.Rooms[room.Name]; ok {
				delete(h.Rooms, room.Name)
			}
		}
	}
}

// ListRooms is a helper function to print what rooms are currently registered with the hub.
func (h *Hub) ListRooms() {
	for k := range h.Rooms {
		fmt.Println("Rooms", k)
	}
}

//ListRoomsAndClients is a helper endpoint that lists all rooms registerd with the hub and all the clients registered in the rooms.
func (h *Hub) ListRoomsAndClients() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ListRooms()
		for _, v := range h.Rooms {
			v.ListClients()
		}
		w.Write([]byte{})
	}
}
