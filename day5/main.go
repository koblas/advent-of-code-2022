package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

var moveRe = regexp.MustCompile("move (\\d+) from (\\d+) to (\\d+)")

type Stack []string
type Move struct {
	count, from, to int64
}

func ParseInput(lines []string) ([]Stack, []Move) {
	moves := []Move{}
	stacks := []Stack{}
	// stacks := make([]Stack, (len(lines[0])/4)+1)
	// fmt.Println(lines[0])
	// fmt.Println("LEN ", len(lines[0]), (len(lines[0])+1)/4)
	// fmt.Println(stacks)

	for _, line := range lines {
		if strings.Contains(line, "[") {
			for idx := 0; idx < len(line); idx += 3 {
				idx += 1
				if line[idx] != ' ' {
					pos := idx / 4
					for pos >= len(stacks) {
						stacks = append(stacks, Stack{})
					}
					// fmt.Println("ADD TO ", idx/4)
					stacks[pos] = append(stacks[pos], string(line[idx]))
				}
			}
		} else if parts := moveRe.FindStringSubmatch(line); parts != nil {
			// fmt.Println("HERE ", line, parts)
			count, _ := strconv.ParseInt(parts[1], 10, 64)
			from, _ := strconv.ParseInt(parts[2], 10, 64)
			to, _ := strconv.ParseInt(parts[3], 10, 64)
			moves = append(moves, Move{count, from - 1, to - 1})
		}
	}

	// fmt.Println(stacks)

	return stacks, moves
}

func PrintStack(stacks []Stack) {
	maxDepth := 0
	for _, stack := range stacks {
		if maxDepth < len(stack) {
			maxDepth = len(stack)
		}
	}

	for depth := 0; depth < maxDepth; depth += 1 {
		for _, stack := range stacks {
			// fmt.Println("D: ", depth, len(stack))
			if depth < len(stack) {
				fmt.Printf("[%s] ", stack[depth])
			} else {
				fmt.Printf("    ")
			}
		}
		fmt.Printf("\n")
	}
}

func PartOneSolution(lines []string) (string, error) {
	stacks, moves := ParseInput(lines)

	PrintStack(stacks)

	for _, move := range moves {
		for idx := int64(0); idx < move.count; idx += 1 {
			item := stacks[move.from][0]
			stacks[move.from] = stacks[move.from][1:]
			stacks[move.to] = append(Stack{item}, stacks[move.to]...)
		}
	}

	result := ""
	for _, stack := range stacks {
		result = result + stack[0]
	}

	return result, nil
}

func PartTwoSolution(lines []string) (string, error) {
	stacks, moves := ParseInput(lines)

	fmt.Println("START === ")
	PrintStack(stacks)

	for _, move := range moves {
		// Must deeply copy the items
		items := append([]string{}, stacks[move.from][0:move.count]...)
		stacks[move.from] = stacks[move.from][move.count:]
		stacks[move.to] = append(items, stacks[move.to]...)

		// OLD OLD old := stacks[move.to]
		// OLD OLD stacks[move.to] = append(Stack{}, items...)
		// OLD OLD stacks[move.to] = append(stacks[move.to], old...)

		// fmt.Println("MOVE === ", move.count, move.from, move.to)
		// PrintStack(stacks)
	}

	result := ""
	for _, stack := range stacks {
		result = result + stack[0]
	}

	return result, nil
}
