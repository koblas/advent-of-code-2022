package main

import (
	"bufio"
	"image"

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

type Direction int

const (
	DIR_NONE = Direction(-1)
	DIR_N    = Direction(0)
	DIR_NE   = Direction(1)
	DIR_E    = Direction(2)
	DIR_SE   = Direction(3)
	DIR_S    = Direction(4)
	DIR_SW   = Direction(5)
	DIR_W    = Direction(6)
	DIR_NW   = Direction(7)
)

var DirName = map[Direction]string{
	DIR_NONE: "NONE",
	DIR_N:    "N",
	DIR_NE:   "NE",
	DIR_E:    "E",
	DIR_SE:   "SE",
	DIR_S:    "S",
	DIR_SW:   "SW",
	DIR_W:    "W",
	DIR_NW:   "NW",
}

type Elf struct {
	current image.Point
}

var DirStep = map[Direction]image.Point{
	DIR_NW: {-1, -1},
	DIR_N:  {0, -1},
	DIR_NE: {1, -1},
	DIR_E:  {1, 0},
	DIR_SE: {1, 1},
	DIR_S:  {0, 1},
	DIR_SW: {-1, 1},
	DIR_W:  {-1, 0},
}

var InverseDirs = map[Direction][]Direction{
	DIR_N:  {DIR_N},
	DIR_E:  {DIR_E},
	DIR_S:  {DIR_S},
	DIR_W:  {DIR_W},
	DIR_NE: {DIR_N, DIR_E},
	DIR_NW: {DIR_N, DIR_W},
	DIR_SE: {DIR_S, DIR_E},
	DIR_SW: {DIR_S, DIR_W},
}

func ParseInput(lines []string) ([]*Elf, error) {
	elves := []*Elf{}
	for y, line := range lines {
		for x := 0; x < len(line); x += 1 {
			if line[x] == '#' {
				elf := Elf{
					current: image.Point{x + 1, y + 1},
				}
				elves = append(elves, &elf)
			}
		}
	}

	return elves, nil
}

func GetAllowedMoves(elf *Elf, elves []*Elf) map[Direction]struct{} {
	inverse := map[image.Point]Direction{}
	for dir, offset := range DirStep {
		inverse[elf.current.Add(offset)] = dir
	}

	possible := map[Direction]struct{}{
		DIR_N: {},
		DIR_S: {},
		DIR_E: {},
		DIR_W: {},
	}
	for _, checkElf := range elves {
		dir, found := inverse[checkElf.current]
		if !found {
			continue
		}
		// fmt.Println("FOUND ELF  ", elf.id, " => ", checkElf.id)
		for _, value := range InverseDirs[dir] {
			// fmt.Println("REMOVING ", dir, value)
			delete(possible, value)
		}
	}

	if len(possible) == 4 || len(possible) == 0 {
		return map[Direction]struct{}{}
	}

	return possible
}

var ORDER = [4]Direction{DIR_N, DIR_S, DIR_W, DIR_E}

func Move(elves []*Elf, step int) bool {
	newpos := map[image.Point][]*Elf{}
	for _, elf := range elves {
		possible := GetAllowedMoves(elf, elves)
		// fmt.Println("ALLOWED", elf.id, possible)
		if len(possible) == 0 {
			if val, found := newpos[elf.current]; found {
				newpos[elf.current] = append(val, elf)
			} else {
				newpos[elf.current] = []*Elf{elf}
			}
			continue
		}

		for orderIdx := 0; orderIdx < len(ORDER); orderIdx += 1 {
			dir := ORDER[(step+orderIdx)%len(ORDER)]
			if _, found := possible[dir]; !found {
				continue
			}

			np := elf.current.Add(DirStep[dir])
			if val, found := newpos[np]; found {
				newpos[np] = append(val, elf)
			} else {
				newpos[np] = []*Elf{elf}
			}

			break
		}
	}

	// Now do moves
	moved := false
	for pos, elfList := range newpos {
		if len(elfList) == 1 {
			elf := elfList[0]
			if !elf.current.Eq(pos) {
				elf.current = pos
				moved = true
			}
		}
	}

	return moved
}

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func MinPoint(a, b image.Point) image.Point {
	return image.Point{Min(a.X, b.X), Min(a.Y, b.Y)}
}

func MaxPoint(a, b image.Point) image.Point {
	return image.Point{Max(a.X, b.X), Max(a.Y, b.Y)}
}

func DrawBoard(elves []*Elf) string {
	min := elves[0].current
	max := elves[0].current
	for _, elf := range elves {
		min = MinPoint(elf.current, min)
		max = MaxPoint(elf.current, max)
	}

	occupied := map[image.Point]*Elf{}
	for _, elf := range elves {
		occupied[elf.current] = elf
	}

	max = max.Add(image.Point{1, 1})

	output := ""
	for y := min.Y - 1; y < max.Y+1; y += 1 {
		for x := min.X - 1; x < max.X+1; x += 1 {
			if _, found := occupied[image.Point{x, y}]; found {
				output = output + "#"
			} else {
				output = output + "."
			}
		}
		output = output + "\n"
	}

	return output
}

func Score(elves []*Elf) int {
	min := elves[0].current
	max := elves[0].current
	for _, elf := range elves {
		min = MinPoint(elf.current, min)
		max = MaxPoint(elf.current, max)
	}
	max = max.Add(image.Point{1, 1})

	return (max.X-min.X)*(max.Y-min.Y) - len(elves)
}

func PartOneSolution(lines []string) (int, error) {
	input, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	for i := 0; i < 10; i += 1 {
		Move(input, i)
	}

	return Score(input), nil
}

func PartTwoSolution(lines []string) (int, error) {
	input, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	for i := 0; true; i += 1 {
		if !Move(input, i) {
			return i + 1, nil
		}
	}

	return 0, nil
}
