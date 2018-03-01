package main

import (
	"log"
	"sync"
	"time"
)

// hub maintains the set of active clients and broadcasts messages to the clients.

type hub struct {
	// mutex to protect connections
	connectionsMx sync.RWMutex

	// Registered connections.
	rooms map[string][]*connection

	// Inbound messages from the connections.
	broadcast chan []byte

	logMx sync.RWMutex
	log   [][]byte
}

func newHub() *hub {
	h := &hub{
		connectionsMx: sync.RWMutex{},
		broadcast:     make(chan []byte),
		rooms:         make(map[string][]*connection),
	}

	go func() {
		for {
			msg := <-h.broadcast
			// decide how to only send message to connections within room
			h.connectionsMx.RLock()
			for c := range h.connections {
				select {
				case c.send <- msg:
					// stop trying to send to this connection after trying for 1 second.
					// if we have to stop, it means that a reader died so remove the connection also.
				case <-time.After(1 * time.Second):
					log.Printf("shutting down connection %s", c)
					h.removeConnection(c)
				}
			}
			h.connectionsMx.RUnlock()
		}
	}()
	return h
}

func (h *hub) addConnection(conn *connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()
	h.connections[conn] = struct{}{}
}

// make methods for rooms
// add connection to room
// remove connection from room
// add room to hub rooms map
// remove room from hub rooms map

func (h *hub) removeRoom(room string) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()
	if _, ok := h.rooms[room]; ok {
		delete(h.rooms, room)
		close(conn.send)
	}
}

func (h *hub) findConnection(conn *connection) bool {
	for k, _ := range h.rooms {
		for _, v := range h.rooms[k] {
			if v == conn {
				return true
			}
		}
	}
	return false
}
