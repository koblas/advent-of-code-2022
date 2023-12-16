package main

import (
	"bufio"
	"errors"
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

type Pair struct {
	x, y int
}

type Grid struct {
	min, max Pair
	data     map[Pair]rune
}

func ParseInput(lines []string) ([][]Pair, error) {
	result := [][]Pair{}

	for _, line := range lines {
		parts := strings.Split(line, "->")
		path := []Pair{}
		for _, part := range parts {
			pair := strings.Split(strings.TrimSpace(part), ",")
			if len(pair) != 2 {
				return nil, errors.New("bad pair " + part + " " + line)
			}

			x, err := strconv.ParseInt(pair[0], 10, 0)
			if err != nil {
				return nil, err
			}
			y, err := strconv.ParseInt(pair[1], 10, 0)
			if err != nil {
				return nil, err
			}

			path = append(path, Pair{int(x), int(y)})
		}

		result = append(result, path)
	}

	return result, nil
}

func BuildGrid(paths [][]Pair) Grid {
	walls := map[Pair]rune{}

	min := paths[0][0]
	max := paths[0][0]

	min.y = 0

	for _, path := range paths {
		for _, pair := range path {
			min.x = Min(min.x, pair.x)
			max.x = Max(max.x, pair.x)
			min.y = Min(min.y, pair.y)
			max.y = Max(max.y, pair.y)
		}

		for idx := 1; idx < len(path); idx += 1 {
			walls[path[idx]] = '#'
			x := path[idx-1].x
			y := path[idx-1].y

			for x != path[idx].x || y != path[idx].y {
				walls[Pair{x, y}] = '#'

				x += Move(x, path[idx].x)
				y += Move(y, path[idx].y)
			}
		}
	}

	return Grid{min: min, max: max, data: walls}
}

func Move(start, dest int) int {
	if start < dest {
		return 1
	}
	if start > dest {
		return -1
	}
	return 0
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Format(grid Grid, path []Pair) string {
	output := []string{}
	for y := grid.min.y - 1; y < grid.max.y+2; y += 1 {
		row := ""
		for x := grid.min.x - 1; x < grid.max.x+2; x += 1 {
			if val, found := grid.data[Pair{x, y}]; found {
				row += string(val)
			} else {
				ch := "."
				for _, pos := range path {
					if pos.x == x && pos.y == y {
						ch = "~"
					}
				}
				row += ch
			}
		}
		output = append(output, row)
	}

	return strings.Join(output, "\n")
}

func DropSand(grid Grid) (bool, []Pair) {
	pos := Pair{500, 0}

	path := []Pair{pos}

	if _, hit := grid.data[pos]; hit {
		return false, path
	}

	for {
		if pos.y > grid.max.y {
			return false, path
		}

		path = append(path, pos)

		np := Pair{pos.x, pos.y + 1}
		if _, hit := grid.data[np]; !hit {
			pos = np
			continue
		}
		np.x = pos.x - 1
		if _, hit := grid.data[np]; !hit {
			pos = np
			continue
		}
		np.x = pos.x + 1
		if _, hit := grid.data[np]; !hit {
			pos = np
			continue
		}

		grid.data[pos] = 'o'

		return true, path
	}

}

func PartOneSolution(lines []string) (int, error) {
	paths, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	grid := BuildGrid(paths)

	count := 0
	for {
		more, _ := DropSand(grid)
		// fmt.Printf("%s\n\n", Format(grid, path))
		if !more {
			break
		}
		count += 1
	}

	// 	fmt.Printf("%s\n", Format(grid, []Pair{}))

	return count, nil
}

func PartTwoSolution(lines []string) (int, error) {
	values, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	grid := BuildGrid(values)

	grid = BuildGrid(append(values, []Pair{{grid.min.x - grid.max.y - 2, grid.max.y + 2}, {grid.max.x + grid.max.y + 2, grid.max.y + 2}}))

	count := 0
	for {
		more, _ := DropSand(grid)
		// fmt.Printf("%s\n\n", Format(grid, path))
		if !more {
			break
		}
		count += 1
	}

	return count, nil
}
