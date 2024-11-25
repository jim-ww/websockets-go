package main

import (
	"log"
	"net/http"
	"text/template"
	// _ "github.com/gorilla/websocket"
)

var homeTmpl = template.Must(template.ParseFiles("./static/home.html"))
var messagesTmpl = template.Must(template.ParseFiles("./static/messages.html"))
var messages = []string{"bob:Hello", "bob:World"}

func main() {
	http.HandleFunc("GET /", home)
	http.HandleFunc("GET /messages", getMessages)
	http.HandleFunc("POST /messages", postMessages)
	http.HandleFunc("DELETE /messages", deleteMessages)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// func sendRecieve(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	for {
// 		messageType, msg, err := conn.ReadMessage()
// 	}
// }
