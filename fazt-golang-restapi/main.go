package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Task struct {
	id      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Tasks []Task

var tasks = Tasks{
	{
		id:      0,
		Name:    "some name 0",
		Content: "some-0-content",
	},
	{
		id:      1,
		Name:    "task one 1",
		Content: "some content 1",
	},
}

func get_tasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my api :)")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/tasks", get_tasks)

	//con esto salta el warning de conexiones entrantes
	//log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
	log.Fatal(http.ListenAndServe("localhost:8000", router))
}
