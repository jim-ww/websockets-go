package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	// _ "github.com/gorilla/websocket"
)

var templ = template.Must(template.ParseFiles("static/chat.html"))
var msgs = []string{"Hello", "World"}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	http.Handle("/static/",http.StripPrefix("/static/", http.FileServer(http.Dir("/static/"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templ.Execute(w, nil)
	})
	http.HandleFunc("/echo", echo)
	log.Printf("starting server on port %s", ":8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	fmt.Println("Upgraded conncection")
	defer conn.Close()
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)
		err = conn.WriteMessage(msgType, msg)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
