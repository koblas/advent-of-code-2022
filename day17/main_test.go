package main

import (
	// "fmt"
	"regexp"
	// "strings"
	"testing"
)

var testData = `>>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>`

func TestPartOne(t *testing.T) {
	lines := regexp.MustCompile("\r?\n").Split(testData, -1)
	value, err := PartOneSolution(lines)

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	expect := 3068
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
	expect := 1514285714288
	if value != expect {
		t.Errorf("Expected %d got %d", expect, value)
	}
}
