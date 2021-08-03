package crud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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

func insert(w http.ResponseWriter, r *http.Request) {
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

func select_all(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func select_one(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//parametro en url
	taskid, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "invalid id")
		return
	}

	for _, task := range tasks {
		if task.Id == taskid {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	fmt.Fprintf(w, "tem not found!")
}

func update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskid, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "invalid id")
		return
	}

	var uptask Task

	repbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "please enter valid data")
	}
	json.Unmarshal(repbody, &uptask)

	for i, task := range tasks {
		if task.Id == taskid {

			tasks = append(tasks[:i], tasks[i+1:]...)
			uptask.Id = taskid
			tasks = append(tasks, uptask)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "task with id %v has been updated", taskid)
			return
		}
	}

	fmt.Fprintf(w, "tem not found!")
}

func delete_one(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//parametro en url
	taskid, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "invalid id")
		return
	}

	for i, task := range tasks {
		if task.Id == taskid {
			//append(conserva lo q esta antes de la i, y unelo con lo que está despues de la i)
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "task with id %v removed", task.Id)
			return
		}
	}

	fmt.Fprintf(w, "tem not found!")
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my api :)")
}
