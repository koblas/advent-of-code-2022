package main

import (
	// "fmt"
	"regexp"
	// "strings"
	"testing"
)

var testData = `498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`

func TestPartOne(t *testing.T) {
	lines := regexp.MustCompile("\r?\n").Split(testData, -1)
	value, err := PartOneSolution(lines)

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	expect := 24
	if value != expect {
		t.Errorf("Expected %d got %d", expect, value)
	}
}

func TestPartTwo(t *testing.T) {
	lines := regexp.MustCompile("\r?\n").Split(testData, -1)
	value, err := PartTwoSolution(lines)

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	expect := 93
	if value != expect {
		t.Errorf("Expected %d got %d", expect, value)
	}
}
