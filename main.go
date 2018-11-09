package main

import (
	"fmt"
	"funPointingPartyTimeServer/server"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	port := getPort()
	s := server.Server{Router: http.NewServeMux()}
	handler := cors.AllowAll().Handler(s.Router)
	s.Routes()
	err := http.ListenAndServe(port, handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
