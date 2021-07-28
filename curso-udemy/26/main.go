//tipo string parte 2
package main

import "fmt"

func main() {
	s1 := "hola mundo"
	s2 := `Esta es una 
	linea partida`

	fmt.Println(s1)
	fmt.Printf("T de s1:%T\n",s1)
	fmt.Println(s2)
	fmt.Printf("T de s2:%T\n",s2)

	fmt.Println()

	bs := []byte(s1)
	fmt.Println(bs)
	fmt.Printf("T de byte: %T\n", bs)
	fmt.Println("")

	for i:= 0; i < len(s1); i++ {
		fmt.Printf("%#U\n", s1[i])
	}

}

