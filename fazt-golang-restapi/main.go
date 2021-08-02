package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Task struct {
	Id      int    `json:"Id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Tasks []Task

var tasks = Tasks{
	{
		Id:      1,
		Name:    "task one 1",
		Content: "some content 1",
	},
}

func create_task(w http.ResponseWriter, r *http.Request) {
	var newtask Task
	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "insert a valid task")
	}

	json.Unmarshal(reqbody, &newtask)

	newtask.Id = len(tasks) + 1
	tasks = append(tasks, newtask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newtask)
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
	router.HandleFunc("/tasks", create_task).Methods("POST")
	router.HandleFunc("/tasks", get_tasks)

	//router.HandleFunc("/tasks/{id}", getOneTask).Methods("GET")
	//router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	//router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	//con esto salta el warning de conexiones entrantes
	//log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
	log.Fatal(http.ListenAndServe("localhost:8000", router))
}
