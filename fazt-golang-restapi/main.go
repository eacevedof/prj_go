package main

import (
	"fmt"
	"log"
	"net/http"
	"crud"

	"github.com/gorilla/mux"
)

func main() {
	host := "localhost"
	port := "3001"

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", crud.Index)

	// TASKS
	router.HandleFunc("/tasks", crud.SelectAll).Methods("GET")
	router.HandleFunc("/tasks", crud.Insert).Methods("POST")
	router.HandleFunc("/tasks/{id}", crud.SelectOne).Methods("GET")
	router.HandleFunc("/tasks/{id}", crud.DeleteOne).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", crud.Update).Methods("PUT")


	//con esto salta el warning de conexiones entrantes
	//log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), router))
}
