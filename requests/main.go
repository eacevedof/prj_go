package main

import (
	"fmt"
	"log"
	"net/http"
	"request"

	"github.com/gorilla/mux"
)

func main() {
	host := "localhost"
	port := "3001"

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", request.Index)

	// TASKS
	router.HandleFunc("/tasks", request.SelectAll).Methods("GET")
	router.HandleFunc("/tasks", request.Insert).Methods("POST")
	router.HandleFunc("/tasks/{id}", request.SelectOne).Methods("GET")
	router.HandleFunc("/tasks/{id}", request.DeleteOne).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", request.Update).Methods("PUT")


	//con esto salta el warning de conexiones entrantes
	//log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), router))
}
