package main

import (
	"bufio"
	"fmt"
	"os"
	// "strings"
)

func main() {
	fd, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(fd)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	value, err := PartOneSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1: ", value)

	value, err = PartTwoSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2: ", value)
}

func seenCheck(seen []byte, value byte) ([]byte, bool) {
	// fmt.Println("Checking ", value, " IN ", seen)
	for _, val := range seen {
		if val == value {
			return seen, true
		}
	}

	seen = append(seen, value)
	return seen, false
}

func findFirst(line string, length int) int {
	part := line[0:length]
	for idx, c := range line[length:] {
		seen := []byte{}
		found := false
		for _, cpart := range part {
			seen, found = seenCheck(seen, byte(cpart))
			if found {
				break
			}
		}
		if !found {
			return idx + length
		}
		ch := string(c)
		part = part[1:] + ch
	}

	return -1
}

func PartOneSolution(lines []string) (int, error) {
	result := 0

	for _, line := range lines {
		value := findFirst(line, 4)

		if value > result {
			result = value
		}
	}

	return result, nil
}

func PartTwoSolution(lines []string) (int, error) {
	result := 0

	for _, line := range lines {
		value := findFirst(line, 14)

		if value > result {
			result = value
		}
	}

	return result, nil
}
