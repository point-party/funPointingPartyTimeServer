package main

import (
	"fmt"
)

// Hub maintains the set of active rooms.
type Hub struct {
	// Registered rooms.
	rooms map[*Room]bool

	// Register requests from the rooms.
	register chan *Room

	// Unregister requests from rooms.
	unregister chan *Room
}

func newHub() *Hub {
	return &Hub{
		register:   make(chan *Room),
		unregister: make(chan *Room),
		rooms:      make(map[*Room]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case room := <-h.register:
			h.rooms[room] = true
		case room := <-h.unregister:
			if _, ok := h.rooms[room]; ok {
				delete(h.rooms, room)
			}
		}
	}
}

func (h *Hub) listRooms() {
	for k := range h.rooms {
		fmt.Println("Rooms", k)
	}
}
