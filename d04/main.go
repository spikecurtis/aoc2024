package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type puzzle [][]byte

type point struct {
	x, y int
}

var directions = []point{
	{0, 1},   // south
	{1, 1},   // southeast
	{1, 0},   // east
	{1, -1},  // northeast
	{0, -1},  // north
	{-1, -1}, // northwest
	{-1, 0},  // west
	{-1, 1},  // southwest
}

var xmas = []byte("XMAS")

func (p point) add(other point) point {
	return point{p.x + other.x, p.y + other.y}
}

func (p puzzle) letterAt(pnt point) byte {
	return p[pnt.y][pnt.x]
}

func (p puzzle) inBounds(pnt point) bool {
	if pnt.x < 0 || pnt.y < 0 {
		return false
	}
	maxY := len(p) - 1
	if pnt.y > maxY {
		return false
	}
	maxX := len(p[0]) - 1
	if pnt.x > maxX {
		return false
	}
	return true
}

func (p puzzle) xmasAt(pnt, dir point) bool {
	for i := range 4 {
		if !p.inBounds(pnt) {
			return false
		}
		if p.letterAt(pnt) != xmas[i] {
			return false
		}
		pnt = pnt.add(dir)
	}
	return true
}

func (p puzzle) p2XMasAt(pnt point) bool {
	if p.letterAt(pnt) != 'A' {
		return false
	}
	nw := pnt.add(point{-1, -1})
	if !p.inBounds(nw) {
		return false
	}
	se := pnt.add(point{1, 1})
	if !p.inBounds(se) {
		return false
	}
	// if both NW and SE are in bounds, then everything must be
	nwl := p.letterAt(nw)
	if !(nwl == 'M' || nwl == 'S') {
		return false
	}
	var sel byte = 'S'
	if nwl == 'S' {
		sel = 'M'
	}
	if p.letterAt(se) != sel {
		return false
	}
	// NE by SW "MAS"
	ne := pnt.add(point{1, -1})
	sw := pnt.add(point{-1, 1})
	nel := p.letterAt(ne)
	if !(nel == 'M' || nel == 'S') {
		return false
	}
	var swl byte = 'S'
	if nel == 'S' {
		swl = 'M'
	}
	if p.letterAt(sw) != swl {
		return false
	}
	return true
}

func main() {
	p := parseInput(getInput())
	p1 := 0
	p2 := 0
	for y := range p {
		for x := range p[y] {
			pnt := point{x, y}
			for _, dir := range directions {
				if p.xmasAt(pnt, dir) {
					p1++
				}
			}
			if p.p2XMasAt(pnt) {
				p2++
			}
		}
	}
	fmt.Println("Part 1:", p1)
	fmt.Println("Part 2:", p2)
}

func parseInput(lines []string) puzzle {
	out := make(puzzle, len(lines))
	for i, line := range lines {
		out[i] = []byte(line)
	}
	return out
}

func getInput() []string {
	f, err := os.Open("./input")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}
