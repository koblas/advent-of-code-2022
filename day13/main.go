package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
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

type Value struct {
	list  []Value
	value int64
}

func ParseList(line string) (int, []Value, error) {
	if line[0] != '[' {
		return 0, nil, errors.New("Bad start character")
	}
	result := []Value{}
	idx := 1

	// fmt.Println("LINE= ", line)

	for line[idx] != ']' {
		if line[idx] == '[' {
			count, children, err := ParseList(line[idx:])
			if err != nil {
				return 0, nil, err
			}
			result = append(result, Value{list: children})
			// fmt.Println("COUNT = ", count, string(line[idx]), string(line[idx+count]))
			idx += count
		} else if line[idx] == ' ' || line[idx] == ',' {
			idx += 1
			continue
		} else {
			capture := ""
			for line[idx] != ']' && line[idx] != ',' && line[idx] != ' ' {
				capture += string(line[idx])
				idx += 1
			}
			val, err := strconv.ParseInt(capture, 10, 64)
			if err != nil {
				return 0, nil, err
			}
			// fmt.Println("CAPTURE ", val)
			result = append(result, Value{value: val})
		}

		// fmt.Printf("NEXT CHAR = %c idx=%d\n", line[idx], idx)
	}

	return idx + 1, result, nil
}

func ParseInput(lines []string) ([]*Value, error) {
	result := []*Value{}

	for _, line := range lines {
		if line == "" {
			continue
		}
		if line[0] != '[' {
			return nil, errors.New("Bad start character")
		}

		count, row, err := ParseList(line)
		if err != nil {
			return nil, err
		}
		if count != len(line) {
			return nil, fmt.Errorf("Not all characters consumed count=%d len=%d", count, len(line))
		}

		val := Value{list: row}
		result = append(result, &val)
	}

	return result, nil
}

func Format(value []Value) string {
	parts := []string{}
	for _, item := range value {
		if item.list != nil {
			parts = append(parts, Format(item.list))
		} else {
			parts = append(parts, fmt.Sprintf("%d", item.value))
		}
	}

	return "[" + strings.Join(parts, ",") + "]"
}

func CheckList(left, right []Value) (bool, error) {
	// fmt.Println("Compare ", Format(left), " vs ", Format(right))
	for idx, item := range right {
		// Left runs out of items first OK
		if len(left) == idx {
			return true, nil
		}
		done, err := Check(left[idx], item)
		if err != nil {
			return false, err
		}
		if done {
			return done, nil
		}
	}

	if len(left) != len(right) {
		return false, errors.New("right list out of items")
	}

	// No determination
	return false, nil
}

func Check(left, right Value) (bool, error) {
	if left.list == nil && right.list == nil {
		// fmt.Println("CHECK VALUE", left.value, right.value)
		if left.value == right.value {
			return false, nil
		} else if left.value > right.value {
			return false, fmt.Errorf("int compare: %d > %d", left.value, right.value)
		}

		return true, nil
	}

	ll := left.list
	rl := right.list

	if left.list == nil {
		ll = []Value{{value: left.value}}
	}
	if right.list == nil {
		rl = []Value{{value: right.value}}
	}

	return CheckList(ll, rl)
}

func PartOneSolution(lines []string) (int, error) {
	values, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	sum := 0
	for idx := 0; idx < len(values); idx += 2 {
		_, err := Check(*values[idx], *values[idx+1])
		// fmt.Println("PAIR ", idx+1, err)
		// fmt.Println("")
		if err == nil {
			sum += (idx / 2) + 1
		}
	}

	return sum, nil
}

type ByValue []*Value

func (a ByValue) Len() int      { return len(a) }
func (a ByValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByValue) Less(i, j int) bool {
	_, err := Check(*a[i], *a[j])
	return err == nil
}

func PartTwoSolution(lines []string) (int, error) {
	values, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	v2 := &Value{list: []Value{{value: 2}}}
	v6 := &Value{list: []Value{{value: 6}}}

	values = append(values, v2, v6)

	sort.Sort(ByValue(values))

	idx2 := 0
	idx6 := 0
	for idx, value := range values {
		// fmt.Println(Format(value.list))
		if value == v2 {
			idx2 = idx + 1
		}
		if value == v6 {
			idx6 = idx + 1
		}
	}

	return idx2 * idx6, nil
}
