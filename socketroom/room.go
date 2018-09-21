package socketroom

import (
	"fmt"
	"math/rand"
	"time"
)

// Room will be the place clients use to create a pointing session.
type Room struct {
	// The Hub handles all rooms.
	hub *Hub

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan GameMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	Name string
}

// RoomName contains the name in json format to send to the client.
type RoomName struct {
	Name string `json:"roomName"`
}

// CreateRoom creates a new room and registers it with the hub.
func CreateRoom(hub *Hub) *Room {
	room := &Room{
		hub:        hub,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Name:       createRoomName(),
		broadcast:  make(chan GameMessage),
	}
	room.hub.register <- room
	return room
}

// Start begins the goroutine and channels for the room
func (r *Room) Start() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
			}
		case gameMessage := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- gameMessage:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}

// ListClients prints to console the clients in the room
func (r *Room) ListClients() {
	for k := range r.clients {
		fmt.Println("Clients", k)
	}
}

// Logic to create random room name
const charset = "abcdefghijklmnopqrstuvwxyz"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func createRoomName() string {
	return stringWithCharset(6, charset)
}
