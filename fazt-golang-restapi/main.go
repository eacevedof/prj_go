package main

import (
	"fmt", "mux"
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

func main() {
	r:= mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/products", ProductHandler)
	r.HandleFunc("/articles", ArticleHandler)
	fmt.Println("Hello world")
}
