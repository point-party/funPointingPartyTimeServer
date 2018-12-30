package socketroom

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Room     *Room
	Name     string
	observer bool
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send         chan GameMessage
	CurrentPoint string
	ID           string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.Room.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg json.RawMessage
		gameMessage := GameMessage{
			Payload: &msg,
		}
		err := c.conn.ReadJSON(&gameMessage)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// pass decoded payload -> msg into wrapped GameMessage
		gameEvent := GameEvent{
			gameMessage: gameMessage,
			rawPayload:  msg,
		}
		c.Room.broadcast <- gameMessage
		c.Room.determineGameAction <- gameEvent
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case gameMessage, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteJSON(gameMessage)

			if err != nil {
				log.Printf("could not send json correctly: %v", err)
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// JoinRoom handles inserting a client into a room and upgrading to WS.
func JoinRoom(hub *Hub, roomName string, clientName string, role string, id string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error during connection upgrade: %e", err)
		return
	}
	room, err := findRoom(hub, roomName)
	if err != nil {
		fmt.Printf("Error finding room: %e", err)
		return
	}
	client := &Client{
		Room:     room,
		conn:     conn,
		Name:     clientName,
		observer: determineObserver(role),
		send:     make(chan GameMessage),
		ID:       id,
	}

	client.Room.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func findRoom(hub *Hub, name string) (*Room, error) {
	var room *Room
	for k, v := range hub.Rooms {
		if k == name {
			room = v
			return room, nil
		}
	}
	return nil, fmt.Errorf("Could not find room")
}

func determineObserver(role string) bool {
	if role == "OBSERVER" {
		return true
	}
	return false
}
