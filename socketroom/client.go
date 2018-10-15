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
		fmt.Println("gameMessage", gameMessage)
		determineGameAction(c, &gameMessage, msg)
		if err != nil {
			fmt.Println("REALLY ENCOUNTERED AN ERROR")
			log.Printf("error getting json message: %v", err)
		}
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.Room.broadcast <- gameMessage
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
			fmt.Println("How does this work???", gameMessage)
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
func JoinRoom(hub *Hub, roomName string, clientName string, observer string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	room, err := findRoom(hub, roomName)
	if err != nil {
		fmt.Println(err)
	}
	client := &Client{
		Room:     room,
		conn:     conn,
		Name:     clientName,
		observer: determineObserver(observer),
		send:     make(chan GameMessage),
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
		}
		return room, nil
	}
	return nil, fmt.Errorf("Could not find room")
}

func determineObserver(observer string) bool {
	if observer == "true" {
		return true
	}
	return false
}
