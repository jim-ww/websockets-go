package main

import (
	"fmt"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	homeTmpl.ExecuteTemplate(w, "home.html", struct{ Messages []string }{messages})
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	messagesTmpl.ExecuteTemplate(w, "messages.html", struct{ Messages []string }{messages})
}

func postMessages(w http.ResponseWriter, r *http.Request) {
	messages = append(messages, fmt.Sprintf("%s: %s", r.FormValue("user"), r.FormValue("message")))
	messagesTmpl.ExecuteTemplate(w, "messages.html", struct{ Messages []string }{messages})
}

func deleteMessages(w http.ResponseWriter, r *http.Request) {
	messages = []string{}
	messagesTmpl.ExecuteTemplate(w, "messages.html", struct{ Messages []string }{messages})
}
