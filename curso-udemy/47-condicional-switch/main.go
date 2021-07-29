package main

import (
	"fmt"
)

func main() {
	switch {
	case 2==4:
		fmt.Println("no deberia imprimir")
	case 3==3:
		fmt.Println("deberia imprimir")
		fallthrough //se indica que el siguiente case se ejecute
	case 4==5:
		fmt.Println("no deberia")		

	}
}

