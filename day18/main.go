package main

import (
	"bufio"
	"strconv"
	"strings"

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

type Point struct {
	X, Y, Z int
}
type Grid map[Point]int

func (p Point) Add(dir Point) Point {
	return Point{p.X + dir.X, p.Y + dir.Y, p.Z + dir.Z}
}

var Directions = []Point{
	{0, 0, 1},
	{0, 0, -1},
	{0, 1, 0},
	{0, -1, 0},
	{1, 0, 0},
	{-1, 0, 0},
}

func ParseInput(lines []string) ([]Point, error) {
	result := []Point{}
	for _, line := range lines {
		parts := strings.Split(line, ",")
		x, _ := strconv.ParseInt(parts[0], 10, 0)
		y, _ := strconv.ParseInt(parts[1], 10, 0)
		z, _ := strconv.ParseInt(parts[2], 10, 0)

		result = append(result, Point{X: int(x), Y: int(y), Z: int(z)})
	}

	return result, nil
}

func BuildGrid(points []Point, value int) Grid {
	grid := Grid{}

	for _, point := range points {
		grid[point] = value
	}

	return grid
}

func Basic(points []Point, grid Grid) (int, error) {
	count := 0

	for _, point := range points {
		for _, dir := range Directions {
			if _, found := grid[point.Add(dir)]; !found {
				count += 1
			}
		}
	}

	return count, nil
}

func FlowShape(index int, point Point, grid Grid) Grid {
	grid[point] = index
	seen := Grid{}
	seen[point] = 0

	neighbors := []Point{point}
	for len(neighbors) != 0 {
		check := neighbors[0]
		neighbors = neighbors[1:]

		for _, dir := range Directions {
			np := check.Add(dir)
			if _, found := grid[np]; !found {
				continue
			}
			if _, found := seen[np]; found {
				continue
			}
			seen[np] = 0
			grid[np] = index
			neighbors = append(neighbors, np)
		}
	}

	return seen
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

func TouchAir(point Point, airGrid, grid Grid) {
	grid[point] = -1
	for _, dir := range Directions {
		np := point.Add(dir)

		if _, found := airGrid[np]; !found {
			continue
		}
		if _, found := grid[np]; found {
			continue
		}
		TouchAir(np, airGrid, grid)
	}
}

func GetMinMax(points []Point) (Point, Point) {
	min := points[0]
	max := points[0]

	for _, point := range points {
		min.X = Min(point.X, min.X)
		min.Y = Min(point.Y, min.Y)
		min.Z = Min(point.Z, min.Z)
		max.X = Max(point.X, max.X)
		max.Y = Max(point.Y, max.Y)
		max.Z = Max(point.Z, max.Z)
	}

	return min, max
}

func BuildAir(points []Point, grid Grid) {
	min, max := GetMinMax(points)

	min = min.Add(Point{-1, -1, -1})
	max = max.Add(Point{1, 1, 1})

	airPoints := []Point{}
	for x := min.X; x <= max.X; x += 1 {
		for y := min.Y; y <= max.Y; y += 1 {
			for z := min.Z; z <= max.Z; z += 1 {
				np := Point{x, y, z}
				if _, found := grid[np]; !found {
					airPoints = append(airPoints, np)
				}
			}
		}
	}

	airGrid := BuildGrid(airPoints, 0)

	TouchAir(min, airGrid, grid)
}

func Hard(points []Point, grid Grid) (int, error) {
	BuildAir(points, grid)

	count := 0

	for _, point := range points {
		for _, dir := range Directions {
			if val, found := grid[point.Add(dir)]; found && val == -1 {
				count += 1
			}
		}
	}

	return count, nil
}

func PartOneSolution(lines []string) (int, error) {
	points, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}
	grid := BuildGrid(points, 1)
	sides, err := Basic(points, grid)

	return sides, err
}

func PartTwoSolution(lines []string) (int, error) {
	points, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}
	grid := BuildGrid(points, 1)
	sides, err := Hard(points, grid)

	return sides, err
}
