package main

import (
	"bufio"
	"fmt"
	"os"
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

	values, err := PartTwoSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2: ", values)
}

type Step struct {
	opcode string
	value  int
}

func ParseInput(lines []string) []Step {
	result := []Step{}
	for _, line := range lines {
		parts := strings.Split(line, " ")
		step := Step{
			opcode: parts[0],
		}
		if len(parts) == 2 {
			value, _ := strconv.ParseInt(parts[1], 10, 64)
			step.value = int(value)
		}
		result = append(result, step)
	}

	return result
}

func Compute(steps []Step) []int {
	value := 1
	result := []int{value}

	for _, step := range steps {
		if step.opcode == "noop" {
			result = append(result, value)
		} else if step.opcode == "addx" {
			result = append(result, value)
			value += step.value
			result = append(result, value)
		}
	}

	return result
}

func getSums(values []int, points ...int) int {
	sum := 0
	for _, idx := range points {
		// fmt.Println(idx, values[idx], values[idx-1])
		// Minus one since we want the value at the start of the cycle
		sum += values[idx-1] * idx
	}

	return sum
}

func PartOneSolution(lines []string) (int, error) {
	steps := ParseInput(lines)
	values := Compute(steps)

	//for idx, v := range values {
	//fmt.Println(idx, v)
	//}

	return getSums(values, 20, 60, 100, 140, 180, 220), nil
}

func PartTwoSolution(lines []string) ([]string, error) {
	steps := ParseInput(lines)
	values := Compute(steps)
	output := [6]string{}

	for idx := 0; idx < 240; idx += 1 {
		val := values[idx]
		diff := val - (idx % 40)
		ch := "."
		if diff < 2 && diff > -2 {
			ch = "#"
		}

		output[(idx)/40] = output[(idx)/40] + ch
	}

	fmt.Println("")
	for _, line := range output {
		fmt.Println(line)
	}
	fmt.Println("")

	return output[:], nil
}
