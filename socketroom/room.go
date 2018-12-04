package socketroom

import (
	"fmt"
	"math/rand"
	"time"
)

type voteEvent struct {
	point string
	id    string
}

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

	vote  chan voteEvent
	clear chan bool
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
		vote:       make(chan voteEvent),
		clear:      make(chan bool),
	}
	room.hub.register <- room
	return room
}

func (r *Room) sendPlayers() []PlayerStatus {
	var ps []PlayerStatus
	for p := range r.clients {
		value := PlayerStatus{p.Name, p.CurrentPoint, p.ID}
		ps = append(ps, value)
	}
	return ps
}

// Start begins the goroutine and channels for the room
func (r *Room) Start() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			fmt.Println("registered with room", client.Name)
			joinMsg := GameMessage{
				Event: joinRoom,
				Payload: PlayerUpdate{
					Players: r.sendPlayers(),
				},
			}
			for client := range r.clients {
				client.send <- joinMsg
			}

		case client := <-r.unregister:
			delete(r.clients, client)
			close(client.send)
			exitMsg := GameMessage{
				Event: leaveRoom,
				Payload: PlayerUpdate{
					Players: r.sendPlayers(),
				},
			}
			for client := range r.clients {
				client.send <- exitMsg
			}
		case gameMessage := <-r.broadcast:
			// DECODE JSON here into different stuff -> decide actions
			for client := range r.clients {
				select {
				case client.send <- gameMessage:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		case voteEvent := <-r.vote:
			r.updateVote(voteEvent.point, voteEvent.id)
		case clear := <-r.clear:
			if clear {
				r.clearPoints()
			}
		}
	}
}

func (r *Room) updateVote(point string, id string) {
	for c := range r.clients {
		if c.ID == id {
			c.CurrentPoint = point
		}
	}
}

// ListClients prints to console the clients in the room
func (r *Room) ListClients() {
	for k := range r.clients {
		fmt.Println("Clients", k)
	}
}

func (r *Room) clearPoints() {
	for c := range r.clients {
		c.CurrentPoint = ""
	}
	fmt.Println("All points cleared")
}

// Logic to create random room name
const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

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
