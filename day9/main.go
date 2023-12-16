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

type Step struct {
	dir   rune
	count int
}

type Pos [2]int

func ParseInput(lines []string) []Step {
	result := []Step{}
	for _, line := range lines {
		parts := strings.Split(line, " ")
		count, _ := strconv.ParseInt(parts[1], 10, 64)
		result = append(result, Step{
			dir:   rune(parts[0][0]),
			count: int(count),
		})
	}

	return result
}

func doMove(pos Pos, dir rune) Pos {
	switch dir {
	case 'L':
		return Pos{pos[0] - 1, pos[1]}
	case 'R':
		return Pos{pos[0] + 1, pos[1]}
	case 'U':
		return Pos{pos[0], pos[1] - 1}
	case 'D':
		return Pos{pos[0], pos[1] + 1}
	}

	panic("Illegal move")
}

func sign(value int) int {
	if value < 0 {
		return -1
	}
	if value > 0 {
		return 1
	}
	return 0
}

func doCatch(head, tail Pos) Pos {
	dX := head[0] - tail[0]
	dY := head[1] - tail[1]

	if dX <= 1 && dX >= -1 && dY <= 1 && dY >= -1 {
		return tail
	}

	mX := sign(dX)
	mY := sign(dY)

	// fmt.Printf("CATCH: delta=%d,%d move=%d,%d\n", dX, dY, mX, mY)

	return Pos{tail[0] + mX, tail[1] + mY}
}

func Draw(head, tail Pos, size int) {
	fmt.Printf("\n------\n")
	for y := 0; y < size; y += 1 {
		for x := 0; x < size; x += 1 {
			char := '.'
			if tail[0] == x && tail[1] == y {
				char = 'T'
			}
			if head[0] == x && head[1] == y {
				char = 'H'
			}

			fmt.Printf("%c", char)
		}
		fmt.Printf("\n")
	}
}

func runRope(lines []string, rope []Pos) (int, error) {
	moves := ParseInput(lines)
	visits := map[Pos]struct{}{}

	visits[rope[0]] = struct{}{}

	// draw(head, tail, 6)
	for _, move := range moves {
		// fmt.Printf("MOVE %c %d\n", move.dir, move.count)
		for i := 0; i < move.count; i++ {
			// fmt.Println("===START===")
			rope[0] = doMove(rope[0], move.dir)

			for i := 1; i < len(rope); i += 1 {
				rope[i] = doCatch(rope[i-1], rope[i])
			}
			// draw(head, tail, 6)

			// fmt.Println("===DONE===")

			visits[rope[len(rope)-1]] = struct{}{}
		}
	}

	return len(visits), nil
}

func PartOneSolution(lines []string) (int, error) {
	rope := make([]Pos, 2)

	return runRope(lines, rope)
}

func PartTwoSolution(lines []string) (int, error) {
	rope := make([]Pos, 10)

	return runRope(lines, rope)
}
