package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
)

type lab [][]byte

type point struct {
	x, y int
}

type pointDir struct {
	p, d point
}

var (
	north = point{0, -1}
	east  = point{1, 0}
	south = point{0, 1}
	west  = point{-1, 0}
)

var right = map[point]point{
	north: east,
	east:  south,
	south: west,
	west:  north,
}

func (p point) add(other point) point {
	return point{p.x + other.x, p.y + other.y}
}

func (l lab) at(pnt point) byte {
	return l[pnt.y][pnt.x]
}

func (l lab) inBounds(pnt point) bool {
	if pnt.x < 0 || pnt.y < 0 {
		return false
	}
	maxY := len(l) - 1
	if pnt.y > maxY {
		return false
	}
	maxX := len(l[0]) - 1
	if pnt.x > maxX {
		return false
	}
	return true
}

func (l lab) withObstacle(p point) lab {
	out := make([][]byte, len(l))
	for y := range l {
		out[y] = slices.Clone(l[y])
	}
	out[p.y][p.x] = '#'
	return out
}

func parseInput(lines []string) (lab, pointDir) {
	var guard point
	out := make(lab, len(lines))
	for y, line := range lines {
		out[y] = []byte(line)
		for x := range out[y] {
			if out[y][x] == '^' {
				guard = point{x, y}
			}
		}
	}
	return out, pointDir{guard, north}
}

func main() {
	l, init := parseInput(getInput())
	guard := init
	p1 := make(map[point]bool)
	p2 := 0
	for {
		p1[guard.p] = true
		next := pointDir{guard.p.add(guard.d), guard.d}
		if !l.inBounds(next.p) {
			break
		}
		if l.at(next.p) == '#' {
			next = pointDir{guard.p, right[guard.d]}
		}
		guard = next
	}
	fmt.Println("Part 1:", len(p1))
	for p := range p1 {
		if init.p == p {
			continue
		}
		l2 := l.withObstacle(p)
		guard = init
		steps := make(map[pointDir]bool)
		for {
			if steps[guard] { // loop'd!
				p2++
				break
			}
			steps[guard] = true
			next := pointDir{guard.p.add(guard.d), guard.d}
			if !l2.inBounds(next.p) {
				break
			}
			if l2.at(next.p) == '#' {
				next = pointDir{guard.p, right[guard.d]}
			}
			guard = next
		}
	}
	fmt.Println("Part 2:", p2)
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
