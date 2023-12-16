package main

import (
	// "fmt"
	"regexp"
	// "strings"
	"testing"
)

var testData = `[1,1,3,1,1]
[1,1,5,1,1]

[[1],[2,3,4]]
[[1],4]

[9]
[[8,7,6]]

[[4,4],4,4]
[[4,4],4,4,4]

[7,7,7,7]
[7,7,7]

[]
[3]

[[[]]]
[[]]

[1,[2,[3,[4,[5,6,7]]]],8,9]
[1,[2,[3,[4,[5,6,0]]]],8,9]`

var testInput = `[[[1,6,[1,9,0,9],6]]]
[[],[[[5,6,3],6,[6,5,3,3]],8,3],[],[4]]

[[1,9,2]]
[[[[],[0],[1,8,10,6]],7,2,[[]]],[6,9],[[[3],[9,7,8],4,[8,1,5],10],2],[1,[[8,10,10,4,1],9,1],8,[[5,1,2],2,0,7,[0,1,7]]]]

[[10]]
[[3,[],[7,4,8,[]],1]]`

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

func TestPartOneInput(t *testing.T) {
	lines := regexp.MustCompile("\r?\n").Split(testInput, -1)
	value, err := PartOneSolution(lines)

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	expect := 0
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
	expect := 140
	if value != expect {
		t.Errorf("Expected %d got %d", expect, value)
	}
}
