package main

import (
	"fmt"
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
	fmt.Println("Hello world")
}
