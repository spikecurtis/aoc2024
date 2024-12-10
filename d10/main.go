package main

import (
	"fmt"
	"strings"

	"github.com/spikecurtis/aoc2024/util"
)

func main() {
	g := parseInput(util.GetInputLines())
	//g := parseInput(getExampleLines())
	p1 := 0
	p2 := 0
	for p, h := range g.Iterate() {
		if h != 0 {
			continue
		}
		this := map[util.Point2D[int]]int{
			p: 1,
		}
		for i := 1; i <= 9; i++ {
			this = findNeighborsAt(g, this, i)
		}
		p1 += len(this)
		for _, ways := range this {
			p2 += ways
		}
	}
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)
}

func findNeighborsAt(g util.Grid2D[int], this map[util.Point2D[int]]int, i int,
) map[util.Point2D[int]]int {
	next := make(map[util.Point2D[int]]int)
	for here, waysToHere := range this {
		for _, n := range here.CardinalNeighbors() {
			h, err := g.At(n)
			if err != nil {
				// OOB
				continue
			}
			if h == i {
				next[n] += waysToHere
			}
		}
	}
	return next
}

func parseInput(lines []string) util.Grid2D[int] {
	heights := make([][]int, len(lines))
	for i, line := range lines {
		chars := []byte(line)
		h := make([]int, len(chars))
		for j, char := range chars {
			h[j] = int(char - '0')
		}
		heights[i] = h
	}
	return util.Grid2D[int]{
		Values: heights,
		Bounds: util.Boundary2D[int]{
			MaxX: len(heights[0]) - 1,
			MaxY: len(heights) - 1,
			MinX: 0,
			MinY: 0,
		},
	}
}

func getExampleLines() []string {
	return strings.Split(`89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`, "\n")
}
