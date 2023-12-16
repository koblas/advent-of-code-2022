package main

import (
	"bufio"
	// "errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const MAX_STEP = 30

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

type Node struct {
	name     string
	flow     int
	edges    []string
	isOpen   bool
	distance int
}

type State struct {
	valve       string
	valvesOpen  string
	minutesLeft int
}

type Graph map[string]*Node

var (
	lineRe = regexp.MustCompile(`^Valve (\w+) has flow rate=(\d+); tunnels? leads? to valves? (.*)`)
)

func ParseInput(lines []string) (Graph, error) {
	result := Graph{}

	for _, line := range lines {
		if parts := lineRe.FindStringSubmatch(line); parts != nil {
			flow, _ := strconv.ParseInt(parts[2], 10, 64)

			tparts := []string{}
			for _, val := range strings.Split(parts[3], ",") {
				tparts = append(tparts, strings.TrimSpace(val))
			}

			result[parts[1]] = &Node{
				name:     parts[1],
				flow:     int(flow),
				edges:    tparts,
				distance: -1,
				// isOpen:  flow == 0,
			}
		}
	}

	return result, nil
}

type Pair [2]string

func BFSSearch(graph Graph, current string) {
	visit := map[Pair]struct{}{}

	queue := []string{current}

	for len(queue) != 0 {
		head := queue[0]
		queue = queue[1:]

		node := graph[head]
		for _, connected := range node.edges {
			p := Pair{head, connected}
			if _, found := visit[p]; found {
				continue
			}
			visit[p] = struct{}{}
		}
	}
}

func ScoreXX(graph Graph, path []string) int {
	score := 0
	seen := map[string]struct{}{}

	for idx, key := range path {
		if _, found := seen[key]; found {
			continue
		}
		seen[key] = struct{}{}
		score += (len(path) - idx) * graph[key].flow
	}

	return score
}

func ResetDistance(graph Graph) {
	for _, node := range graph {
		node.distance = -1
	}
}

func SetDistance(graph Graph, current string) {
	dist := 0

	for queue := []string{current}; len(queue) != 0; dist += 1 {
		newq := []string{}

		for _, key := range queue {
			node := graph[key]
			if node.distance != -1 {
				continue
			}
			node.distance = dist
			newq = append(newq, node.edges...)
		}

		queue = newq
	}
}

func GetFlow(graph Graph) (int, []string) {
	total := 0
	open := []string{}
	for _, node := range graph {
		if node.isOpen {
			total += node.flow
			open = append(open, node.name)
		}
	}

	return total, open
}

// Priority Queue
type PriorityQueueItem struct {
	weight int
	node   *Node
}

type PriorityQueue []PriorityQueueItem

func (queue PriorityQueue) Add(node PriorityQueueItem) PriorityQueue {
	for idx := 0; idx < len(queue); idx += 1 {
		if queue[idx].weight > node.weight {
			continue
		}

		nq := PriorityQueue{}
		if idx > 0 {
			nq = append(nq, queue[0:idx]...)
		}
		nq = append(nq, node)
		nq = append(nq, queue[idx:]...)

		return nq
	}

	return append(queue, node)
}

// Keep this
func GetBestClosed(graph Graph, timeleft int) PriorityQueue {
	result := PriorityQueue{}
	for _, node := range graph {
		if node.isOpen || node.flow == 0 {
			continue
		}
		result = result.Add(PriorityQueueItem{
			weight: node.flow * (timeleft - node.distance - 1),
			node:   node,
		})
	}

	return result
}

func BFSrecurse(graph Graph, seen map[string]struct{}, current, best string, timeleft int) (int, string) {
	if _, found := seen[current+best]; found {
		return 0, ""
	}
	seen[current+best] = struct{}{}

	bestVal := 0
	bestPath := ""

	if timeleft < 0 {
		return bestVal, bestPath
	}

	node := graph[current]
	visited := []string{}

	// fmt.Println("HERE ", current, node.isOpen)

	if !node.isOpen && node.flow != 0 {
		graph[current].isOpen = true
		sVal, pVal := BFSrecurse(graph, seen, current, best, timeleft-1)
		bestVal = sVal + node.flow*(timeleft-1)
		bestPath = strings.ToLower(current) + pVal

		graph[current].isOpen = false
	}

	for _, key := range node.edges {
		if _, found := seen[current+key]; found {
			continue
		}
		visited = append(visited, current+key)
		seen[current+key] = struct{}{}

		sVal, pVal := BFSrecurse(graph, seen, key, best, timeleft-1)

		flow := 0
		if !graph[key].isOpen && graph[key].flow != 0 {
			flow = graph[key].flow * timeleft
		}
		scoreVal := sVal + flow
		if scoreVal > bestVal {
			bestVal = scoreVal
			bestPath = key + pVal
		}
	}

	for _, key := range visited {
		delete(seen, key)
	}

	delete(seen, current+best)

	return bestVal, bestPath
}

func BFSBestPath(graph Graph, current string, timeleft int) string {
	seen := map[string]struct{}{}

	tries := []string{}
	for key, node := range graph {
		if node.flow != 0 && !node.isOpen {
			tries = append(tries, key)
		}
	}

	bestScore := 0
	bestPath := ""
	for _, item := range tries {
		score, path := BFSrecurse(graph, seen, current, item, timeleft-1)
		if score > bestScore {
			bestScore = score
			bestPath = path
		}
	}

	fmt.Println("PATH = ", bestPath)

	if bestPath == "" {
		return ""
	}

	return bestPath[0:2]
}

func Best(graph Graph, current string, timeleft int) string {
	canidates := GetBestClosed(graph, timeleft)

	fmt.Printf("Canidates ")
	for _, item := range canidates {
		fmt.Printf("%s(flow=%d,weight=%d,dist=%d), ", item.node.name, item.node.flow, item.weight, item.node.distance)
	}
	fmt.Printf("\n")

	if len(canidates) == 0 {
		return current
	}

	return canidates[0].node.name
}

func DistanceBFS(graph Graph, current, dest string) (int, []string) {
	queue := []string{current}
	visited := map[string]struct{}{}

	dist := 0
	path := []string{}
	for len(queue) != 0 {
		front := queue[0]
		queue = queue[1:]

		path = append(path, front)
		if front == dest {
			fmt.Println("DONE ", current, dest, dist, path)
			return dist, path[1:]
		}

		for _, name := range graph[front].edges {
			if _, found := visited[name]; found {
				continue
			}
			visited[name] = struct{}{}
			queue = append(queue, name)
		}

		dist += 1
	}

	return dist, path[1:]
}

func Step(graph Graph, current string, timeleft int) string {
	// best := Best(graph, current, timeleft)

	next := BFSBestPath(graph, current, timeleft)
	if next == "" {
		return current
	}

	// dist, path := DistanceBFS(graph, current, best)

	// fmt.Println("   Steps to=", best, " dist=", dist, " path=", path)

	//if len(path) == 0 {
	//return current
	//}

	fmt.Println("You move to valve ", next)
	return next
}

func PartOneSolution(lines []string) (int, error) {
	graph, err := ParseInput(lines)
	if err != nil {
		return 0, err
	}

	maxTime := MAX_STEP
	pos := "AA"
	totalFlow := 0
	for step := 0; step < maxTime; step += 1 {
		ResetDistance(graph)
		SetDistance(graph, pos)

		fmt.Println("== Minute ", step+1, " ==")
		flow, open := GetFlow(graph)
		totalFlow += flow

		fmt.Println("Valves ", open, " are open, releasing ", flow, " pressure.")

		pos = Step(graph, pos, maxTime-step-1)

		upos := strings.ToUpper(pos)
		if pos == upos {
			fmt.Println("You move to valve ", pos)
		} else {
			pos = upos
			fmt.Println("You open valve ", pos)
			graph[pos].isOpen = true
		}
	}

	return totalFlow, nil
}

func PartTwoSolution(lines []string) (int, error) {
	//graph, err := ParseInput(lines)
	//if err != nil {
	//return 0, err
	//}

	return 0, nil
}
