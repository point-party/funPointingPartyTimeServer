package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alecthomas/template"
)

func init() {
	fmt.Println("I'm up and running!")
}

func main() {
	flag.Parse()
	tpl := template.Must(template.ParseFiles("index.html"))
	h := newHub()
	router := http.NewServeMux()
	router.Handle("/", homeHandler(tpl))
	router.Handle("/wakeup", wakeUp())
	router.Handle("/ws", wsHandler{h: h})
	log.Printf("serving on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func homeHandler(tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r)
	})
}

func wakeUp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2000)
		w.Write([]byte("API is up and running"))
	})
}
