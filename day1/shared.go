package main

import (
	"bufio"
	"io"
	"sort"
	"strconv"
	"strings"
)

func process(fd io.Reader) ([]int, error) {
	elfs := []int{0}
	elfIdx := 0
	scanner := bufio.NewScanner(fd)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			elfIdx = len(elfs)
			elfs = append(elfs, 0)
			continue
		}
		value, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return elfs, err
		}
		elfs[elfIdx] += int(value)
	}

	if err := scanner.Err(); err != nil {
		return elfs, err
	}

	sort.Ints(elfs)

	return elfs, nil
}

func processOLD(fd io.Reader) ([]int, error) {
	elfs := []int{0}
	brd := bufio.NewReader(fd)
	for {
		line, err := brd.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return elfs, err
		}
		if line == "\n" {
			elfs = append(elfs, 0)
			continue
		}
		value, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		if err != nil {
			return elfs, err
		}
		elfs[len(elfs)-1] += int(value)
	}

	sort.Ints(elfs)

	return elfs, nil
}
