package main

import (
	"fmt"
	"os"
)

// Answer: 71506
func main() {
	fd, err := os.Open("input1.data")
	if err != nil {
		panic(err)
	}
	values, err := process(fd)
	if err != nil {
		panic(err)
	}
	fmt.Println(values[len(values)-1])
}
