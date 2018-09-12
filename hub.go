package main

import (
	"fmt"
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

func newHub() *Hub {
	return &Hub{
		register:   make(chan *Room),
		unregister: make(chan *Room),
		rooms:      make(map[string]*Room),
	}
}

func (h *Hub) run() {
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

func (h *Hub) listRooms() {
	for k := range h.rooms {
		fmt.Println("Rooms", k)
	}
}
