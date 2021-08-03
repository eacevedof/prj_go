package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)

	// TASKS
	router.HandleFunc("/tasks", select_all).Methods("GET")
	router.HandleFunc("/tasks", insert).Methods("POST")
	router.HandleFunc("/tasks/{id}", select_one).Methods("GET")
	router.HandleFunc("/tasks/{id}", delete_one).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", update).Methods("PUT")

	//con esto salta el warning de conexiones entrantes
	//log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
	log.Fatal(http.ListenAndServe("localhost:8000", router))
}
