// nombre del paquete
package main

import "fmt"

// con var definimos a nivel global
var z = 41

func main() {
	//var inicializa w con 0, se puede usar dentro de func pero no es buena practica
	var w int

	//el operador corto := se usa dentro de bloques, dentro de funcioens 
	x := 42 + 7
	y := "james bond"
	fmt.Println(x)
	fmt.Println(y)
	x = 50
	fmt.Println(x)
	fmt.Println(z)
	fmt.Println("w:",w)
	numero()
}

func numero() {

	fmt.Println(z)
}