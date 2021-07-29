package main

import (
	"fmt"
)

func main() {
	for i := 0; i<= 100; i++ {
		fmt.Printf("%v\n", i)
		for j:=0; j<3; j++ {
			fmt.Printf("\t%v\n", j)
		}
	}
}
