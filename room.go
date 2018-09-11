package main

import "fmt"

// Room will be the place clients use to create a pointing session.
type Room struct {
	// The Hub handles all rooms.
	hub *Hub

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	name string
}

// createRoom creates a new room and registers it with the hub.
func createRoom(hub *Hub, name string) *Room {
	room := &Room{
		hub:        hub,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		name:       name,
		broadcast:  make(chan []byte),
	}
	room.hub.register <- room
	return room
}

func (r *Room) start() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
			}
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}

func (r *Room) listClients() {
	for k := range r.clients {
		fmt.Println("Clients", k)
	}
}
