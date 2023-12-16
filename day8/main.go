package main

import (
	"bufio"
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

	value, err = PartTwoSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2: ", value)
}

func ParseInput(lines []string) [][]int {
	result := [][]int{}
	for _, line := range lines {
		data := make([]int, len(line))
		result = append(result, data)
		for idx, ch := range line {
			data[idx] = int(ch) - int('0')
		}
	}

	return result
}

func CheckVisible(trees [][]int, x, y int) bool {
	height := trees[y][x]

	obscured := false
	for px := x - 1; px >= 0; px -= 1 {
		if trees[y][px] >= height {
			obscured = true
			break
		}
	}
	if !obscured {
		return true
	}

	obscured = false
	for px := x + 1; px < len(trees[y]); px += 1 {
		if trees[y][px] >= height {
			obscured = true
			break
		}
	}
	if !obscured {
		return true
	}

	obscured = false
	for py := y - 1; py >= 0; py -= 1 {
		if trees[py][x] >= height {
			obscured = true
			break
		}
	}
	if !obscured {
		return true
	}

	obscured = false
	for py := y + 1; py < len(trees); py += 1 {
		if trees[py][x] >= height {
			obscured = true
			break
		}
	}
	if !obscured {
		return true
	}

	return false
}

func GetScore(trees [][]int, x, y int) int {
	height := trees[y][x]

	count := 0
	for px := x - 1; px >= 0; px -= 1 {
		count += 1
		if trees[y][px] >= height {
			break
		}
	}
	score := count

	count = 0
	for px := x + 1; px < len(trees[y]); px += 1 {
		count += 1
		if trees[y][px] >= height {
			break
		}
	}
	score *= count

	count = 0
	for py := y - 1; py >= 0; py -= 1 {
		count += 1
		if trees[py][x] >= height {
			break
		}
	}
	score *= count

	count = 0
	for py := y + 1; py < len(trees); py += 1 {
		count += 1
		if trees[py][x] >= height {
			break
		}
	}
	score *= count

	return score
}

func PartOneSolution(lines []string) (int, error) {
	trees := ParseInput(lines)

	result := 0
	for y := 0; y < len(lines); y += 1 {
		for x := 0; x < len(lines[y]); x += 1 {
			if CheckVisible(trees, x, y) {
				// fmt.Printf("%d! ", trees[y][x])
				result += 1
			} else {
				// fmt.Printf("%d  ", trees[y][x])
			}
		}
		// fmt.Printf("\n")
	}

	return result, nil
}

func PartTwoSolution(lines []string) (int, error) {
	trees := ParseInput(lines)
	scoreMax := 1

	for y := 0; y < len(lines); y += 1 {
		for x := 0; x < len(lines[y]); x += 1 {
			score := GetScore(trees, x, y)
			if score > scoreMax {
				scoreMax = score
			}
			// fmt.Printf("%d(%4d)  ", trees[y][x], score)
		}
		// fmt.Printf("\n")
	}
	return scoreMax, nil
}
