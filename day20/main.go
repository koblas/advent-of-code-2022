package main

import (
	"bufio"
	"strconv"

	// "errors"
	"fmt"
	"os"
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

	values, err := PartTwoSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2: ", values)
}

func ParseInput(lines []string) ([]int, error) {
	result := []int{}
	for _, line := range lines {
		value, err := strconv.ParseInt(line, 10, 0)
		if err != nil {
			return nil, err
		}

		result = append(result, int(value))
	}

	return result, nil
}

type RingNode struct {
	value            int
	next, prev, step *RingNode
}

func Mix(input []int, rounds int) []int {
	lengthMinusOne := len(input) - 1

	var root, current *RingNode
	for _, value := range input {
		value := &RingNode{
			value: value,
			prev:  current,
		}
		if current != nil {
			current.next = value
			current.step = value
		}
		current = value
		if root == nil {
			root = value
		}
	}
	root.prev = current
	current.next = root

	for round := 0; round < rounds; round += 1 {
		for item := root; item != nil; item = item.step {
			if item.value == 0 {
				continue
			}

			// Remove from ring
			pItem := item.prev
			nItem := item.next
			pItem.next = nItem
			nItem.prev = pItem

			var pos *RingNode
			shift := item.value % lengthMinusOne
			if shift > 0 {
				pos = nItem
				for v := 1; v < shift; v += 1 {
					pos = pos.next
				}
			} else {
				pos = pItem
				for v := 0; v > shift; v -= 1 {
					pos = pos.prev
				}
			}

			pos.next.prev = item
			item.next = pos.next
			item.prev = pos
			pos.next = item
		}
	}

	result := []int{root.value}
	for item := root.next; item != root; item = item.next {
		result = append(result, item.value)
	}

	return result
}

func Sum(values []int) int {
	length := len(values)
	zeroIdx := 0
	for values[zeroIdx] != 0 {
		zeroIdx += 1
	}

	return values[(zeroIdx+1000)%length] + values[(zeroIdx+2000)%length] + values[(zeroIdx+3000)%length]
}

func PartOneSolution(lines []string) (int, error) {
	input, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	mixed := Mix(input, 1)

	return Sum(mixed), err
}

func PartTwoSolution(lines []string) (int, error) {
	input, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	for idx := 0; idx < len(input); idx += 1 {
		input[idx] *= 811589153
	}

	mixed := Mix(input, 10)

	return Sum(mixed), err
}
