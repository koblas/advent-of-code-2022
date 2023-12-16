package main

import (
	// "fmt"
	"regexp"
	// "strings"
	"testing"
)

var testData = `root: pppw + sjmn
dbpl: 5
cczh: sllz + lgvd
zczc: 2
ptdq: humn - dvpt
dvpt: 3
lfqf: 4
humn: 5
ljgn: 2
sjmn: drzm * dbpl
sllz: 4
pppw: cczh / lfqf
lgvd: ljgn * ptdq
drzm: hmdt - zczc
hmdt: 32`

func TestPartOne(t *testing.T) {
	lines := regexp.MustCompile("\r?\n").Split(testData, -1)
	value, err := PartOneSolution(lines)

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	expect := 152
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
	expect := 301
	if value != expect {
		t.Errorf("Expected %d got %d", expect, value)
	}
}
