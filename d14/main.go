package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spikecurtis/aoc2024/util"
)

var parseRegex = regexp.MustCompile(`p=(\d+),(\d+)\sv=(-?\d+),(-?\d+)`)

type robot struct {
	p, v util.Point2D[int]
}

func (r robot) Move(b util.Boundary2D[int], steps int) robot {
	x := r.p.X + (steps * r.v.X)
	x = x % (b.MaxX + 1)
	if x < 0 {
		x = x + b.MaxX + 1
	}
	y := r.p.Y + (steps * r.v.Y)
	y = y % (b.MaxY + 1)
	if y < 0 {
		y = y + b.MaxY + 1
	}
	return robot{
		p: util.Point2D[int]{X: x, Y: y},
		v: r.v,
	}
}

func main() {
	//b := util.Boundary2D[int]{
	//	MinX: 0,
	//	MinY: 0,
	//	MaxX: 10,
	//	MaxY: 6,
	//}
	//robots := parseInput(getExampleLines())
	b := util.Boundary2D[int]{
		MinX: 0,
		MinY: 0,
		MaxX: 100,
		MaxY: 102,
	}
	robots := parseInput(util.GetInputLines())
	midX := b.MaxX / 2
	midY := b.MaxY / 2
	quads := make([]int, 4) // NW, NE, SE, SW
	for _, r := range robots {
		r = r.Move(b, 100)
		if r.p.X < midX { // West
			if r.p.Y < midY { // North
				quads[0]++
			}
			if r.p.Y > midY { // South
				quads[3]++
			}
		}
		if r.p.X > midX { // East
			if r.p.Y < midY { // North
				quads[1]++
			}
			if r.p.Y > midY { // South
				quads[2]++
			}
		}
	}
	p1 := 1
	for _, q := range quads {
		p1 *= q
	}
	fmt.Println("Part 1: ", p1)
	for i := range 100000 {
		if i > 7300 {
			fmt.Printf("\n\n%d seconds\n", i)
			printRobots(b, robots)
			time.Sleep(500 * time.Millisecond)
		}
		for r := range robots {
			robots[r] = robots[r].Move(b, 1)
		}
	}
}

func printRobots(b util.Boundary2D[int], robots []robot) {
	out := make([][]int, b.MaxY+1)
	for i := range out {
		out[i] = make([]int, b.MaxX+1)
	}
	for _, r := range robots {
		out[r.p.Y][r.p.X]++
	}
	for i := range out {
		bld := strings.Builder{}
		for _, n := range out[i] {
			if n == 0 {
				bld.WriteString(".")
				continue
			}
			bld.WriteString(strconv.Itoa(n))
		}
		fmt.Println(bld.String())
	}
}

func parseInput(lines []string) []robot {
	out := make([]robot, len(lines))
	for i, line := range lines {
		groups := parseRegex.FindStringSubmatch(line)
		if groups == nil {
			log.Fatalf("Error parsing line %d: %s", i, line)
		}
		x, err := strconv.Atoi(groups[1])
		if err != nil {
			log.Fatalf("Error parsing line %d: %s", i, line)
		}
		y, err := strconv.Atoi(groups[2])
		if err != nil {
			log.Fatalf("Error parsing line %d: %s", i, line)
		}
		vx, err := strconv.Atoi(groups[3])
		if err != nil {
			log.Fatalf("Error parsing line %d: %s", i, line)
		}
		vy, err := strconv.Atoi(groups[4])
		if err != nil {
			log.Fatalf("Error parsing line %d: %s", i, line)
		}
		out[i] = robot{
			p: util.Point2D[int]{
				X: x,
				Y: y,
			},
			v: util.Point2D[int]{
				X: vx,
				Y: vy,
			},
		}
	}
	return out
}

func getExampleLines() []string {
	return strings.Split(`p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3`, "\n")
}
