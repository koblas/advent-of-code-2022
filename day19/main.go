package main

import (
	"bufio"
	"regexp"
	"strconv"

	// "errors"
	"fmt"
	"os"
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

	values, err := PartTwoSolution(lines)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2: ", values)
}

var (
	// Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 17 clay. Each geode robot costs 4 ore and 20 obsidian.
	lineRe = regexp.MustCompile(`(\d+)`)
)

const (
	ORE      = 0
	CLAY     = 1
	OBSIDIAN = 2
	GEODE    = 3
)

type ValueBase [4]int
type BuildCost ValueBase
type Robots ValueBase
type Resources ValueBase

func (r Robots) Add(kind int) Robots {
	val := r
	val[kind] += 1
	return val
}

type State struct {
	step      int
	robots    Robots
	resources Resources
}

func (s State) ToKey() uint64 {
	v := uint64(s.robots[0])<<0 |
		uint64(s.robots[1])<<6 |
		uint64(s.robots[2])<<12 |
		uint64(s.robots[3])<<18 |
		uint64(s.step)<<24 |
		uint64(s.resources[0])<<32 |
		uint64(s.resources[1])<<40 |
		uint64(s.resources[2])<<48 |
		uint64(s.resources[3])<<56

	return v
}

func (state State) Build(cost BuildCost, kind int) State {
	return State{
		step:      state.step - 1,
		resources: state.resources.Build(state.robots).Spend(cost),
		robots:    state.robots.Add(kind),
	}

}

type Blueprint struct {
	number int

	robotCost []BuildCost
}

func ParseInput(lines []string) ([]Blueprint, error) {
	result := []Blueprint{}
	for _, line := range lines {
		parts := []int{}

		split := lineRe.FindAllString(line, -1)
		if len(split) != 7 {
			panic("PARSE FAILED: " + line)
		}
		for _, part := range split {
			value, _ := strconv.ParseInt(part, 10, 0)
			parts = append(parts, int(value))
		}

		result = append(result, Blueprint{
			number: parts[0],
			robotCost: []BuildCost{
				{parts[1]},
				{parts[2]},
				{parts[3], parts[4]},
				{parts[5], 0, parts[6]},
			},
		})
	}

	return result, nil
}

func CanBuild(cost BuildCost, resources Resources) bool {
	for idx := 0; idx < 3; idx += 1 {
		if resources[idx] < cost[idx] {
			return false
		}
	}
	return true
}

func (resources Resources) Spend(cost BuildCost) Resources {
	return Resources{
		resources[0] - cost[0],
		resources[1] - cost[1],
		resources[2] - cost[2],
		resources[3],
	}
}

func (resources Resources) Build(robots Robots) Resources {
	return Resources{
		resources[0] + robots[0],
		resources[1] + robots[1],
		resources[2] + robots[2],
		resources[3] + robots[3],
	}
}

type Memo map[uint64]int

