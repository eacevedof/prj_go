package main

import (
	"fmt"
)

type Task struct {
	id      int    `json:id`
	name    string `json:name`
	content string `json:content`
}

type alltasks []Task

func main() {
	fmt.Println("Hello world")
}
