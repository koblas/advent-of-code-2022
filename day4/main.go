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

	value, err = PartTwoSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2: ", value)
}

func ToRange(value string) (int, int, error) {
	parts := strings.Split(value, "-")

	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("Bad range")
	}

	a, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	b, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return int(a), int(b), nil
}

func PartOneSolution(lines []string) (int, error) {
	total := 0

	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return 0, fmt.Errorf("Bad overlap: %v ||", parts)
		}

		aStart, aEnd, err := ToRange(parts[0])
		if err != nil {
			return 0, err
		}
		bStart, bEnd, err := ToRange(parts[1])
		if err != nil {
			return 0, err
		}

		hit := false

		if aStart == bStart {
			hit = true
		} else if aStart < bStart {
			if aEnd >= bEnd {
				hit = true
			}
		} else {
			if bEnd >= aEnd {
				hit = true
			}
		}

		// fmt.Printf("%d-%d, %d-%d, %v ---- %s\n", aStart, aEnd, bStart, bEnd, hit, line)
		if hit {
			total += 1
		}
	}

	return total, nil
}

func PartTwoSolution(lines []string) (int, error) {
	total := 0

	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return 0, fmt.Errorf("Bad overlap: %v ||", parts)
		}

		aStart, aEnd, err := ToRange(parts[0])
		if err != nil {
			return 0, err
		}
		bStart, bEnd, err := ToRange(parts[1])
		if err != nil {
			return 0, err
		}

		// Make A the smaller
		if bStart < aStart {
			aStart, aEnd, bStart, bEnd = bStart, bEnd, aStart, aEnd
		}

		hit := true

		if aStart < bStart && aEnd < bStart {
			hit = false
		} else if aStart > bEnd {
			hit = false
		}

		// fmt.Printf("%d-%d, %d-%d, %v ---- %s\n", aStart, aEnd, bStart, bEnd, hit, line)
		if hit {
			total += 1
		}
	}

	return total, nil
}