func StepSearch(input Blueprint, robots Robots, resources Resources, step int) int {
	seen := Memo{}
	queue := []State{
		{
			robots: robots,
			step:   step,
		},
	}

	maxOreRobotCost := 0
	for idx := 0; idx < 4; idx += 1 {
		val := input.robotCost[idx][ORE]
		if val > maxOreRobotCost {
			maxOreRobotCost = val
		}
	}

	best := 0

	for len(queue) != 0 {
		state := queue[0]
		queue = queue[1:]

		if best < state.resources[GEODE] {
			best = state.resources[GEODE]
		}
		if state.step == 0 {
			continue
		}
		key := state.ToKey()
		if _, found := seen[key]; found {
			continue
		}
		seen[key] = 0

		if CanBuildRobot(input, state.resources, GEODE) {
			// fmt.Println("BUILD GEODE", state.step)
			queue = append(queue, state.Build(input.robotCost[GEODE], GEODE))
		} else {
			if CanBuildRobot(input, state.resources, ORE) && state.robots[ORE] < maxOreRobotCost {
				// fmt.Println("BUILD ORE", state.step)
				queue = append(queue, state.Build(input.robotCost[ORE], ORE))
			}
			if CanBuildRobot(input, state.resources, CLAY) {
				// fmt.Println("BUILD CLAY", state.step)
				queue = append(queue, state.Build(input.robotCost[CLAY], CLAY))
			}
			if CanBuildRobot(input, state.resources, OBSIDIAN) {
				// fmt.Println("BUILD OBSIDIAN", state.step)
				queue = append(queue, state.Build(input.robotCost[OBSIDIAN], OBSIDIAN))
			}

			queue = append(queue, State{
				step:      state.step - 1,
				robots:    state.robots,
				resources: state.resources.Build(state.robots),
			})
		}
	}
	// fmt.Println("DEPTH ", step, robots, resources

	// if step == 1 {
	// 	return resources.Build(robots)[GEODE]
	// }

	// key := ValueBase(robots).ToString() + "|" + ValueBase(resources).ToString() + "|" + strconv.FormatInt(int64(step), 10)
	// if value, found := memo[key]; found {
	// 	return value
	// }

	// best := StepSearch(memo, input, robots, resources.Build(robots), step-1)
	// for idx := 3; idx >= 0; idx -= 1 {
	// 	if !CanBuild(input.robotCost[idx], resources) {
	// 		continue
	// 	}

	// 	newRobots := Robots{robots[0], robots[1], robots[2], robots[3]}
	// 	newRobots[idx] += 1

	// 	if newRobots[ORE] > 3 {
	// 		continue
	// 	}
	// 	if newRobots[CLAY] > 6 {
	// 		continue
	// 	}
	// 	if newRobots[OBSIDIAN] > 6 {
	// 		continue
	// 	}

	// 	newResources := resources.Build(robots).Spend(input.robotCost[idx])

	// 	value := StepSearch(memo, input, newRobots, newResources, step-1)
	// 	if value > best {
	// 		best = value
	// 	}
	// }

	// memo[key] = best

	return best
}

func Simulate(input Blueprint, robots Robots, resources Resources, maxSteps int) int {
	return StepSearch(input, robots, resources, maxSteps)
}

func CanBuildRobot(input Blueprint, resources Resources, robot int) bool {
	flag := (resources[ORE] >= input.robotCost[robot][ORE]) &&
		(resources[CLAY] >= input.robotCost[robot][CLAY]) &&
		(resources[OBSIDIAN] >= input.robotCost[robot][OBSIDIAN])

	return flag
}

// func SimulateTwo(input Blueprint, robots Robots, resources Resources, maxSteps int) int {
// 	var sRobots []Robots

// 	for attempt := 10; attempt > 0; attempt -= 10 {
// 		// Init to zero
// 		sRobots = make([]Robots, maxSteps)

// 		for step := 0; step < maxSteps; step += 1 {
// 			if CanBuildRobot(input, resources[step], GEODE) {
// 				panic(fmt.Sprintf("CAN BUILD AT STEP %d", step))
// 			}
// 		}
// 	}

// 	sum := 0
// 	for idx := 0; idx < maxSteps; idx += 1 {
// 		sum += sRobots[idx][GEODE]
// 	}
// 	return sum
// }

func PartOneSolution(lines []string) (int, error) {
	input, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	total := (0)
	for _, print := range input {
		robots := Robots{1}
		resources := Resources{}
		total += Simulate(print, robots, resources, 24) * print.number
	}

	return int(total), err
}

func PartTwoSolution(lines []string) (int, error) {
	input, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	total := (1)
	for idx := 0; idx < 3 && idx < len(input); idx += 1 {
		robots := Robots{1}
		resources := Resources{}
		value := Simulate(input[idx], robots, resources, 32)
		fmt.Println("INPUT ", idx, value)
		total *= value
	}

	return int(total), err
}
