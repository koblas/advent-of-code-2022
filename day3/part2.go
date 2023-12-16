package main

import (
	"fmt"
	"os"
)

func main() {
	fd, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	value, err := processPart2(fd)
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}
