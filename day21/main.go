package main

import (
	"bufio"
	"regexp"
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

var (
	opRe  = regexp.MustCompile(`^(\w+):\s*(\w+)\s*([+/*-])\s*(\w+)$`)
	numRe = regexp.MustCompile(`^(\w+):\s*(\d+)$`)
)

type Node struct {
	value       int
	left, right string
	op          string
}

type Tree map[string]Node

func ParseInput(lines []string) (Tree, error) {
	result := Tree{}
	for _, line := range lines {
		if parts := opRe.FindStringSubmatch(line); parts != nil {
			result[parts[1]] = Node{left: parts[2], right: parts[4], op: parts[3]}
		} else if parts := numRe.FindStringSubmatch(line); parts != nil {
			val, err := strconv.Atoi(parts[2])
			if err != nil {
				return nil, err
			}
			result[parts[1]] = Node{value: val}
		} else {
			panic("BAD LINE " + line)
		}
	}

	return result, nil
}

func HasHuman(tree Tree, name string) bool {
	if name == "humn" {
		return true
	}
	node := tree[name]
	if node.op == "" {
		return false
	}
	return HasHuman(tree, node.left) || HasHuman(tree, node.right)
}

func Compute(tree Tree, name string) int {
	node := tree[name]
	if node.op == "" {
		return node.value
	}
	v1 := Compute(tree, node.left)
	v2 := Compute(tree, node.right)
	switch node.op {
	case "+":
		return v1 + v2
	case "-":
		return v1 - v2
	case "*":
		return v1 * v2
	case "/":
		return v1 / v2
	}
	panic("Bad Opcode " + node.op)
}

func Invert(tree Tree, name string, value int) int {
	if name == "humn" {
		return value
	}
	node := tree[name]

	fixed := node.right
	variable := node.left
	if HasHuman(tree, fixed) {
		fixed, variable = variable, fixed
	}

	computed := Compute(tree, fixed)
	if name == "root" {
		return Invert(tree, variable, computed)
	}

	// if fixed == node.left {
	// 	fmt.Println("INVERT ", name, " | ", computed, node.op, variable, " = ", value)
	// } else {
	// 	fmt.Println("INVERT ", name, " | ", variable, node.op, computed, " = ", value)
	// }
	find := 0
	switch node.op {
	case "+":
		find = value - computed
	case "-":
		if fixed == node.left {
			find = computed - value
		} else {
			find = value + computed
		}
	case "*":
		find = value / computed
	case "/":
		if fixed == node.left {
			find = value / computed
		} else {
			find = computed * value
		}
	default:
		panic("BAD OPCODE " + node.op)
	}

	// fmt.Println("   STEP ", variable, find)

	return Invert(tree, variable, find)
}

func PartOneSolution(lines []string) (int, error) {
	input, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	return Compute(input, "root"), nil
}

func PartTwoSolution(lines []string) (int, error) {
	input, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	result := Invert(input, "root", 0)

	return result, nil
}
