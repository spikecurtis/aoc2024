package main

import (
	"fmt"

	"github.com/spikecurtis/aoc2024/util"
)

func main() {
	locs, bounds := parseInput(util.GetInputLines())
	antinodes1 := make(map[util.Point2D[int]]bool)
	antinodes2 := make(map[util.Point2D[int]]bool)
	for _, points := range locs {
		for i := range points {
			for j := i + 1; j < len(points); j++ {
				d := points[j].Minus(points[i])
				a := points[i]
				for n := 0; ; n++ {
					b := a.Plus(d.Scale(n))
					if !bounds.In(b) {
						break
					}
					antinodes2[b] = true
					if n == 2 {
						antinodes1[b] = true
					}
				}
				for n := -1; ; n-- {
					b := a.Plus(d.Scale(n))
					if !bounds.In(b) {
						break
					}
					antinodes2[b] = true
					if n == -1 {
						antinodes1[b] = true
					}
				}
			}
		}
	}
	fmt.Println("Part 1:", len(antinodes1))
	fmt.Println("Part 2:", len(antinodes2))
}

func parseInput(lines []string) (map[byte][]util.Point2D[int], util.Boundary2D[int]) {
	out := make(map[byte][]util.Point2D[int])
	for y, line := range lines {
		for x, c := range []byte(line) {
			if c != '.' {
				out[c] = append(out[c], util.Point2D[int]{x, y})
			}
		}
	}
	return out, util.Boundary2D[int]{
		MinX: 0,
		MinY: 0,
		MaxX: len([]byte(lines[0])) - 1,
		MaxY: len(lines) - 1,
	}
}
