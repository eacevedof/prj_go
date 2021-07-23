// nombre del paquete
package main

import (
	"fmt"
)

func main() {
	nb, err := fmt.Println("hola mundo")

	//blank identifer: le indica al compilador que deseche lo que retorna, forma incorrecta asignacion := (no new variables on left side of :=), correcta =
	_, _  = fmt.Println(nb, err) 
}