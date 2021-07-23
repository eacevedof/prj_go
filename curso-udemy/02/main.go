// nombre del paquete
package main

import "fmt"


//y := 40 // esto daría error (outside func body)

func main() {
	// ... interface{} //todos los tipos implmenetan la interfaz vacía
	//por eso println acepta cualquier tipo como parametro

	// := operador corto de declaraciones cortas, solo se puede usar en el cuerpo de la func no fuera
	// := inicia la variable, si se desea reasignar hay que usar = 
	x := 42
	y := 2 + 4
	x = 5
	fmt.Println(x, y)
}