package main

// conversion de tipos, No casting 
import "fmt"

type dinero int

var a int
var b dinero

func main() {
	b = 23
	fmt.Printf("valor y tipo de b: %v, %T\n", b, b)

	a = int(b)
	fmt.Printf("a: %v, %T\n", a,a)

}

