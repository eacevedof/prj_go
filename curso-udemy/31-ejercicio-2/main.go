package main

import (
	"fmt"
)

const (
	a = 42 //constante sin tipo
	b int = 43 //constante con tipo
)

func main() {
	fmt.Printf("%v %v\n", a, b)
}
