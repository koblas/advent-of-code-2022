package main

import (
	"bufio"
	"fmt"
	"os"
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

type Point [2]int
type QueueItem struct {
	point Point
	dist  int
}
type Empty struct{}
type Visited map[Point]Empty

type Queue []QueueItem

func (queue Queue) Add(item QueueItem) Queue {
	for idx := 0; idx < len(queue); idx += 1 {
		if queue[idx].dist >= item.dist {
			nq := Queue{}
			if idx > 0 {
				nq = append(nq, queue[0:idx]...)
			}
			nq = append(nq, item)
			nq = append(nq, queue[idx:]...)

			// fmt.Println(nq)

			return nq
		}
	}

	return append(queue, item)
}

func ParseInput(lines []string) ([][]int, Point, Point, error) {
	start := [2]int{}
	end := [2]int{}
	grid := [][]int{}
	for y, line := range lines {
		row := []int{}

		for x, ch := range line {
			height := int(ch - 'a')
			if ch == 'S' {
				height = 0
				start = [2]int{x, y}
			} else if ch == 'E' {
				height = 25
				end = [2]int{x, y}
			}
			row = append(row, height)
		}

		grid = append(grid, row)
	}

	return grid, start, end, nil
}

func canMove(grid [][]int, visited Visited, current, point Point) bool {
	if point[0] < 0 || point[1] < 0 || point[0] >= len(grid[0]) || point[1] >= len(grid) {
		return false
	}
	dX := point[0] - current[0]
	dY := point[1] - current[1]
	if (dY == 1 || dY == -1) && dX != 0 {
		return false
	}
	if (dX == 1 || dX == -1) && dY != 0 {
		return false
	}
	if dY < -1 || dY > 1 || dX < -1 || dX > 1 {
		return false
	}

	if _, found := visited[point]; found {
		return false
	}
	curHeight := grid[current[1]][current[0]]
	newHeight := grid[point[1]][point[0]]

	diff := newHeight - curHeight
	return diff <= 1
}

var deltas = []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

var COLOR_RESET = "\u001b[0m"
var COLOR_RED = "\u001b[31m"
var COLOR_GREEN = "\u001b[32m"

func Dump(grid [][]int, visited Visited, point Point) {
	for y, line := range grid {
		for x, ch := range line {
			val := fmt.Sprintf("%c", ch+int('a'))
			if _, found := visited[Point{x, y}]; found {
				val = strings.ToUpper(val)
			}
			cstart, cend := "", ""
			if point[0] == x && point[1] == y {
				cstart = COLOR_RED
			}
			if canMove(grid, visited, point, Point{x, y}) {
				cstart = COLOR_GREEN
			}
			if cstart != "" {
				cend = COLOR_RESET
			}
			fmt.Printf("%s%s%s", cstart, val, cend)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func Move(grid [][]int, start Point, end Point) int {
	queue := Queue{}
	queue = queue.Add(QueueItem{start, 0})

	visited := Visited{}
	visited[start] = Empty{}

	for len(queue) != 0 {
		front := queue[0]
		queue = queue[1:]

		// Dump(grid, visited, front.point)

		if front.point[0] == end[0] && front.point[1] == end[1] {
			return front.dist
		}

		// fmt.Println("TRY ", front.point)

		for _, delta := range deltas {
			move := Point{front.point[0] + delta[0], front.point[1] + delta[1]}
			if canMove(grid, visited, front.point, move) {
				visited[move] = Empty{}

				queue = queue.Add(QueueItem{Point{move[0], move[1]}, front.dist + 1})
			}
		}
	}

	return -1
}

func PartOneSolution(lines []string) (int, error) {
	grid, start, end, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	// count := Move(grid, end, start)
	count := Move(grid, start, end)

	return count, nil
}

func PartTwoSolution(lines []string) (int, error) {
	grid, start, end, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	canidates := []Point{}
	for y, line := range grid {
		for x, ch := range line {
			if ch == 0 {
				canidates = append(canidates, Point{x, y})
			}
		}
	}

	min := -1
	for _, start = range canidates {
		val := Move(grid, start, end)

		// fmt.Println(start, val)
		if val != -1 && (min == -1 || val < min) {
			min = val
		}
	}

	return min, nil
}
