package main

import (
	"fmt"
)

func main() {
	switch "manzana" {
	case "pera", "limon":
		fmt.Println("no deberia imprimir")
	case "manzana", "cierula", "fresas":
		fmt.Println("rojas")
		fallthrough //se indica que el siguiente case se ejecute
	default:
		fmt.Println("default")		

	}
}

