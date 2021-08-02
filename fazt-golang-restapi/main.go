package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Task struct {
	id      int    `json:id`
	name    string `json:name`
	content string `json:content`
}

type Tasks []Task

var tasks = Tasks{
	{
		id:      1,
		name:    "task one",
		content: "some content",
	},
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", YourHandler)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
