// nombre del paquete
package main

import "fmt"

var a int
var b string = "programa"
var c bool
var d bool = true

func main() {
	e := 42
	f := "Hola mundo"

	//raw string literal
	g := `String multi
	linea
	`
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
	fmt.Println(e)
	fmt.Println(f)
	fmt.Println(g)
}
