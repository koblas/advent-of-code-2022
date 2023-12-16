package main

import (
	"bufio"
	"errors"
	"image"
	"regexp"
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

type Direction int

const (
	RIGHT = Direction(0)
	DOWN  = Direction(1)
	LEFT  = Direction(2)
	UP    = Direction(3)
)

var DirName = map[Direction]string{
	RIGHT: "right",
	DOWN:  "down",
	LEFT:  "left",
	UP:    "up",
}

type Faces []Board
type Offsets []image.Point
type Board map[image.Point]rune
type Step map[Direction]image.Point

type FallOff struct {
	direction Direction
	frameOne  int
	frameTwo  int
	reverse   bool
	opcode    string
}

type Ajacency = map[int]map[Direction]FallOff

var AjTest = Ajacency{
	0: {
		RIGHT: FallOff{LEFT, 0, 5, false, "none"},
		DOWN:  FallOff{DOWN, 3, 3, false, "none"},
		LEFT:  FallOff{DOWN, 0, 2, false, "none"},
		UP:    FallOff{DOWN, 4, 1, false, "none"},
	},
	1: {
		RIGHT: FallOff{RIGHT, 2, 2, false, "none"},
		DOWN:  FallOff{UP, 1, 4, false, "none"},
		LEFT:  FallOff{UP, 3, 5, false, "none"},
		UP:    FallOff{DOWN, 1, 0, false, "none"},
	},
	2: {
		RIGHT: FallOff{RIGHT, 3, 3, false, "none"},
		DOWN:  FallOff{RIGHT, 2, 4, false, "none"},
		LEFT:  FallOff{LEFT, 3, 1, false, "none"},
		UP:    FallOff{RIGHT, 2, 0, false, "none"},
	},
	3: {
		RIGHT: FallOff{DOWN, 1, 5, false, "none"},
		DOWN:  FallOff{DOWN, 4, 4, false, "none"},
		LEFT:  FallOff{LEFT, 2, 2, false, "none"},
		UP:    FallOff{UP, 0, 0, false, "none"},
	},
	4: {
		RIGHT: FallOff{RIGHT, 5, 5, false, "none"},
		DOWN:  FallOff{UP, 0, 1, false, "none"},
		LEFT:  FallOff{UP, 5, 2, false, "none"},
		UP:    FallOff{UP, 3, 3, false, "none"},
	},
	5: {
		RIGHT: FallOff{LEFT, 4, 0, false, "none"},
		DOWN:  FallOff{RIGHT, 5, 1, false, "none"},
		LEFT:  FallOff{LEFT, 4, 4, false, "none"},
		UP:    FallOff{LEFT, 5, 3, false, "none"},
	},
}

var AjBig = Ajacency{
	0: {
		RIGHT: FallOff{RIGHT, 1, 1, true, "none"},
		DOWN:  FallOff{DOWN, 2, 2, true, "none"},
		LEFT:  FallOff{RIGHT, 1, 3, false, "flip"},
		UP:    FallOff{RIGHT, 4, 5, false, "L=T"},
	},
	1: {
		RIGHT: FallOff{LEFT, 0, 4, true, "flip"},
		DOWN:  FallOff{LEFT, 1, 2, false, "L=T"},
		LEFT:  FallOff{LEFT, 0, 0, false, "none"},
		UP:    FallOff{UP, 1, 5, false, "none"},
	},
	2: {
		RIGHT: FallOff{UP, 2, 1, true, "B=R"},
		DOWN:  FallOff{DOWN, 4, 4, false, "none"},
		LEFT:  FallOff{DOWN, 2, 3, true, "B=R"},
		UP:    FallOff{UP, 0, 0, false, "none"},
	},
	3: {
		RIGHT: FallOff{RIGHT, 4, 4, false, "none"},
		DOWN:  FallOff{DOWN, 5, 5, false, "none"},
		LEFT:  FallOff{RIGHT, 4, 0, true, "flip"},
		UP:    FallOff{RIGHT, 5, 2, false, "L=T"},
	},
	4: {
		RIGHT: FallOff{LEFT, 3, 1, false, "flip"},
		DOWN:  FallOff{LEFT, 0, 5, false, "L=T"},
		LEFT:  FallOff{LEFT, 3, 3, false, "none"},
		UP:    FallOff{UP, 2, 2, false, "none"},
	},
	5: {
		RIGHT: FallOff{UP, 5, 4, false, "B=R"},
		DOWN:  FallOff{DOWN, 3, 1, false, "none"},
		LEFT:  FallOff{DOWN, 5, 0, true, "B=R"},
		UP:    FallOff{UP, 3, 3, true, "none"},
	},
}

var dirstep = Step{
	RIGHT: image.Point{1, 0},
	LEFT:  image.Point{-1, 0},
	DOWN:  image.Point{0, 1},
	UP:    image.Point{0, -1},
}

func ParseInput(lines []string) (Board, image.Point, []string, error) {
	board := Board{}
	size := image.Point{}
	moves := []string{}
	for x := 0; x < 1000; x += 1 {
		board[image.Point{x, -1}] = ' '
	}
	for y, line := range lines {
		if strings.Contains(line, ".") || strings.Contains(line, "#") {
			if size.X < len(line) {
				size.X = len(line)
			}
			size.Y = y + 1
			for x := 0; x < len(line); x += 1 {
				board[image.Point{x, y}] = rune(line[x])
			}
			board[image.Point{-1, y}] = ' '
			board[image.Point{len(line), y}] = ' '
		} else if line == "" {
			for x := 0; x < 1000; x += 1 {
				board[image.Point{x, y}] = ' '
			}
		} else if line != "" {
			moves = moveRe.FindAllString(line, -1)
		}
	}

	return board, size, moves, nil
}

var moveRe = regexp.MustCompile(`([LR]|\d+)`)

func Fold(board Board, size image.Point) ([]Board, Offsets, int) {
	width := size.X
	height := size.Y
	cell := 0

	// fmt.Println("SIZE ", size)

	if width%3 == 0 && height%4 == 0 {
		cell = width / 3
		// input case
	} else if width%4 == 0 && height%3 == 0 {
		// test case
		cell = width / 4
	}
	width = width / cell
	height = height / cell
	// fmt.Println("SIZE SCALED", width, height)

	faces := []Board{}
	offsets := Offsets{}
	for y := 0; y < height; y += 1 {
		for x := 0; x < width; x += 1 {
			base := image.Point{x * cell, y * cell}
			if val, found := board[base]; found && val != ' ' {
				// fmt.Print("@")

				face := Board{}
				for xx := 0; xx < cell; xx += 1 {
					for yy := 0; yy < cell; yy += 1 {
						p := image.Point{xx, yy}
						face[p] = board[base.Add(p)]
					}
				}
				faces = append(faces, face)
				offsets = append(offsets, image.Point{1, 1}.Add(base))
			} else {
				// fmt.Print(" ")
			}
		}
		// fmt.Println()
	}

	return faces, offsets, cell
}

func Move(current image.Point, direction Direction, is3d bool) image.Point {
	delta := dirstep[direction]
	np := current.Add(delta)

	return np
}

func Turn(direction Direction, step string) Direction {
	if step == "L" {
		if direction == RIGHT {
			return UP
		} else {
			return direction - 1
		}
	}
	// Must be R
	if direction == UP {
		return RIGHT
	} else {
		return direction + 1
	}
}

func GetStart(board Board, size image.Point) (image.Point, error) {
	current := image.Point{0, 0}
	for {
		val, found := board[current]
		if found && val == '.' {
			return current, nil
		}
		current.X += 1
		if current.X == size.X {
			current.X = 0
			current.Y += 1
			if current.Y == size.Y {
				return image.Point{}, errors.New("fell off the board")
			}
		}
	}
}

func Walk(faces []Board, cellSize int, ajacent Ajacency, moves []string, isPartTwo bool) (image.Point, Direction, int, error) {
	faceIdx := 0
	direction := RIGHT
	board := faces[faceIdx]
	current, err := GetStart(board, image.Point{cellSize, cellSize})
	if err != nil {
		return current, direction, faceIdx, err
	}

	// fmt.Println(moveRe.FindAllString(moves, -1))
	for _, step := range moves {
		if step == "L" || step == "R" {
			direction = Turn(direction, step)
			continue
		} else {
			dist, err := strconv.Atoi(step)
			// fmt.Println("MOVING ", dist, " DIR=", dirName[direction])
			if err != nil {
				return current, RIGHT, faceIdx, err
			}

			fmt.Println("==STEP", dist, current)
			for ; dist > 0; dist -= 1 {
				fmt.Println("AT ", faceIdx, current, DirName[direction])
				np := current.Add(dirstep[direction])
				nxtFace := faceIdx
				nxtBoard := board
				nxtDirection := direction

				ch, found := board[np]
				if !found {
					aj := ajacent[faceIdx][direction]
					// fmt.Println("OLD POS", np, DirName[direction], dirstep[direction])
					np = np.Sub(dirstep[direction].Mul(cellSize))
					// fmt.Println(" FIXED POS", np)

					nxtDirection = direction
					nxtFace = aj.frameOne
					if isPartTwo {
						nxtDirection = aj.direction
						nxtFace = aj.frameTwo

						shared := current.X
						if direction == UP || direction == DOWN {
							shared = current.Y
						}
						if aj.reverse {
							shared = cellSize - 1 - shared
						}

						switch nxtDirection {
						case UP:
							np = image.Point{cellSize - 1, shared}
						case RIGHT:
							np = image.Point{shared, 0}
						case DOWN:
							np = image.Point{0, shared}
						case LEFT:
							np = image.Point{shared, cellSize - 1}
						}

						fmt.Println("__ CHANGE", faceIdx, "=>", nxtFace, DirName[direction], " => ", DirName[nxtDirection], current, np)
					}

					// fmt.Println("ADJ ", faceIdx, dirName[direction], " AJ=", aj.board, dirName[nxtDirection])

					nxtBoard = faces[nxtFace]
					// for y := 0; y < cellSize; y += 1 {
					// 	for x := 0; x < cellSize; x += 1 {
					// 		fmt.Print(string(nxtBoard[image.Point{x, y}]))
					// 	}
					// 	fmt.Println()
					// }

					ch, found = nxtBoard[np]
					if !found {
						panic(fmt.Sprintf("MOVED OFF BOARD nxtFace=%v  np=%v", nxtFace, np))
					}
				}
				// fmt.Println("CHECK CHAR frame=", faceIdx, " VAL=", string(ch), " POS=", np, " DIR=", dirName[nxtDirection])
				if ch == '#' {
					break
				}
				// Found cases
				if ch == '.' {
					faceIdx = nxtFace
					board = nxtBoard
					current = np
					direction = nxtDirection
				}
			}
		}
	}

	return current, direction, faceIdx, nil
}

func Score(pos image.Point, dir Direction) int {
	// fmt.Println("FINAL ", pos)
	return 1000*(pos.Y) + 4*(pos.X) + int(dir)
}

func checkTable(table Ajacency) {
	for frame := 0; frame < len(table); frame += 1 {
		for _, item := range table[frame] {
			f2 := item.frameTwo

			backDir := Turn(Turn(item.direction, "L"), "L")

			item2 := table[f2][backDir]
			if item2.frameTwo != frame {
				panic(fmt.Sprintf("MISMATCH BETWEEN %d(%s) and %d(%s)", frame, DirName[item.direction], f2, DirName[backDir]))
			}
		}
	}
}

func init() {
	checkTable(AjTest)
	checkTable(AjBig)
}

func PartOneSolution(lines []string) (int, error) {
	board, size, moves, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	faces, offsets, faceSize := Fold(board, size)

	ajacency := AjTest
	if faceSize == 50 {
		ajacency = AjBig
	}

	pos, dir, bidx, err := Walk(faces, faceSize, ajacency, moves, false)
	if err != nil {
		return 0, err
	}

	// fmt.Println("POS ", pos, "FACE=", bidx)
	// fmt.Println("FOFSETS", offsets)

	return Score(pos.Add(offsets[bidx]), dir), nil
}

func PartTwoSolution(lines []string) (int, error) {
	board, size, moves, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	faces, offsets, faceSize := Fold(board, size)

	ajacency := AjTest
	if faceSize == 50 {
		ajacency = AjBig
	}

	pos, dir, bidx, err := Walk(faces, faceSize, ajacency, moves, true)
	if err != nil {
		return 0, err
	}

	return Score(pos.Add(offsets[bidx]), dir), nil
}
