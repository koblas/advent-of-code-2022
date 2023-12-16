package main

import (
	"bufio"
	// "errors"
	"fmt"
	"image"
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

	value, err := PartOneSolution(lines, 2_000_000)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1: ", value)

	values, err := PartTwoSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2: ", values)
}

var directions = []image.Point{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
}

type Item struct {
	sensor image.Point
	beacon image.Point
	// delta  image.Point
	dist int
}

type Range [2]int

type Grid map[image.Point]rune

var (
	lineRe = regexp.MustCompile("^Sensor at x=(-?\\d+), y=(-?\\d+): closest beacon is at x=(-?\\d+), y=(-?\\d+)")
)

func ParseInput(lines []string) ([]Item, error) {
	result := []Item{}

	for _, line := range lines {
		if parts := lineRe.FindStringSubmatch(line); parts != nil {
			v1, _ := strconv.ParseInt(parts[1], 10, 64)
			v2, _ := strconv.ParseInt(parts[2], 10, 64)
			v3, _ := strconv.ParseInt(parts[3], 10, 64)
			v4, _ := strconv.ParseInt(parts[4], 10, 64)

			sensor := image.Point{int(v1), int(v2)}
			beacon := image.Point{int(v3), int(v4)}
			delta := sensor.Sub(beacon)
			dist := Abs(delta.X) + Abs(delta.Y)

			result = append(result, Item{
				sensor: sensor,
				beacon: beacon,
				dist:   dist,
			})
		}
	}

	return result, nil
}

func DrawItems(grid Grid, items []Item) {
	for _, item := range items {
		grid[item.sensor] = 'S'
		grid[item.beacon] = 'B'
	}
}

func RangeSensor(grid Grid, item Item, row *int) {
	dist := item.sensor.Sub(item.beacon)

	dX := Abs(dist.X)
	dY := Abs(dist.Y)
	dSum := dX + dY
	// fmt.Println("DSUM = ", dSum)

	//for y := 0; y <= dSum; y += 1 {
	//for x := 0; x <= dSum-y; x += 1 {
	//pos := item.sensor.Add(image.Point{x, y})
	//grid[pos] = '#'
	//pos = item.sensor.Add(image.Point{x, -y})
	//grid[pos] = '#'
	//pos = item.sensor.Add(image.Point{-x, y})
	//grid[pos] = '#'
	//pos = item.sensor.Add(image.Point{-x, -y})
	//grid[pos] = '#'
	//}
	//}
	for d := -dSum; d <= dSum; d += 1 {
		absD := Abs(d)
		if row != nil && *row != d+item.sensor.Y {
			continue
		}

		pos := item.sensor.Add(image.Point{-(dSum - absD), d})

		//if pos.Y < 0 || pos.Y > 4_000_000 {
		//continue
		//}
		for x := -(dSum - absD); x <= dSum-absD; x += 1 {
			// grid[image.Point{x, y}] = '#'
			// grid[item.sensor.Add(image.Point{x, d})] = '#'
			grid[pos] = '#'
			pos = pos.Add(image.Point{1, 0})
		}
	}
}

func RangeSensorRange(items []Item, ypos int) []Range {
	ranges := []Range{}

	for _, item := range items {
		diffY := Abs(item.sensor.Y - ypos)
		if diffY > item.dist {
			continue
		}
		distX := item.dist - diffY
		// fmt.Println("VALUES dist=", item.dist, " diffY=", diffY, " distX=", distX)

		ranges = append(ranges, Range{item.sensor.X - distX, item.sensor.X + distX})
	}

	return ranges
}

func MergeRanges(ranges []Range) []Range {
	result := []Range{}

	for _, item := range ranges {
		found := false

		w1 := item[1] - item[0] + 1
		for idx, value := range result {
			w2 := value[1] - value[0] + 1

			min := Min(value[0], item[0])
			max := Max(value[1], item[1])

			if w2+w1 >= max-min+1 {
				result[idx][0] = Min(value[0], item[0])
				result[idx][1] = Max(value[1], item[1])
				found = true
				break
			}

			//	if Max(value[0], item[0]) <= Min(value[1], item[1]) {
			//		result[idx][0] = Min(value[0], item[0])
			//		result[idx][1] = Max(value[1], item[1])
			//		found = true
			//		break
			//	}
		}

		if !found {
			result = append(result, item)
		}
	}

	cnt := len(result)
	if cnt == 1 || cnt == len(ranges) {
		return result
	}

	return MergeRanges(result)
}

