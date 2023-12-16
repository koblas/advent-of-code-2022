package main

import (
	"bufio"
	"io"
	"strings"
)

var scores = "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type RuneSet map[rune]struct{}

func process(fd io.Reader) (int, error) {
	scanner := bufio.NewScanner(fd)
	total := 0

	for scanner.Scan() {
		line := scanner.Text()

		part1, part2 := line[0:len(line)/2], line[len(line)/2:]

		has := RuneSet{}
		for _, c := range part2 {
			has[c] = struct{}{}
		}

		for _, c := range part1 {
			if _, found := has[c]; found {
				total += strings.IndexRune(scores, c)
				break
			}
		}
	}

	return total, scanner.Err()
}

func processPart2(fd io.Reader) (int, error) {
	scanner := bufio.NewScanner(fd)
	total := 0

	has := []RuneSet{}
	for scanner.Scan() {
		line := scanner.Text()

		values := RuneSet{}
		for _, c := range line {
			values[c] = struct{}{}
		}
		has = append(has, values)

		if len(has) != 3 {
			continue
		}

		for c := range has[0] {
			_, found1 := has[1][c]
			_, found2 := has[2][c]
			if found1 && found2 {
				total += strings.IndexRune(scores, c)
				break
			}
		}

		has = []RuneSet{}
	}

	return total, scanner.Err()
}
