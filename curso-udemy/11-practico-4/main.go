//practio-4 tipos primitivos
package main

import "fmt"

type (
	A1 = string
	A2 = A1
)

type (
	B1 = string
	B2 B1
	B3 []B1
	B4 B3
)

type (
	numero int
)

var x numero

func main() {
	fmt.Println(x)
	fmt.Printf("el tipo de x es: %T\n",x)
	x = 42
	fmt.Println(x)
}

