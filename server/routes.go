// Package server contains the implementation structs and functions for all routing.
package server

import (
	"encoding/json"
	"fmt"
	"funPointingPartyTimeServer/socketroom"
	"net/http"
)

// Server will contain all structs such as router, db, config
type Server struct {
	Router *http.ServeMux
}

// Routes contains all available routes
func (s *Server) Routes() {
	h := socketroom.NewHub()
	go h.Run()
	s.Router.Handle("/", http.FileServer(http.Dir("./static")))
	s.Router.HandleFunc("/wakeup", s.wakeup())
	s.Router.HandleFunc("/generateRoom", s.generateRoom(h))
	s.Router.HandleFunc("/listRoomsAndClients", s.listRoomsAndClients(h))
	s.Router.HandleFunc("/joinRoom", s.joinRoom(h))
}

func (s *Server) wakeup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("API is up and running"))
	}
}

func (s *Server) joinRoom(h *socketroom.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomName := r.URL.Query().Get("room")
		playerName := r.URL.Query().Get("name")
		observer := r.URL.Query().Get("observer")
		id := r.URL.Query().Get("id")
		socketroom.JoinRoom(h, roomName, playerName, observer, id, w, r)
		fmt.Println("joined room")
	}
}

func (s *Server) generateRoom(h *socketroom.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		room := socketroom.CreateRoom(h)
		n := socketroom.RoomName{Name: room.Name}
		go room.Start()
		h.ListRooms()
		res, err := json.Marshal(n)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	}
}

func (s *Server) listRoomsAndClients(h *socketroom.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ListRooms()
		for _, v := range h.Rooms {
			v.ListClients()
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Done"))
	}
}
