package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/spikecurtis/aoc2024/util"
)

var clockwise = map[util.Point2D[int]]util.Point2D[int]{
	util.North[int](): util.East[int](),
	util.East[int]():  util.South[int](),
	util.South[int](): util.West[int](),
	util.West[int]():  util.North[int](),
}

var counterclockwise = map[util.Point2D[int]]util.Point2D[int]{
	util.North[int](): util.West[int](),
	util.East[int]():  util.North[int](),
	util.South[int](): util.East[int](),
	util.West[int]():  util.South[int](),
}

type maze struct {
	wall        util.Grid2D[bool]
	start, goal util.Point2D[int]
}

func (m maze) IsGoal(n node) bool {
	return n.p == m.goal
}

func (m maze) Neighbors(n node) (neighbors []node, costs []int) {
	f := n.p.Plus(n.d)
	if w, _ := m.wall.At(f); !w {
		neighbors = append(neighbors, node{
			p: f,
			d: n.d,
		})
		costs = append(costs, 1)
	}
	neighbors = append(neighbors,
		node{
			p: n.p,
			d: clockwise[n.d],
		},
		node{
			p: n.p,
			d: counterclockwise[n.d],
		},
	)
	costs = append(costs, 1000, 1000)
	return neighbors, costs
}

func (m maze) Heuristic(n node) int {
	diff := m.goal.Minus(n.p)
	cost := diff.L1Norm()
	if cost == 0 {
		return 0
	}
	if diff.X == 0 {
		if util.Abs(n.d.X) == 1 {
			// facing E/W and need to go N/S
			return cost + 1000
		}
		if (diff.Y > 0 && n.d.Y > 0) || (diff.Y < 0 && n.d.Y < 0) {
			// facing the right direction
			return cost
		}
		// facing the wrong direction
		return cost + 2000
	}
	if diff.Y == 0 {
		if util.Abs(n.d.Y) == 1 {
			// facing N/S and need to go E/W
			return cost + 1000
		}
		if (diff.X > 0 && n.d.X > 0) || (diff.X < 0 && n.d.X < 0) {
			// facing the right direction
			return cost
		}
		// facing the wrong direction
		return cost + 2000
	}
	if (diff.X > 0 && n.p.X > 0) ||
		(diff.X < 0 && n.p.X < 0) ||
		(diff.Y > 0 && n.p.Y > 0) ||
		(diff.Y < 0 && n.p.Y < 0) {
		// going the right way in X or Y
		return cost + 1000
	}
	return cost + 2000
}

func (m maze) Start() node {
	return node{
		p: m.start,
		d: util.East[int](),
	}
}

type node struct {
	p, d util.Point2D[int]
}

func main() {
	m := parseInput(util.GetInputLines())
	//m := parseInput(getExampleLines())
	p1, _, err := util.AStar[node](m)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1: ", p1)
	p2 := numBestSeats(m)
	fmt.Println("Part 2: ", p2)
}

func parseInput(lines []string) maze {
	walls := make([][]bool, len(lines))
	var start, goal util.Point2D[int]
	for y, line := range lines {
		w := make([]bool, len(line))
		for x, b := range []byte(line) {
			switch b {
			case 'S':
				start.X = x
				start.Y = y
			case 'E':
				goal.X = x
				goal.Y = y
			case '#':
				w[x] = true
			case '.':
				// OK
			default:
				panic("unhandled")
			}
		}
		walls[y] = w
	}
	return maze{
		wall: util.Grid2D[bool]{
			Values: walls,
			Bounds: util.Boundary2D[int]{
				MinX: 0,
				MinY: 0,
				MaxX: len(walls[0]) - 1,
				MaxY: len(walls) - 1,
			},
		},
		start: start,
		goal:  goal,
	}
}

func getExampleLines() []string {
	return strings.Split(`###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############`, "\n")
}

// numBestSeats counts the best seats with a modifed A* that keeps searching after finding the goal
// to include any additional paths at equal cost.
func numBestSeats(m maze) int {
	start := m.Start()
	open := map[node]struct{}{
		start: {},
	}
	gScore := make(map[node]int)
	gScore[start] = 0
	fScore := make(map[node]int)
	fScore[start] = m.Heuristic(start)
	cameFrom := make(map[node][]node)

	countSeats := func(en map[node]struct{}) int {
		bestSeats := make(map[util.Point2D[int]]bool)
		nodesToCheck := make([]node, 0)
		for n := range en {
			nodesToCheck = append(nodesToCheck, n)
		}
		for len(nodesToCheck) > 0 {
			this := nodesToCheck[len(nodesToCheck)-1]
			nodesToCheck = nodesToCheck[:len(nodesToCheck)-1]
			bestSeats[this.p] = true
			ps := cameFrom[this]
			nodesToCheck = append(nodesToCheck, ps...)
		}
		return len(bestSeats)
	}

	bestNext := func() node {
		bestScore := math.MaxInt
		bestNode := start
		for n := range open {
			if score, ok := fScore[n]; ok && score < bestScore {
				bestScore = score
				bestNode = n
			}
		}
		if bestScore == math.MaxInt {
			panic("no best")
		}
		delete(open, bestNode)
		return bestNode
	}

	bestCost := math.MaxInt
	endNodes := make(map[node]struct{})
	for len(open) > 0 {
		current := bestNext()
		if m.IsGoal(current) {
			bestCost = gScore[current]
			fmt.Println("Best cost: ", bestCost, current)
			endNodes[current] = struct{}{}
			continue
		}
		neighbors, costs := m.Neighbors(current)
		if len(neighbors) != len(costs) {
			panic("unequal neighbors and costs")
		}
		for i, neighbor := range neighbors {
			edgeCost := costs[i]
			score := gScore[current] + edgeCost
			if oldScore, ok := gScore[neighbor]; !ok || score <= oldScore {
				if score == oldScore {
					cameFrom[neighbor] = append(cameFrom[neighbor], current)
				} else {
					cameFrom[neighbor] = []node{current}
				}

				gScore[neighbor] = score
				fs := score + m.Heuristic(neighbor)
				fScore[neighbor] = fs
				if fs <= bestCost {
					open[neighbor] = struct{}{}
				}
			}
		}
	}
	fmt.Println("found!")
	return countSeats(endNodes)
}
