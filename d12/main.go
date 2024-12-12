package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spikecurtis/aoc2024/util"
)

type region struct {
	plant                  byte
	area, perimeter, sides int
}

func main() {
	//g := parseInput(getExample())
	g := parseInput(util.GetInputLines())
	regions, regGrid := findRegions(g)
	// corners are indexed by the grid location they are north-west of
	corners := util.Boundary2D[int]{
		MinX: g.Bounds.MinX,
		MinY: g.Bounds.MinY,
		MaxX: g.Bounds.MaxX + 1,
		MaxY: g.Bounds.MaxY + 1,
	}
	for corner := range corners.Iterate() {
		// go clockwise
		edges := make([]bool, 4)   // N, E, S, W
		regs := make([]*region, 4) // NW, NE, SE, SW
		regs[0], _ = regGrid.At(corner.Plus(util.NorthWest[int]()))
		regs[1], _ = regGrid.At(corner.Plus(util.North[int]()))
		regs[2], _ = regGrid.At(corner)
		regs[3], _ = regGrid.At(corner.Plus(util.West[int]()))
		edges[0] = regs[0] != regs[1]
		edges[1] = regs[1] != regs[2]
		edges[2] = regs[2] != regs[3]
		edges[3] = regs[3] != regs[0]
		// Perimeter:
		// only count edges from north and west to avoid double-counting
		// North
		if edges[0] {
			if regs[0] != nil {
				regs[0].perimeter++
			}
			if regs[1] != nil {
				regs[1].perimeter++
			}
		}
		// West
		if edges[3] {
			if regs[0] != nil {
				regs[0].perimeter++
			}
			if regs[3] != nil {
				regs[3].perimeter++
			}
		}
		// Sides
		for i := range 4 {
			j := (i + 1) % 4

			if edges[i] && edges[j] && regs[j] != nil {
				regs[j].sides++
				// if the next two edges, clockwise, are missing, this is a concave corner for the
				// next region.
				k := (i + 2) % 4
				l := (i + 3) % 4
				if !edges[k] && !edges[l] && regs[k] != nil {
					regs[k].sides++
				}
			}
		}
	}
	p1 := 0
	p2 := 0
	for _, r := range regions {
		p1 += r.perimeter * r.area
		p2 += r.sides * r.area
	}
	fmt.Println("Part 1: ", p1)
	fmt.Println("Part 2: ", p2)
}

func findRegions(g util.Grid2D[byte]) ([]*region, util.Grid2D[*region]) {
	regions := make([]*region, 0)
	regGrid := newEmptyRegionGrid(g.Bounds)
	for pnt, plant := range g.Iterate() {
		if r, _ := regGrid.At(pnt); r != nil {
			// already visited
			continue
		}
		r := &region{plant: plant, area: 1}
		regions = append(regions, r)
		err := regGrid.Set(pnt, r)
		if err != nil {
			log.Fatal(err)
		}
		q := []util.Point2D[int]{pnt}
		for len(q) > 0 {
			this := q[len(q)-1]
			q = q[:len(q)-1]
			for _, next := range this.CardinalNeighbors() {
				np, err := g.At(next)
				if err != nil {
					// OOB
					continue
				}
				if nr, _ := regGrid.At(next); nr != nil {
					// already visited
					continue
				}
				if np == plant {
					r.area++
					err := regGrid.Set(next, r)
					if err != nil {
						log.Fatal(err)
					}
					q = append(q, next)
				}
			}
		}
	}
	return regions, regGrid
}

func newEmptyRegionGrid(b util.Boundary2D[int]) util.Grid2D[*region] {
	values := make([][]*region, b.MaxY+1)
	for i := range values {
		values[i] = make([]*region, b.MaxX+1)
	}
	return util.Grid2D[*region]{
		Values: values,
		Bounds: b,
	}
}

func parseInput(lines []string) util.Grid2D[byte] {
	plants := make([][]byte, len(lines))
	for i := range plants {
		plants[i] = []byte(lines[i])
	}
	return util.Grid2D[byte]{
		Values: plants,
		Bounds: util.Boundary2D[int]{
			MinX: 0,
			MinY: 0,
			MaxX: len(plants[0]) - 1,
			MaxY: len(plants) - 1,
		},
	}
}

func getExample() []string {
	example := `RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE`
	return strings.Split(example, "\n")
}
