package main

import (
	"fmt"
	"os"
)

// Answer: 209603
func main() {
	fd, err := os.Open("input1.data")
	if err != nil {
		panic(err)
	}
	values, err := process(fd)
	if err != nil {
		panic(err)
	}

	sum := 0
	for idx := len(values) - 3; idx < len(values); idx += 1 {
		sum += values[idx]
	}
	fmt.Println(sum)
}
