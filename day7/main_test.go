package main

import (
	// "fmt"
	"regexp"
	// "strings"
	"testing"
)

var testData = `$ cd /
$ ls
dir a
14848514 b.txt
8504156 c.dat
dir d
$ cd a
$ ls
dir e
29116 f
2557 g
62596 h.lst
$ cd e
$ ls
584 i
$ cd ..
$ cd ..
$ cd d
$ ls
4060174 j
8033020 d.log
5626152 d.ext
7214296 k`

func TestPartOne(t *testing.T) {
	lines := regexp.MustCompile("\r?\n").Split(testData, -1)
	value, err := PartOneSolution(lines)

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if value != 95437 {
		t.Errorf("Expected %d got %d", 95437, value)
	}
}

func TestPartTwo(t *testing.T) {
	lines := regexp.MustCompile("\r?\n").Split(testData, -1)
	value, err := PartTwoSolution(lines)

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if value != 24933642 {
		t.Errorf("Expected %d got %d", 24933642, value)
	}
}
