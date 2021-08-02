package main

import (
	"fmt"
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

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my api :)")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)

	//con esto salta el warning de conexiones entrantes
	//log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
	log.Fatal(http.ListenAndServe("localhost:8000", router))
}
