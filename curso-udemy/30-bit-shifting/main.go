//big shifting
package main

import "fmt"

const (
	//cancela el 0
	_ = iota
	kb = 1 << (iota * 10) //iota = 1
	gb = 1 << (iota * 10) //iota = 2
	tb = 1 << (iota * 10) //iota = 3
)

func main() {
	fmt.Printf("%d\t\t\t%b\n", kb, kb)
	fmt.Printf("%d\t\t\t%b\n", gb, gb)
	fmt.Printf("%d\t\t%b\n", tb, tb)
}