// Package server contains the implementation structs and functions for all routing.
package server

import (
	"encoding/json"
	"fmt"
	"funPointingPartyTimeServer/authentication"
	"funPointingPartyTimeServer/socketroom"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

// Server will contain all structs such as router, db, config
type Server struct {
	Router *http.ServeMux
}

// Routes contains all available routes
func (s *Server) Routes() {
	h := socketroom.NewHub()
	sS := authentication.StateString(10)
	endpoint := oauth2.Endpoint{
		AuthURL:  "https://auth.atlassian.com/authorize?audience=api.atlassian.com",
		TokenURL: "https://auth.atlassian.com/oauth/token?audience=api.atlassian.com",
	}
	config := oauth2.Config{
		ClientID:     os.Getenv("JIRA_CLIENT_ID"),
		ClientSecret: os.Getenv("JIRA_SECRET_ID"),
		Endpoint:     endpoint,
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"read:jira-work"},
	}
	url := config.AuthCodeURL(sS)
	go h.Run()
	s.Router.Handle("/", http.FileServer(http.Dir("./static")))
	s.Router.HandleFunc("/wakeup", s.wakeup())
	s.Router.HandleFunc("/generateRoom", s.generateRoom(h))
	s.Router.HandleFunc("/listRoomsAndClients", s.listRoomsAndClients(h))
	s.Router.HandleFunc("/joinRoom", s.joinRoom(h))
	s.Router.HandleFunc("/login", s.login(url))
	s.Router.HandleFunc("/callback", s.callback(config, sS))
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
		defer r.Body.Close()
		playerName := r.URL.Query().Get("name")
		role := r.URL.Query().Get("role")
		id := r.URL.Query().Get("id")
		socketroom.JoinRoom(h, roomName, playerName, role, id, w, r)
		fmt.Println("joined room")
	}
}

func (s *Server) generateRoom(h *socketroom.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pS := r.URL.Query().Get("pointScale")
		room := socketroom.CreateRoom(h, pS)
		defer r.Body.Close()
		n := socketroom.RoomName{Name: room.Name}
		go room.Start()
		res, err := json.Marshal(n)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	}
}

func (s *Server) callback(config oauth2.Config, state string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		stateResp := r.FormValue("state")
		if stateResp != state {
			http.Error(w, "State code doesn't match", 400)
			return
		}
		code := r.FormValue("code")
		token, err := config.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Println("err", err)
			http.Error(w, "Couldn't get token from code", 500)
			return
		}
		resp := fmt.Sprintf("State: %s \n, Token: %s", stateResp, token)
		w.Write([]byte(resp))
	}
}

func (s *Server) login(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, 301)
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
