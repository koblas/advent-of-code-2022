package main

import (
	"bufio"
	"image"
	"time"

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

	start := time.Now()
	value, err := PartOneSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Part 1 (in %s): %v\n", time.Since(start), value)

	start = time.Now()
	values, err := PartTwoSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Part 2 (in %s): %v\n", time.Since(start), values)
}

type Direction rune

const (
	DIR_N = Direction('^')
	DIR_E = Direction('>')
	DIR_S = Direction('v')
	DIR_W = Direction('<')
)

var DirStep = map[Direction]image.Point{
	DIR_N: {0, -1},
	DIR_E: {1, 0},
	DIR_S: {0, 1},
	DIR_W: {-1, 0},
}

type Piece struct {
	current image.Point
	dir     Direction
	step    image.Point
}
type Board struct {
	size   image.Point
	walls  map[image.Point]image.Point
	start  image.Point
	end    image.Point
	pieces []Piece
}

func ParseInput(lines []string) (Board, error) {
	board := Board{
		walls:  map[image.Point]image.Point{},
		pieces: []Piece{},
	}
	for y, line := range lines {
		for x := 0; x < len(line); x += 1 {
			if line[x] == '#' {
				board.walls[image.Point{x, y}] = image.Point{}
			} else if line[x] == '.' {
				continue
			} else {
				board.pieces = append(board.pieces, Piece{
					current: image.Point{x, y},
					step:    DirStep[Direction(line[x])],
					dir:     Direction(line[x]),
				})
			}
		}
		board.size.X = len(line) - 1
		board.size.Y = y
	}

	for x := 0; x < board.size.X; x += 1 {
		if _, found := board.walls[image.Point{x, 0}]; found {
			board.walls[image.Point{x, 0}] = image.Point{x, board.size.Y - 1}
		} else {
			board.start = image.Point{x, 0}
		}
		if _, found := board.walls[image.Point{x, board.size.Y}]; found {
			board.walls[image.Point{x, board.size.Y}] = image.Point{x, 1}
		} else {
			board.end = image.Point{x, board.size.Y}
		}
	}
	for y := 0; y < board.size.Y; y += 1 {
		if _, found := board.walls[image.Point{0, y}]; found {
			board.walls[image.Point{0, y}] = image.Point{board.size.X - 1, y}
		}
		if _, found := board.walls[image.Point{board.size.X, y}]; found {
			board.walls[image.Point{board.size.X, y}] = image.Point{1, y}
		}
	}

	return board, nil
}

var deltas = []image.Point{
	{0, 0},
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
}

func MovePlayer(board Board, pieces []Piece, player image.Point) []image.Point {
	occupied := map[image.Point]struct{}{}
	for _, piece := range pieces {
		occupied[piece.current] = struct{}{}
	}
	canMove := func(pos image.Point) bool {
		if _, found := occupied[pos]; found {
			return false
		}
		if _, found := board.walls[pos]; found {
			return false
		}
		if pos.X < 0 || pos.Y < 0 {
			return false
		}
		if pos.X == board.size.X+1 || pos.Y == board.size.Y+1 {
			return false
		}
		// if pos.Eq(board.start) {
		// 	return false
		// }
		return true
	}

	result := []image.Point{}
	for _, delta := range deltas {
		pos := player.Add(delta)
		if canMove(pos) {
			result = append(result, pos)
		}
	}
	return result
}

func MoveSnow(board Board, pieces []Piece) []Piece {
	result := []Piece{}
	for _, piece := range pieces {
		np := piece.current.Add(piece.step)
		if np.X == 0 {
			np.X = board.size.X - 1
		}
		if np.X == board.size.X {
			np.X = 1
		}
		if np.Y == 0 {
			np.Y = board.size.Y - 1
		}
		if np.Y == board.size.Y {
			np.Y = 1
		}
		result = append(result, Piece{
			current: np,
			step:    piece.step,
			dir:     piece.dir,
		})
	}

	return result
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

func DrawBoard(board Board, pieces []Piece, player image.Point) string {
	output := ""
	for y := 0; y <= board.size.Y; y += 1 {
		for x := 0; x <= board.size.X; x += 1 {
			count := 0
			pos := image.Point{x, y}
			last := Piece{}
			for _, piece := range pieces {
				if pos.Eq(piece.current) {
					count += 1
					last = piece
				}
			}
			if count != 0 {
				if count == 1 {
					output += string(last.dir)
				} else {
					output += fmt.Sprintf("%d", count)
				}
			} else if _, found := board.walls[pos]; found {
				output += "#"
			} else if pos.Eq(player) {
				output += "E"
			} else {
				output += "."
			}
		}
		output += "\n"
	}

	return output
}

type State struct {
	depth  int
	player image.Point
	pieces []Piece
}

type Visit struct {
	depth int
	pos   image.Point
}

func Simulate(board Board, pieces []Piece, start, end image.Point) (int, []Piece) {
	queue := []State{
		{
			1,
			start,
			pieces[:],
		},
	}

	snowMemo := map[int][]Piece{
		0: pieces[:],
	}

	visited := map[Visit]struct{}{}

	for len(queue) != 0 {
		first := queue[0]
		queue = queue[1:]

		// fmt.Println("AT ", first.depth)
		// fmt.Println("== Minute ", first.depth)
		// fmt.Println(DrawBoard(board, snowMemo[first.depth-1], first.player))

		newSnow := snowMemo[first.depth]
		if newSnow == nil {
			newSnow = MoveSnow(board, snowMemo[first.depth-1])
			snowMemo[first.depth] = newSnow
		}

		options := MovePlayer(board, newSnow, first.player)
		// fmt.Println("OPT ", options)
		for _, option := range options {
			if option.Eq(end) {
				return first.depth, newSnow
			}
			v := Visit{first.depth + 1, option}
			if _, found := visited[v]; found {
				continue
			}
			visited[v] = struct{}{}
			queue = append(queue, State{
				depth:  first.depth + 1,
				player: option,
			})
		}
	}

	return -1, []Piece{}
}

func PartOneSolution(lines []string) (int, error) {
	board, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	count, _ := Simulate(board, board.pieces, board.start, board.end)

	return count, nil
}

func PartTwoSolution(lines []string) (int, error) {
	board, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	c1, state := Simulate(board, board.pieces, board.start, board.end)
	c2, state := Simulate(board, state, board.end, board.start)
	c3, _ := Simulate(board, state, board.start, board.end)
	fmt.Println("VALS ", c1, c2, c3)

	return c1 + c2 + c3, nil
}
