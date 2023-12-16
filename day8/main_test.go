package main

import (
	// "fmt"
	"regexp"
	// "strings"
	"testing"
)

var testData = `30373
25512
65332
33549
35390`

func TestPartOne(t *testing.T) {
	lines := regexp.MustCompile("\r?\n").Split(testData, -1)
	value, err := PartOneSolution(lines)

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if value != 21 {
		t.Errorf("Expected %d got %d", 21, value)
	}
}

func TestPartTwo(t *testing.T) {
	lines := regexp.MustCompile("\r?\n").Split(testData, -1)
	value, err := PartTwoSolution(lines)

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if value != 8 {
		t.Errorf("Expected %d got %d", 8, value)
	}
}
