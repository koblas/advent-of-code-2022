package main

import (
	// "fmt"
	// "regexp"
	// "strings"
	"testing"
)

var testData = []struct {
	value   string
	partOne int
	partTwo int
}{
	{"bvwbjplbgvbhsrlpgdmjqwftvncz", 5, 0},
	{"nppdvjthqldpwncqszvftbrmjlhg", 6, 0},
	{"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 10, 0},
	{"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 11, 0},

	{"mjqjpqmgbljsphdztnvjfqwrcgsmlb", 0, 19},
	{"bvwbjplbgvbhsrlpgdmjqwftvncz", 0, 23},
	{"nppdvjthqldpwncqszvftbrmjlhg", 0, 23},
	{"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", 0, 29},
	{"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw", 0, 26},
}

func TestPartOne(t *testing.T) {
	for _, item := range testData {
		if item.partOne == 0 {
			continue
		}
		result := findFirst(item.value, 4)
		if result != item.partOne {
			t.Errorf("%s: Expected %d got %d", item.value, item.partOne, result)
		}
	}
}

func TestPartTwo(t *testing.T) {
	for _, item := range testData {
		if item.partTwo == 0 {
			continue
		}
		result := findFirst(item.value, 14)
		if result != item.partTwo {
			t.Errorf("%s: Expected %d got %d", item.value, item.partTwo, result)
		}
	}
}
