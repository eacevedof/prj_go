// nombre del paquete
package main

import "fmt"

// esto es una declaracion, a nivel de paquete, hay que evitar vars a nivel de paquete
var z int = 21
//z = 21 //es una expresion de asignación (no una declaracion) por lo tanto da error


func main() {
	fmt.Println(z)
}
