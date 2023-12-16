package main

import (
	"bufio"
	// "errors"
	"fmt"
	"image"
	"os"
)

type Rock []image.Point
type Grid struct {
	height       int
	points       map[image.Point]struct{}
	values       []int
	heightToRock []int
}

var rocks = []Rock{
	{{0, 0}, {1, 0}, {2, 0}, {3, 0}},         // ####
	{{1, 0}, {0, 1}, {1, 1}, {2, 1}, {1, 2}}, // +
	{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}}, // inverse L
	{{0, 0}, {0, 1}, {0, 2}, {0, 3}},         // |
	{{0, 0}, {1, 0}, {0, 1}, {1, 1}},         // box
}

const MAX_STEP = 2022

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

func ParseInput(lines []string) ([]int, error) {
	result := []int{}
	for _, line := range lines {
		for _, ch := range line {
			dir := 1
			if ch == '<' {
				dir = -1
			}
			result = append(result, dir)
		}
	}

	return result, nil
}

func (grid *Grid) DrawGrid(rock *Rock, pos *image.Point) {
	height := grid.Height() + 4
	if rock != nil {
		height += 4
	}

	for y := height; y >= 0; y -= 1 {
		fmt.Print("|")
		for x := 0; x < 7; x += 1 {
			np := image.Point{x, y}
			ch := "."
			if grid.IsOccupied(np) {
				ch = "#"
			}
			if rock != nil {
				for _, val := range *rock {
					if val.Add(*pos).Eq(np) {
						ch = "@"
					}
				}
			}
			fmt.Print(ch)
		}
		fmt.Println("|")
	}
	fmt.Println("")
}

func (grid *Grid) FindRepeat() (int, int) {
	var match []int
	if len(grid.values) > 20_000 {
		match = make([]int, 20_000)
	} else {
		match = make([]int, 20)
	}

	for y := 0; y < grid.height-len(match); y += 1 {
		for idx := 0; idx < len(match); idx += 1 {
			match[idx] = grid.values[y+idx]
		}

		found := false
		var yy int
		for yy = y + len(match); yy < grid.height-len(match) && !found; yy += 1 {
			found = true
			for idx := 0; idx < len(match) && idx+yy < grid.height && found; idx += 1 {
				found = found && (match[idx] == grid.values[idx+yy])
			}
		}

		if found {
			pStart := y
			pLen := yy - y - 1

			fmt.Println("FOUND REPEAT ", y, yy-y)
			fmt.Print("PATTTERN ")
			for idx := pStart; idx < pLen+pStart; idx += 1 {
				fmt.Print(grid.values[idx], " ")
			}
			fmt.Println()

			return pStart, pLen
		}
	}

	fmt.Println("NO REPEAT")
	return -1, 0
}

func (grid *Grid) Set(points []image.Point, rockIdx int) {
	minY := points[0].Y
	maxY := points[0].Y
	for _, point := range points {
		if grid.height < point.Y {
			grid.height = point.Y
			grid.values = append(grid.values, 0)
			grid.heightToRock = append(grid.heightToRock, rockIdx)
		}
		if maxY < point.Y {
			maxY = point.Y
		}
		grid.points[point] = struct{}{}
	}

	for y := minY; y <= maxY; y += 1 {
		sum := 0
		for x := 0; x < 7; x += 1 {
			sum = sum * 2
			if _, found := grid.points[image.Point{x, y}]; found {
				sum += 1
			}
		}
		grid.values[y] = sum
	}
}

func (grid *Grid) IsOccupied(pos image.Point) bool {
	_, found := grid.points[pos]

	return found
}

func (grid *Grid) Height() int {
	return grid.height
}

func (rock *Rock) CanMoveDown(grid *Grid, pos image.Point) bool {
	if pos.Y == 1 {
		return false
	}
	delta := image.Point{0, -1}.Add(pos)
	for _, point := range *rock {
		if grid.IsOccupied(point.Add(delta)) {
			return false
		}
	}
	return true
}

func (rock *Rock) CanBlow(grid *Grid, pos image.Point, dir int) bool {
	delta := image.Point{dir, 0}.Add(pos)
	for _, point := range *rock {
		np := point.Add(delta)
		if np.X < 0 || np.X > 6 {
			return false
		}
		if grid.IsOccupied(np) {
			return false
		}
	}
	return true
}

func (rock *Rock) Place(grid *Grid, pos image.Point, rockIdx int) {
	vals := []image.Point{}
	for _, point := range *rock {
		vals = append(vals, point.Add(pos))
	}
	grid.Set(vals, rockIdx)
}

func Simulate(lines []string, maxSteps int) (int, *Grid, error) {
	input, err := ParseInput(lines)
	if err != nil {
		return 0, nil, err
	}

	grid := &Grid{
		points:       map[image.Point]struct{}{},
		values:       []int{0},
		heightToRock: []int{0},
		height:       0,
	}

	isMoveStep := true
	blowIdx := 0
	var rock *Rock
	var pos image.Point
	rockCount := 0

	for {
		if rock == nil {
			rock = &rocks[rockCount%len(rocks)]
			rockCount += 1
			if rockCount == maxSteps+1 {
				break
			}
			isMoveStep = false

			pos = image.Point{2, grid.Height() + 4}
			// fmt.Println("PLACE NEW ROCK")
		} else if isMoveStep {
			if rock.CanMoveDown(grid, pos) {
				// fmt.Println("MOVE DOWN")
				pos.Y -= 1
			} else {
				// fmt.Println("PLACE")
				rock.Place(grid, pos, rockCount+1)
				rock = nil
			}
			isMoveStep = false
		} else {
			dir := input[blowIdx%len(input)]
			blowIdx += 1
			if rock.CanBlow(grid, pos, dir) {
				// fmt.Println("MOVE LEFT/RIGHT")
				pos.X += dir
			} else {
				// fmt.Println("ON EDGE")
			}
			isMoveStep = true
		}
		// DrawGrid(grid, rock, &pos)
	}

	// grid.DrawGrid(nil, nil)
	fmt.Println(grid.values)

	return grid.Height(), grid, nil
}

func PartOneSolution(lines []string) (int, error) {
	height, grid, err := Simulate(lines, 2022)

	fmt.Println("GRID ", grid.values)
	fmt.Println("STARTS ", grid.heightToRock)

	grid.FindRepeat()

	return height, err
}

func PartTwoSolution(lines []string) (int, error) {
	_, grid, err := Simulate(lines, 100_000)

	pStart, pLen := grid.FindRepeat()

	rockStart := grid.heightToRock[pStart]
	rocks := grid.heightToRock[pStart+pLen] - rockStart

	copies := (1_000_000_000_000 - rockStart) / rocks
	remain := (1_000_000_000_000 - rockStart) % rocks
	fmt.Println("HEIGHTS ", grid.heightToRock[pStart:pStart+pLen])
	fmt.Println("ROCKS ", rockStart, rocks, "COPIES ", copies, remain)

	height := 0
	for height = 0; height < pLen; height += 1 {
		rockIdx := grid.heightToRock[pStart+height] - rockStart
		if rockIdx > remain {
			break
		}
	}

	newHeight := pStart + copies*pLen + height - 1
	fmt.Println("DONE ", pStart, pLen, rockStart, rocks, copies, newHeight)

	// height, _, err := Simulate(lines, 1_000_000_000_000)
	return newHeight, err
}
