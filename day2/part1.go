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
	value, err := process(fd)
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}
