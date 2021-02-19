package main

import (
	"fmt"
	"log"

	"github.com/apogeeoak/dircmp/compare"
)

func main() {
	config := compare.ParseConfig()
	fmt.Println("Searching through", config.Compared)

	stat, err := compare.Compare(config)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(stat)
}
