package main

import (
	// "fmt"
	"regexp"
	// "strings"
	"testing"
)

var testData = `2-4,6-8
2-3,4-5
5-7,7-9
2-8,3-7
6-6,4-6
2-6,4-8`

func TestPartOne(t *testing.T) {
	checks := []struct {
		Value  string
		Result int
	}{
		{testData, 2},
		{"1-2,1-2", 1},

		{"1-2,1-12", 1},
		{"2-2,1-12", 1},
		{"2-3,1-12", 1},
		{"1-2,2-12", 0},

		{"1-12,1-2", 1},
		{"1-12,2-2", 1},
		{"1-12,2-3", 1},
		{"2-12,1-2", 0},
	}

	for _, check := range checks {
		lines := regexp.MustCompile("\r?\n").Split(check.Value, -1)
		value, err := PartOneSolution(lines)

		if err != nil {
			t.Errorf("%s: Got error: %v", check.Value, err)
		}
		if value != check.Result {
			t.Errorf("%s: Expected %d got %d", check.Value, check.Result, value)
		}
	}
}

func TestPartTwo(t *testing.T) {
	checks := []struct {
		Value  string
		Result int
	}{
		{testData, 4},
		{"1-2,1-2", 1},

		{"1-2,1-12", 1},
		{"2-2,1-12", 1},
		{"2-3,1-12", 1},
		{"1-2,2-12", 1},

		{"1-12,1-2", 1},
		{"1-12,2-2", 1},
		{"1-12,2-3", 1},
		{"2-12,1-2", 1},
	}

	for _, check := range checks {
		lines := regexp.MustCompile("\r?\n").Split(check.Value, -1)
		value, err := PartTwoSolution(lines)

		if err != nil {
			t.Errorf("%s: Got error: %v", check.Value, err)
		}
		if value != check.Result {
			t.Errorf("%s: Expected %d got %d", check.Value, check.Result, value)
		}
	}
}
