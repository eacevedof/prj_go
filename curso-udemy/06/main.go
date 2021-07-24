package main

//creando nuestro propio tipo
import "fmt"

type dinero int

var a int
var b dinero

func main() {
	b = 1000
	fmt.Printf("valor y tipo de b: %v, %T\n", b, b)
	//valor y tipo de b: 1000, main.dinero

	//si hago esto daría error:
	//a = b
	//cannot use b (type dinero) as type int in assignment

}

