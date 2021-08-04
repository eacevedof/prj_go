package main

import (
    "fmt"
	"request"
)

func main() {
    strproducts := request.FetchProducts()
    fmt.Println(strproducts)
}
