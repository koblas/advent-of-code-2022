package main

import (
	// "fmt"
	"regexp"
	// "strings"
	"testing"
)

var testData = `R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2`

var testDataTwo = `R 5
U 8
L 8
D 3
R 17
D 10
L 25
U 20`

func TestPartOne(t *testing.T) {
	lines := regexp.MustCompile("\r?\n").Split(testData, -1)
	value, err := PartOneSolution(lines)

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	expect := 13
	if value != expect {
		t.Errorf("Expected %d got %d", expect, value)
	}
}

func TestPartTwo(t *testing.T) {
	lines := regexp.MustCompile("\r?\n").Split(testDataTwo, -1)
	value, err := PartTwoSolution(lines)

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	expect := 36
	if value != expect {
		t.Errorf("Expected %d got %d", expect, value)
	}
}
