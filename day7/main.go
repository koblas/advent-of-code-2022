package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	fd, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(fd)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	value, err := PartOneSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1: ", value)

	value, err = PartTwoSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2: ", value)
}

var cdRe = regexp.MustCompile("^\\$ cd (.*)$")
var lsRe = regexp.MustCompile("^\\$ ls$")
var dirRe = regexp.MustCompile("^dir (.*)$")
var fileRe = regexp.MustCompile("^(\\d+) (.*)$")

type File struct {
	name string
	size int
}

type Directory struct {
	name     string
	children []*Directory
	files    []File
	parent   *Directory
}

func newDir(name string, parent *Directory) *Directory {
	dir := Directory{
		name:     name,
		children: []*Directory{},
		files:    []File{},
		parent:   parent,
	}

	if parent == nil {
		dir.parent = &dir
	} else {
		parent.children = append(parent.children, &dir)
	}

	return &dir
}

func BuildTree(lines []string) (*Directory, error) {
	root := newDir("", nil)

	curwd := root
	for _, line := range lines {
		if parts := cdRe.FindStringSubmatch(line); parts != nil {
			// fmt.Println("CD", parts[1])
			if parts[1] == ".." {
				curwd = curwd.parent
			} else {
				for _, child := range curwd.children {
					// fmt.Println("LOOKING ", parts[1], child.name)
					if parts[1] == child.name {
						curwd = child
						break
					}
				}
			}
		} else if parts := lsRe.FindStringSubmatch(line); parts != nil {
			// Skip
		} else if parts := dirRe.FindStringSubmatch(line); parts != nil {
			newDir(parts[1], curwd)
		} else if parts := fileRe.FindStringSubmatch(line); parts != nil {
			size, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, err
			}
			curwd.files = append(curwd.files, File{
				name: parts[2],
				size: int(size),
			})
		} else {
			fmt.Println("No match ", line)
		}
	}

	return root, nil
}

func PrintTree(curwd *Directory, indent int) {
	spaces := strings.Repeat(" ", indent*2)
	for _, item := range curwd.children {
		size, _ := WalkPartOne(item)
		chars := ""
		if size >= 2_677_139 {
			chars = " ===="
		}
		fmt.Printf("%s- %s (dir, size=%d%s)\n", spaces, item.name, size, chars)
		PrintTree(item, indent+1)
	}
	for _, item := range curwd.files {
		fmt.Printf("%s- %s (file, size=%d)\n", spaces, item.name, item.size)
	}
}

func WalkPartOne(curwd *Directory) (int, int) {
	total := 0
	matchTotal := 0
	for _, item := range curwd.children {
		subTotal, subMatch := WalkPartOne(item)
		total += subTotal
		matchTotal += subMatch
	}
	for _, item := range curwd.files {
		total += item.size
	}
	if total <= 100_000 {
		return total, matchTotal + total
	}
	return total, matchTotal
}

func WalkPartTwo(curwd *Directory, need int, best int) (int, int) {
	total := 0
	for _, item := range curwd.children {
		subTotal, subBest := WalkPartTwo(item, need, best)
		total += subTotal

		if subBest >= need && subBest < best {
			best = subBest
		}
	}
	for _, item := range curwd.files {
		total += item.size
	}

	if total >= need && total < best {
		return total, total
	}
	return total, best
}

func PartOneSolution(lines []string) (int, error) {
	tree, err := BuildTree(lines)
	if err != nil {
		return 0, err
	}

	PrintTree(tree, 0)

	_, total := WalkPartOne(tree)

	return total, nil
}

func PartTwoSolution(lines []string) (int, error) {
	tree, err := BuildTree(lines)
	if err != nil {
		return 0, err
	}

	total, _ := WalkPartOne(tree)

	free := 70000000 - total
	need := 30000000 - free

	fmt.Println("NEED ", need, total)

	_, best := WalkPartTwo(tree, need, total)

	return best, nil
}
