package main

import (
	"fmt"
	"log"
	"os"

	"github.com/apogeeoak/dircmp/compare"
)

func main() {
	original := os.Args[1]
	compared := os.Args[2]

	fmt.Println("Searching through", compared)

	err := compare.Compare(original, compared)
	if err != nil {
		log.Fatalln(err)
	}
}
