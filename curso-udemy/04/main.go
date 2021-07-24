// nombre del paquete
package main

import "fmt"

// esto es una declaracion, a nivel de paquete, hay que evitar vars a nivel de paquete
var z int 
//z = 21 //es una expresion de asignación (no una declaracion) por lo tanto da error
// usamos var cuando deseamos inicializar con su valor neutro, 0, "", false

func main() {
	fmt.Println(z)
}
