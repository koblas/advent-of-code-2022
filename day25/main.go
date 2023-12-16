package main

import (
	"bufio"
	"time"

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

	start := time.Now()
	value, err := PartOneSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Part 1 (in %s): %v\n", time.Since(start), value)

	start = time.Now()
	values, err := PartTwoSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Part 2 (in %s): %v\n", time.Since(start), values)
}

var CharToValue = map[rune]int{
	'=': -2,
	'-': -1,
	'0': 0,
	'1': 1,
	'2': 2,
}

var ValueToChar = map[int]string{
	-2: "=",
	-1: "-",
	0:  "0",
	1:  "1",
	2:  "2",
}

type Value []int

func (value Value) ToInt() int {
	sum := 0
	for _, v := range value {
		sum = (sum * 5) + v
	}
	return sum
}

func (value Value) ToString() string {
	result := ""
	for _, v := range value {
		result += ValueToChar[v]
	}
	return result
}

func FromString(input string) Value {
	value := Value{}
	for _, ch := range input {
		value = append(value, CharToValue[rune(ch)])
	}
	return value
}

func FromInt(input int) Value {
	value := Value{}
	if input == 0 {
		return Value{0}
	}

	// fmt.Println("START ", input)

	// carry := 0
	for input != 0 {
		v := (input % 5)
		digit := v
		if v > 2 {
			digit = v - 5
			input += 5
		}
		// fmt.Println("  HERE ", input, v, digit, carry)
		value = append(Value{digit}, value...)
		input = input / 5
	}
	// fmt.Println("  VALUE ", value)

	return value
}

func ParseInput(lines []string) ([]Value, error) {
	values := []Value{}

	for _, line := range lines {
		values = append(values, FromString(line))
	}

	return values, nil
}

func PartOneSolution(lines []string) (string, error) {
	values, err := ParseInput(lines)
	if err != nil {
		return "", err
	}

	sum := 0
	for _, value := range values {
		sum += value.ToInt()
	}

	val := FromInt(sum)

	return val.ToString(), nil
}

func PartTwoSolution(lines []string) (int, error) {
	// board, err := ParseInput(lines)
	// if err != nil {
	// 	return 0, err
	// }

	return 0, nil
}
