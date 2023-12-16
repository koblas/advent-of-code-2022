package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

type Monkey struct {
	items []int64

	exprA  string
	opcode rune
	exprB  string

	divisor int64
	ifTrue  int64
	ifFalse int64

	inspectCount int
}

var (
	monkeyRe    = regexp.MustCompile("^Monkey (\\d+):")
	itemsRe     = regexp.MustCompile("^\\s+Starting items:(.*)")
	operationRe = regexp.MustCompile("^\\s+Operation: new = (\\S+) (\\S+) (\\S+)")
	testRe      = regexp.MustCompile("^\\s+Test: divisible by\\s+(\\d+)")
	ifTrueRe    = regexp.MustCompile("^\\s+If true: throw to monkey\\s+(\\d+)")
	ifFalseRe   = regexp.MustCompile("^\\s+If false: throw to monkey\\s+(\\d+)")
)

func ParseInput(lines []string) ([]*Monkey, error) {
	result := []*Monkey{}
	var monkey *Monkey
	for _, line := range lines {
		if parts := monkeyRe.FindStringSubmatch(line); parts != nil {
			monkey = &Monkey{
				items: []int64{},
			}
			result = append(result, monkey)
		} else if parts := itemsRe.FindStringSubmatch(line); parts != nil {
			items := strings.Split(parts[1], ",")
			if items != nil {
				for _, item := range items {
					if item == "" {
						continue
					}
					value, err := strconv.ParseInt(strings.TrimSpace(item), 10, 64)
					if err != nil {
						return nil, err
					}

					monkey.items = append(monkey.items, value)
				}
			}
		} else if parts := operationRe.FindStringSubmatch(line); parts != nil {
			monkey.exprA = parts[1]
			monkey.opcode = rune(parts[2][0])
			monkey.exprB = parts[3]
		} else if parts := testRe.FindStringSubmatch(line); parts != nil {
			value, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, err
			}
			monkey.divisor = value
		} else if parts := ifTrueRe.FindStringSubmatch(line); parts != nil {
			value, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, err
			}
			monkey.ifTrue = value
		} else if parts := ifFalseRe.FindStringSubmatch(line); parts != nil {
			value, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, err
			}
			monkey.ifFalse = value
		} else if line == "" {
			// ignore
		} else {
			return nil, fmt.Errorf("Unknown line: %s", line)
		}
	}

	return result, nil
}

func getValue(value int64, expr string) int64 {
	if expr != "old" {
		value, _ = strconv.ParseInt(expr, 10, 64)
	}
	return value
}

func UpdateMonkey(idx int, monkey *Monkey, monkeys []*Monkey, hasWorry bool, commonDivisor int64) {
	for _, item := range monkey.items {
		monkey.inspectCount += 1

		v1 := getValue(item, monkey.exprA)
		v2 := getValue(item, monkey.exprB)
		var value int64
		switch monkey.opcode {
		case '+':
			value = v1 + v2
		case '*':
			value = v1 * v2
		default:
			panic("Unknown opcode " + string(monkey.opcode))
			// case '-':
			//	value = v1 - v2
			//case '/':
			//	value = v1 / v2

		}

		if hasWorry {
			value = value / 3
		}

		throw := monkey.ifFalse
		value = value % commonDivisor
		if value%monkey.divisor == 0 {
			throw = monkey.ifTrue
		}
		monkeys[throw].items = append(monkeys[throw].items, value)
	}
	monkey.items = []int64{}
}

func Calculate(lines []string, rounds int, hasWorry bool) (int, error) {
	monkeys, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	commonDivisor := int64(1)
	for _, monkey := range monkeys {
		commonDivisor *= monkey.divisor
	}

	for i := 0; i < rounds; i += 1 {
		for idx, monkey := range monkeys {
			UpdateMonkey(idx, monkey, monkeys, hasWorry, commonDivisor)
		}
	}

	for idx, monkey := range monkeys {
		fmt.Printf("Monkey %d (%d): %v\n", idx, monkey.inspectCount, monkey.items)
	}

	counts := []int{}
	for _, monkey := range monkeys {
		counts = append(counts, monkey.inspectCount)
	}
	sort.Ints(counts)

	return counts[len(counts)-2] * counts[len(counts)-1], nil
}

func PartOneSolution(lines []string) (int, error) {
	return Calculate(lines, 20, true)
}

func PartTwoSolution(lines []string) (int, error) {
	return Calculate(lines, 10_000, false)
}
