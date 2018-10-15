package main

import (
	"funPointingPartyTimeServer/server"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	s := server.Server{Router: http.NewServeMux()}
	handler := cors.AllowAll().Handler(s.Router)
	s.Routes()
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
