package main

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type Item int
type Outcome int

const (
	Rock     = Item(1)
	Paper    = Item(2)
	Scissors = Item(3)
	Win      = Outcome(6)
	Loss     = Outcome(0)
	Draw     = Outcome(3)
)

var piece = map[string]Item{
	"A": Rock,
	"B": Paper,
	"C": Scissors,
	"X": Rock,
	"Y": Paper,
	"Z": Scissors,
}

var expected = map[string]Outcome{
	"X": Loss,
	"Y": Draw,
	"Z": Win,
}

var outcome = map[Item]map[Item]int{
	Rock: {
		Rock:     3,
		Paper:    0,
		Scissors: 6,
	},
	Paper: {
		Rock:     6,
		Paper:    3,
		Scissors: 0,
	},
	Scissors: {
		Rock:     0,
		Paper:    6,
		Scissors: 3,
	},
}

var predict = map[Item]map[Outcome]Item{
	Rock: {
		Win:  Paper,
		Loss: Scissors,
		Draw: Rock,
	},
	Paper: {
		Win:  Scissors,
		Loss: Rock,
		Draw: Paper,
	},
	Scissors: {
		Win:  Rock,
		Loss: Paper,
		Draw: Scissors,
	},
}

func process(fd io.Reader) (int, error) {
	scanner := bufio.NewScanner(fd)
	total := 0

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return total, errors.New("Not enoguht parts")
		}

		opponent, ok := piece[parts[0]]
		if !ok {
			return total, errors.New("invalid opponent")
		}
		self, ok := piece[parts[1]]
		if !ok {
			return total, errors.New("invalid self")
		}

		total += outcome[self][opponent] + int(self)
	}

	return total, scanner.Err()
}

func processPart2(fd io.Reader) (int, error) {
	scanner := bufio.NewScanner(fd)
	total := 0

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return total, errors.New("Not enoguht parts")
		}

		opponent, ok := piece[parts[0]]
		if !ok {
			return total, errors.New("invalid opponent")
		}
		result, ok := expected[parts[1]]
		if !ok {
			return total, errors.New("invalid self")
		}

		total += int(predict[opponent][result]) + int(result)
	}

	return total, scanner.Err()
}
