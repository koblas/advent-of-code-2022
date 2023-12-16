package main

import (
	// "fmt"
	"regexp"
	// "strings"
	"testing"
)

var testData = `1=-0-2
12111
2=0=
21
2=01
111
20012
112
1=-1=
1-12
12
1=
122`

func TestPartOne(t *testing.T) {
	lines := regexp.MustCompile("\r?\n").Split(testData, -1)
	value, err := PartOneSolution(lines)

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	expect := "2=-1=0"
	if value != expect {
		t.Errorf("Expected %s got %s", expect, value)
	}
}

func TestPartTwo(t *testing.T) {
	lines := regexp.MustCompile("\r?\n").Split(testData, -1)
	value, err := PartTwoSolution(lines)

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	expect := 54
	if value != expect {
		t.Errorf("Expected %d got %d", expect, value)
	}
}
