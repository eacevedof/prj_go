//big shifting
package main

import "fmt"

const (
	//cancela el 0
	_ = iota
	// bit shifting = rellena con 0
	kb = 1 << (iota * 10) //iota = 1
	gb = 1 << (iota * 10) //iota = 2
	tb = 1 << (iota * 10) //iota = 3
)

func main() {

	a := 4
	fmt.Printf("%d\t\t\t%b\n", a, a) //esto imprime 4 	100  4 y su rep en binario

	// nuevo bits en 1 posición hacia la izquierda. Se coge el bit más significativo y se desplaza una pos hacia la iz
	b := a << 1
	fmt.Printf("%d\t\t\t%b\n", b, b) //esto imprime 8 	1000 se desplza 1 hacia la izq y se completa con 0	

	fmt.Printf("%d\t\t\t%b\n", kb, kb)
	fmt.Printf("%d\t\t\t%b\n", gb, gb)
	fmt.Printf("%d\t\t%b\n", tb, tb)
}