func RangeSensorOld(grid Grid, item Item) {
	round := []image.Point{item.sensor}
	found := false

	seen := map[image.Point]struct{}{}

	for !found {
		next := []image.Point{}

		for _, step := range round {
			for _, dir := range directions {
				pos := step.Add(dir)
				found = found || pos.Eq(item.beacon)
				if _, hit := seen[pos]; hit {
					continue
				}

				next = append(next, pos)
				seen[pos] = struct{}{}

				if _, hit := grid[pos]; !hit {
					grid[pos] = '#'
				}
			}
		}

		// fmt.Printf("sensor=%s beacon=%s\n", item.sensor, item.beacon)
		// fmt.Printf("%s\n\n", Format(grid))

		round = next
	}
}

func Min[T int](a, b T) T {
	if a < b {
		return a
	}
	return b
}
func Max[T int](a, b T) T {
	if a > b {
		return a
	}
	return b
}
func Abs[T int](a T) T {
	if a < 0 {
		return -1 * a
	}
	return a
}

func Format(grid Grid, maxSize int) string {
	output := []string{}

	min := image.Point{}
	max := image.Point{}

	first := true
	for item := range grid {
		if first {
			first = false
			min.X = item.X
			min.Y = item.Y
			max.X = item.X
			max.Y = item.Y
		} else {
			min.X = Min(min.X, item.X)
			min.Y = Min(min.Y, item.Y)
			max.X = Max(max.X, item.X)
			max.Y = Max(max.Y, item.Y)
		}
	}

	for y := min.Y - 1; y < max.Y+2; y += 1 {
		if maxSize > 0 && (y < 0 || y > maxSize) {
			continue
		}
		row := fmt.Sprintf("%3d ", y)
		for x := min.X - 1; x < max.X+2; x += 1 {
			if maxSize > 0 && (x < 0 || x > maxSize) {
				continue
			}
			ch := "."
			if val, found := grid[image.Point{x, y}]; found {
				ch = string(val)
			}
			row += ch
		}
		output = append(output, row)
	}

	return strings.Join(output, "\n")
}

func PartOneSolution(lines []string, row int) (int, error) {
	input, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	grid := Grid{}
	for _, item := range input {
		// fmt.Printf("%s %s\n\n", item.sensor, item.beacon)
		RangeSensor(grid, item, &row)
		// RangeSensor(grid, item, nil)
		// fmt.Printf("%s %s\n%s\n\n", item.sensor, item.beacon, Format(grid))
	}
	// DrawItems(grid, input)

	for _, item := range input {
		delete(grid, item.sensor)
		delete(grid, item.beacon)
	}

	fmt.Println("Doing count ", len(grid))
	count := 0
	for pos, item := range grid {
		if pos.Y != row {
			continue
		}
		if item == '#' {
			count += 1
		}
	}

	fmt.Printf("%s\n", Format(grid, 20))

	return count, nil
}

func PartTwoSolution(lines []string) (int, error) {
	input, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	for y := 0; y < 4_000_000; y += 1 {
		ranges := RangeSensorRange(input, y)
		// fmt.Println("RANGES ", y, " = ", ranges)
		// fmt.Println("RANGES CLEAN", MergeRanges(ranges))

		clean := MergeRanges(ranges)

		// fmt.Println("Y= ", y, "RANGES=", ranges, "CLEAN=", clean, "HIT = ", len(clean) != 1)
		if len(clean) != 1 {
			fmt.Println("Y= ", y, "RANGES=", ranges, "CLEAN=", clean)
			fmt.Println(clean)

			return (clean[0][1]+1)*4000000 + y, nil
		}
	}

	//grid := Grid{}
	//for _, item := range input {
	//RangeSensor(grid, item, nil)
	//}
	//DrawItems(grid, input)
	//fmt.Printf("%s\n", Format(grid, -1))

	return 0, nil
}
