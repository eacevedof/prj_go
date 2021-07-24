// nombre del paquete
package main

//alclarando el pquete fmt
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

	print_format()
}

func print_format() {
	//%v verbo que apunta al valor en formato raw
	fmt.Printf("valor de la variable a es %v\n",a)
	fmt.Printf("valor de la variable a es %d\n",a) //flag formato decimales
	fmt.Printf("el valor de b es: %v \n",b)
	fmt.Printf("el valor de b es: %s \n",b) //fla formato string
	fmt.Printf("tipo de a es %T\n",a)
	fmt.Printf("tipo de b es %T\n",b)

	s1 := fmt.Sprint("el ",b," dice hola mundo.")
	fmt.Println(s1)
	fmt.Printf("%T",s1)

}