// Package server contains the implementation structs and functions for all routing.
package server

import (
	"fmt"
	"funPointingPartyTime/socketroom"
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
	s.Router.HandleFunc("/generateRoom", h.GenerateRoom())
	s.Router.HandleFunc("/listRoomsAndClients", h.ListRoomsAndClients())
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
		fmt.Println("roomName", roomName)
		fmt.Println("playerName", playerName)
		fmt.Println("observer", observer)
		socketroom.JoinRoom(h, roomName, playerName, observer, w, r)
		fmt.Println("joined room")
	}
}
